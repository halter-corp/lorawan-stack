// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
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

package rpclog

import (
	"context"
	"time"

	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryClientInterceptor returns a new unary client interceptor that optionally logs the execution of external gRPC calls.
func UnaryClientInterceptor(ctx context.Context, opts ...Option) grpc.UnaryClientInterceptor {
	o := evaluateClientOpt(opts)
	logger := log.FromContext(ctx).WithField("namespace", "grpc")
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		onceFields, propagatedFields := logFieldsForCall(ctx, method)
		logger := logger.WithFields(propagatedFields)
		newCtx := log.NewContext(ctx, logger)

		var md metadata.MD
		startTime := time.Now()
		err := invoker(newCtx, method, req, reply, cc, append(opts, grpc.Header(&md))...)
		if requestID := md.Get("x-request-id"); len(requestID) > 0 {
			onceFields = onceFields.WithField(
				"request_id", requestID[0],
			)
		}
		onceFields = onceFields.WithField(
			"duration", time.Since(startTime).Round(time.Microsecond*100),
		)
		if err != nil {
			onceFields = onceFields.WithFields(logFieldsForError(err))
		}

		level := o.levelFunc(o.codeFunc(err))
		if level > log.InfoLevel {
			level = log.InfoLevel
		}
		entry := logger.WithFields(onceFields)
		if err != nil {
			entry = entry.WithError(err)
		}
		commit(entry, level, "Finished unary call")
		return err
	}
}

// StreamClientInterceptor returns a new streaming client interceptor that optionally logs the execution of external gRPC calls.
func StreamClientInterceptor(ctx context.Context, opts ...Option) grpc.StreamClientInterceptor {
	o := evaluateClientOpt(opts)
	logger := log.FromContext(ctx).WithField("namespace", "grpc")
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		onceFields, propagatedFields := logFieldsForCall(ctx, method)
		logger := logger.WithFields(propagatedFields)
		newCtx := log.NewContext(ctx, logger)

		var md metadata.MD
		startTime := time.Now()
		clientStream, err := streamer(newCtx, desc, cc, method, append(opts, grpc.Header(&md))...)
		if requestID := md.Get("x-request-id"); len(requestID) > 0 {
			onceFields = onceFields.WithField(
				"request_id", requestID[0],
			)
		}
		if err != nil {
			onceFields = onceFields.WithField(
				"duration", time.Since(startTime).Round(time.Microsecond*100),
			)
			onceFields = onceFields.WithFields(logFieldsForError(err))
			level := o.levelFunc(o.codeFunc(err))
			entry := logger.WithFields(onceFields)
			if err != nil {
				entry = entry.WithError(err)
			}
			commit(entry, level, "Failed streaming call")
			return clientStream, err
		}
		go func() {
			<-clientStream.Context().Done()
			err := clientStream.Context().Err()
			onceFields = onceFields.WithField(
				"duration", time.Since(startTime).Round(time.Microsecond*100),
			)
			if err != nil {
				onceFields = onceFields.WithFields(logFieldsForError(err))
			}
			level := o.levelFunc(o.codeFunc(err))
			entry := logger.WithFields(onceFields)
			if err != nil {
				entry = entry.WithError(err)
			}
			commit(entry, level, "Finished streaming call")
		}()
		return clientStream, err
	}
}
