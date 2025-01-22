// Copyright © 2020 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package javascript contains the Javascript payload formatter message processors.
package javascript

import (
	"context"
	"fmt"
	"reflect"
	"runtime/trace"
	"strings"
	"time"

	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/goproto"
	"go.thethings.network/lorawan-stack/v3/pkg/messageprocessors"
	"go.thethings.network/lorawan-stack/v3/pkg/messageprocessors/normalizedpayload"
	"go.thethings.network/lorawan-stack/v3/pkg/scripting"
	js "go.thethings.network/lorawan-stack/v3/pkg/scripting/javascript"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"google.golang.org/protobuf/types/known/structpb"
)

type host struct {
	engine scripting.AheadOfTimeEngine
}

// New creates and returns a new Javascript payload encoder and decoder.
func New() messageprocessors.CompilablePayloadEncoderDecoder {
	return &host{
		engine: js.New(scripting.DefaultOptions),
	}
}

type encodeDownlinkInput struct {
	Data  map[string]any `json:"data"`
	FPort *uint8         `json:"fPort"`
}

type encodeDownlinkOutput struct {
	Bytes    []uint8  `json:"bytes"`
	FPort    *uint8   `json:"fPort"`
	Warnings []string `json:"warnings"`
	Errors   []string `json:"errors"`
}

var (
	errInput          = errors.DefineInvalidArgument("input", "invalid input")
	errOutput         = errors.Define("output", "invalid output")
	errOutputErrors   = errors.DefineAborted("output_errors", "{errors}")
	errOutputEncoding = errors.DefineInvalidArgument("output_encoding", "{errors}")
)

func wrapDownlinkEncoderScript(script string) string {
	// Fallback to Encoder() for backwards compatibility with The Things Network Stack V2 payload functions.
	return fmt.Sprintf(`
		%s

		function main(input) {
			const { data, fPort } = input;
			if (typeof encodeDownlink === 'function') {
				return encodeDownlink({ data, fPort });
			}
			return {
				bytes: Encoder(data, fPort),
				fPort: input.fPort
			}
		}
	`, script)
}

// CompileDownlinkEncoder generates a downlink encoder from the provided script.
func (h *host) CompileDownlinkEncoder(
	ctx context.Context, script string,
) (
	func(
		context.Context,
		*ttnpb.EndDeviceIdentifiers,
		*ttnpb.EndDeviceVersionIdentifiers,
		*ttnpb.ApplicationDownlink,
	) error,
	error,
) {
	defer trace.StartRegion(ctx, "compile downlink encoder").End()

	run, err := h.engine.Compile(ctx, wrapDownlinkEncoderScript(script))
	if err != nil {
		return nil, err
	}

	return func(
		ctx context.Context,
		_ *ttnpb.EndDeviceIdentifiers,
		_ *ttnpb.EndDeviceVersionIdentifiers,
		msg *ttnpb.ApplicationDownlink,
	) error {
		return h.encodeDownlink(ctx, msg, run)
	}, nil
}

// EncodeDownlink encodes the message's DecodedPayload to FRMPayload using the given script.
func (h *host) EncodeDownlink(
	ctx context.Context,
	_ *ttnpb.EndDeviceIdentifiers,
	_ *ttnpb.EndDeviceVersionIdentifiers,
	msg *ttnpb.ApplicationDownlink,
	script string,
) error {
	run := func(ctx context.Context, fn string, params ...any) (func(any) error, error) {
		return h.engine.Run(ctx, wrapDownlinkEncoderScript(script), fn, params...)
	}
	return h.encodeDownlink(ctx, msg, run)
}

func (*host) encodeDownlink(
	ctx context.Context,
	msg *ttnpb.ApplicationDownlink,
	run func(context.Context, string, ...any) (func(any) error, error),
) error {
	defer trace.StartRegion(ctx, "encode downlink message").End()

	decoded := msg.DecodedPayload
	if decoded == nil {
		return nil
	}
	data, err := goproto.Map(decoded)
	if err != nil {
		return errInput.WithCause(err)
	}
	fPort := uint8(msg.FPort) // nolint:gosec
	input := encodeDownlinkInput{
		Data:  data,
		FPort: &fPort,
	}

	valueAs, err := run(ctx, "main", input)
	if err != nil {
		return err
	}

	var output encodeDownlinkOutput
	err = valueAs(&output)
	if err != nil {
		return errOutput.WithCause(err)
	}
	if len(output.Errors) > 0 {
		return errOutputErrors.WithAttributes("errors", strings.Join(output.Errors, ", "))
	}

	msg.FrmPayload = output.Bytes
	msg.DecodedPayloadWarnings = output.Warnings
	if output.FPort != nil {
		fPort := *output.FPort
		msg.FPort = uint32(fPort)
	} else if msg.FPort == 0 {
		msg.FPort = 1
	}
	return nil
}

type decodeUplinkInput struct {
	Bytes    []uint8 `json:"bytes"`
	FPort    uint8   `json:"fPort"`
	RecvTime int64   `json:"recvTime"` // UnixNano
}

type decodeUplinkOutput struct {
	Data     map[string]any `json:"data"`
	Warnings []string       `json:"warnings"`
	Errors   []string       `json:"errors"`
}

type normalizeUplinkOutput struct {
	Data     any      `json:"data"`
	Warnings []string `json:"warnings"`
	Errors   []string `json:"errors"`
}

type uplinkDecoderOutput struct {
	Decoded    decodeUplinkOutput     `json:"decoded"`
	Normalized *normalizeUplinkOutput `json:"normalized"`
}

func wrapUplinkDecoderScript(script string) string {
	// This wrapper executes decodeUplink() if it is defined. Then, it executes normalizeUplink() if it is defined too,
	// and if the output of decodeUplink() didn't return errors.
	// Fallback to Decoder() for backwards compatibility with The Things Network Stack V2 payload functions.
	return fmt.Sprintf(`
		%s

		function main(input) {
			const bytes = input.bytes.slice();
			const { fPort, recvTime } = input;

			// Convert UnixNano to JavaScript Date.
			const jsDate = new Date(Number(BigInt(recvTime) / 1000000n));

			if (typeof decodeUplink === 'function') {
				const decoded = decodeUplink({ bytes, fPort, recvTime: jsDate });
				let normalized;
				const { data, errors } = decoded;
				if ((!errors || !errors.length) && data && typeof normalizeUplink === 'function') {
					normalized = normalizeUplink({ data });
				}
				return { decoded, normalized };
			}
			return {
				decoded: {
					data: Decoder(bytes, fPort)
				}
			}
		}
	`, script)
}

// CompileUplinkDecoder generates an uplink decoder from the provided script.
func (h *host) CompileUplinkDecoder(
	ctx context.Context, script string,
) (
	func(
		context.Context,
		*ttnpb.EndDeviceIdentifiers,
		*ttnpb.EndDeviceVersionIdentifiers,
		*ttnpb.ApplicationUplink,
	) error,
	error,
) {
	defer trace.StartRegion(ctx, "compile uplink decoder").End()

	run, err := h.engine.Compile(ctx, wrapUplinkDecoderScript(script))
	if err != nil {
		return nil, err
	}

	return func(
		ctx context.Context,
		_ *ttnpb.EndDeviceIdentifiers,
		_ *ttnpb.EndDeviceVersionIdentifiers,
		msg *ttnpb.ApplicationUplink,
	) error {
		return h.decodeUplink(ctx, msg, run)
	}, nil
}

// DecodeUplink decodes the message's FRMPayload to DecodedPayload using the given script.
func (h *host) DecodeUplink(
	ctx context.Context,
	_ *ttnpb.EndDeviceIdentifiers,
	_ *ttnpb.EndDeviceVersionIdentifiers,
	msg *ttnpb.ApplicationUplink,
	script string,
) error {
	run := func(ctx context.Context, fn string, params ...any) (func(any) error, error) {
		return h.engine.Run(ctx, wrapUplinkDecoderScript(script), fn, params...)
	}
	return h.decodeUplink(ctx, msg, run)
}

func appendValidationErrors(dst []string, measurements []normalizedpayload.ParsedMeasurement) []string {
	for i, m := range measurements {
		for _, err := range m.ValidationErrors {
			var (
				errString string
				ttnErr    *errors.Error
			)
			if errors.As(err, &ttnErr) {
				errString = ttnErr.FormatMessage(ttnErr.PublicAttributes())
			} else {
				errString = err.Error()
			}
			dst = append(dst, fmt.Sprintf("measurement %d: %s", i+1, errString))
		}
	}
	return dst
}

func (*host) decodeUplink( // nolint: gocyclo
	ctx context.Context,
	msg *ttnpb.ApplicationUplink,
	run func(context.Context, string, ...any) (func(any) error, error),
) error {
	defer trace.StartRegion(ctx, "decode uplink message").End()

	input := decodeUplinkInput{
		Bytes:    msg.FrmPayload,
		FPort:    uint8(msg.FPort), // nolint:gosec
		RecvTime: msg.ReceivedAt.AsTime().UnixNano(),
	}

	valueAs, err := run(ctx, "main", input)
	if err != nil {
		return err
	}

	var output uplinkDecoderOutput
	err = valueAs(&output)
	if err != nil {
		return errOutput.WithCause(err)
	}

	if errs := output.Decoded.Errors; len(errs) > 0 {
		return errOutputErrors.WithAttributes("errors", strings.Join(errs, ", "))
	}

	// goproto.Struct does not support time.Time, use UnixNano instead.
	for key, item := range output.Decoded.Data {
		if t, ok := item.(time.Time); ok {
			output.Decoded.Data[key] = t.UnixNano()
		}
	}

	decodedPayload, err := goproto.Struct(output.Decoded.Data)
	if err != nil {
		return errOutput.WithCause(err)
	}
	if errs := goproto.ValidateStruct(decodedPayload); len(errs) > 0 {
		return errOutputEncoding.WithAttributes("errors", strings.Join(errs, ", "))
	}

	msg.DecodedPayload, msg.DecodedPayloadWarnings = decodedPayload, output.Decoded.Warnings
	msg.NormalizedPayload, msg.NormalizedPayloadWarnings = nil, nil

	if normalized := output.Normalized; normalized != nil {
		if errs := normalized.Errors; len(errs) > 0 {
			return errOutputErrors.WithAttributes("errors", strings.Join(errs, ", "))
		}
		if normalized.Data == nil {
			return nil
		}
		// The returned data can be an array of measurements or a single measurement object.
		var measurements []map[string]any
		if val := reflect.ValueOf(normalized.Data); val.Kind() == reflect.Slice {
			measurements = make([]map[string]any, val.Len())
			for i := 0; i < val.Len(); i++ {
				measurement, ok := val.Index(i).Interface().(map[string]any)
				if !ok {
					return errOutput.New()
				}
				measurements[i] = measurement
			}
		} else {
			measurement, ok := normalized.Data.(map[string]any)
			if !ok {
				return errOutput.New()
			}
			measurements = []map[string]any{measurement}
		}
		normalizedPayload := make([]*structpb.Struct, len(measurements))
		for i := range measurements {
			pb, err := goproto.Struct(measurements[i])
			if err != nil {
				return errOutput.WithCause(err)
			}
			normalizedPayload[i] = pb
		}
		// Validate the normalized payload.
		normalizedMeasurements, err := normalizedpayload.Parse(normalizedPayload)
		if err != nil {
			return errOutput.WithCause(err)
		}
		msg.NormalizedPayload = make([]*structpb.Struct, 0, len(normalizedMeasurements))
		for _, measurement := range normalizedMeasurements {
			if len(measurement.Valid.GetFields()) == 0 {
				continue
			}
			msg.NormalizedPayload = append(msg.NormalizedPayload, measurement.Valid)
		}
		msg.NormalizedPayloadWarnings = make([]string, 0, len(normalized.Warnings))
		msg.NormalizedPayloadWarnings = append(msg.NormalizedPayloadWarnings, normalized.Warnings...)
		msg.NormalizedPayloadWarnings = appendValidationErrors(msg.NormalizedPayloadWarnings, normalizedMeasurements)
	} else {
		// If the normalizer is not set, the decoder may return already normalized payload.
		// This is a best effort attempt to parse the decoded payload as normalized payload.
		// If that does not return an error, the decoded payload is assumed to be normalized.
		normalizedPayload := []*structpb.Struct{
			decodedPayload,
		}
		normalizedMeasurements, err := normalizedpayload.Parse(normalizedPayload)
		if err == nil {
			msg.NormalizedPayload = make([]*structpb.Struct, 0, len(normalizedMeasurements))
			for _, measurement := range normalizedMeasurements {
				if len(measurement.Valid.GetFields()) == 0 {
					continue
				}
				msg.NormalizedPayload = append(msg.NormalizedPayload, measurement.Valid)
			}
			msg.NormalizedPayloadWarnings = appendValidationErrors(msg.NormalizedPayloadWarnings, normalizedMeasurements)
		}
	}

	return nil
}

type decodeDownlinkInput struct {
	Bytes []uint8 `json:"bytes"`
	FPort uint8   `json:"fPort"`
}

type decodeDownlinkOutput struct {
	Data     map[string]any `json:"data"`
	Warnings []string       `json:"warnings"`
	Errors   []string       `json:"errors"`
}

func wrapDownlinkDecoderScript(script string) string {
	return fmt.Sprintf(`
		%s

		function main(input) {
			const bytes = input.bytes.slice();
			const { fPort } = input;
			return decodeDownlink({ bytes, fPort });
		}
	`, script)
}

// CompileDownlinkDecoder generates a downlink decoder from the provided script.
func (h *host) CompileDownlinkDecoder(
	ctx context.Context, script string,
) (
	func(
		context.Context,
		*ttnpb.EndDeviceIdentifiers,
		*ttnpb.EndDeviceVersionIdentifiers,
		*ttnpb.ApplicationDownlink,
	) error,
	error,
) {
	defer trace.StartRegion(ctx, "compile downlink decoder").End()

	run, err := h.engine.Compile(ctx, wrapDownlinkDecoderScript(script))
	if err != nil {
		return nil, err
	}

	return func(
		ctx context.Context,
		_ *ttnpb.EndDeviceIdentifiers,
		_ *ttnpb.EndDeviceVersionIdentifiers,
		msg *ttnpb.ApplicationDownlink,
	) error {
		return h.decodeDownlink(ctx, msg, run)
	}, nil
}

// DecodeUplink decodes the message's FRMPayload to DecodedPayload using the given script.
func (h *host) DecodeDownlink(
	ctx context.Context,
	_ *ttnpb.EndDeviceIdentifiers,
	_ *ttnpb.EndDeviceVersionIdentifiers,
	msg *ttnpb.ApplicationDownlink,
	script string,
) error {
	run := func(ctx context.Context, fn string, params ...any) (func(any) error, error) {
		return h.engine.Run(ctx, wrapDownlinkDecoderScript(script), fn, params...)
	}
	return h.decodeDownlink(ctx, msg, run)
}

func (*host) decodeDownlink(
	ctx context.Context,
	msg *ttnpb.ApplicationDownlink,
	run func(context.Context, string, ...any) (func(any) error, error),
) error {
	defer trace.StartRegion(ctx, "decode downlink message").End()

	input := decodeDownlinkInput{
		Bytes: msg.FrmPayload,
		FPort: uint8(msg.FPort), // nolint:gosec
	}

	valueAs, err := run(ctx, "main", input)
	if err != nil {
		return err
	}

	var output decodeDownlinkOutput
	err = valueAs(&output)
	if err != nil {
		return errOutput.WithCause(err)
	}
	if len(output.Errors) > 0 {
		return errOutputErrors.WithAttributes("errors", strings.Join(output.Errors, ", "))
	}

	s, err := goproto.Struct(output.Data)
	if err != nil {
		return errOutput.WithCause(err)
	}
	if errs := goproto.ValidateStruct(s); len(errs) > 0 {
		return errOutputEncoding.WithAttributes("errors", strings.Join(errs, ", "))
	}

	msg.DecodedPayload = s
	msg.DecodedPayloadWarnings = output.Warnings

	return nil
}
