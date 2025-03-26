// Copyright © 2023 The Things Network Foundation, The Things Industries B.V.
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

package mac

import (
	"context"
	"fmt"

	"go.thethings.network/lorawan-stack/v3/pkg/band"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
)

// RelayKeyService provides common relay related cryptographic operations.
type RelayKeyService interface {
	// BatchDeriveRootWorSKey derives the RootWorSKey for the provided end devices.
	// For devices with a pending session, the derived RootWorSKey is derived rom the
	// pending NwkSEncKey. For devices with an active session, the derived RootWorSKey
	// is derived from the active NwkSEncKey.
	BatchDeriveRootWorSKey(
		ctx context.Context, appID *ttnpb.ApplicationIdentifiers, deviceIDs []string, sessionKeyIDs [][]byte,
	) (devAddrs []*types.DevAddr, keys []*types.AES128Key, err error)
}

func secondChFields(secondCh *ttnpb.RelaySecondChannel) []any {
	if secondCh == nil {
		return nil
	}
	return []any{
		"relay_second_ch_ack_offset", secondCh.AckOffset,
		"relay_second_ch_data_rate_index", secondCh.DataRateIndex,
		"relay_second_ch_frequency", secondCh.Frequency,
	}
}

func servingRelayFields(serving *ttnpb.ServingRelayParameters) log.Fielder {
	if serving == nil {
		return log.Fields()
	}
	return log.Fields(
		append(
			secondChFields(serving.SecondChannel),
			"relay_default_ch_index", serving.DefaultChannelIndex,
			"relay_cad_periodicity", serving.CadPeriodicity,
		)...,
	)
}

func relayForwardLimitsFields(limits *ttnpb.RelayForwardLimits, prefix string) []any {
	if limits == nil {
		return nil
	}
	return []any{
		fmt.Sprintf("relay_%v_limit_bucket_size", prefix), limits.BucketSize,
		fmt.Sprintf("relay_%v_limit_reload_rate", prefix), limits.ReloadRate,
	}
}

func relayConfigureForwardLimitsFields(limits *ttnpb.ServingRelayForwardingLimits) log.Fielder {
	if limits == nil {
		return log.Fields()
	}
	fields := []any{"relay_limit_reset_behavior", limits.ResetBehavior}
	fields = append(fields, relayForwardLimitsFields(limits.JoinRequests, "join_requests")...)
	fields = append(fields, relayForwardLimitsFields(limits.Notifications, "notifications")...)
	fields = append(fields, relayForwardLimitsFields(limits.UplinkMessages, "uplink_messages")...)
	fields = append(fields, relayForwardLimitsFields(limits.Overall, "overall")...)
	return log.Fields(fields...)
}

func servedRelayFields(served *ttnpb.ServedRelayParameters) log.Fielder {
	if served == nil {
		return log.Fields()
	}
	fields := []any{}
	switch {
	case served.GetAlways() != nil:
		fields = append(fields, "relay_mode", "always")
	case served.GetDynamic() != nil:
		fields = append(
			fields,
			"relay_mode", "dynamic",
			"relay_smart_enable_level", served.GetDynamic().SmartEnableLevel,
		)
	case served.GetEndDeviceControlled() != nil:
		fields = append(fields, "relay_mode", "end_device_controlled")
	default:
		panic("unreachable")
	}
	fields = append(fields, "relay_backoff", served.Backoff)
	fields = append(fields, secondChFields(served.SecondChannel)...)
	return log.Fields(fields...)
}

func relayUpdateUplinkListReqFields(req *ttnpb.MACCommand_RelayUpdateUplinkListReq) log.Fielder {
	fields := []any{
		"relay_rule_index", req.RuleIndex,
		"relay_served_dev_addr", types.MustDevAddr(req.DevAddr),
		"relay_served_w_f_cnt", req.WFCnt,
		"relay_served_session_key_id", req.SessionKeyId,
	}
	if limits := req.ForwardLimits; limits != nil {
		fields = append(fields,
			"relay_served_bucket_size", limits.BucketSize,
			"relay_served_reload_rate", limits.ReloadRate,
		)
	}
	return log.Fields(fields...)
}

func relayCtrlUplinkListReqFields(req *ttnpb.MACCommand_RelayCtrlUplinkListReq) log.Fielder {
	return log.Fields(
		"relay_rule_index", req.RuleIndex,
		"relay_ctrl_action", req.Action,
	)
}

// DeviceDefaultRelaySettings returns the default relay parameters for the given device.
func DeviceDefaultRelaySettings(
	dev *ttnpb.EndDevice,
	defaults *ttnpb.MACSettings,
	profile *ttnpb.MACSettings,
) *ttnpb.RelaySettings {
	switch {
	case profile.GetRelay() != nil:
		return profile.Relay
	case dev.GetMacSettings().GetRelay() != nil:
		return dev.MacSettings.Relay
	case defaults.GetRelay() != nil:
		return defaults.Relay
	default:
		return nil
	}
}

// DeviceDesiredRelaySettings returns the desired relay parameters for the given device.
func DeviceDesiredRelaySettings(
	dev *ttnpb.EndDevice,
	defaults *ttnpb.MACSettings,
	profile *ttnpb.MACSettings,
) *ttnpb.RelaySettings {
	switch {
	case profile.GetDesiredRelay() != nil:
		return profile.DesiredRelay
	case dev.GetMacSettings().GetDesiredRelay() != nil:
		return dev.MacSettings.DesiredRelay
	case defaults.GetDesiredRelay() != nil:
		return defaults.DesiredRelay
	default:
		return DeviceDefaultRelaySettings(dev, defaults, profile)
	}
}

// ServedRelayParametersFromServedRelaySettings returns the served relay parameters for the given settings.
func ServedRelayParametersFromServedRelaySettings(
	settings *ttnpb.ServedRelaySettings, phy *band.Band,
) *ttnpb.ServedRelayParameters {
	if settings == nil {
		return nil
	}
	parameters := &ttnpb.ServedRelayParameters{
		Backoff:         phy.ServedRelayBackoff,
		SecondChannel:   settings.SecondChannel,
		ServingDeviceId: settings.ServingDeviceId,
	}
	if backoff := settings.Backoff; backoff != nil {
		parameters.Backoff = backoff.Value
	}
	switch {
	case settings.GetAlways() != nil:
		parameters.Mode = &ttnpb.ServedRelayParameters_Always{
			Always: settings.GetAlways(),
		}
	case settings.GetDynamic() != nil:
		parameters.Mode = &ttnpb.ServedRelayParameters_Dynamic{
			Dynamic: settings.GetDynamic(),
		}
	case settings.GetEndDeviceControlled() != nil:
		parameters.Mode = &ttnpb.ServedRelayParameters_EndDeviceControlled{
			EndDeviceControlled: settings.GetEndDeviceControlled(),
		}
	default:
		panic("unreachable")
	}
	return parameters
}

// ServingRelayParametersFromServingRelaySettings returns the serving relay parameters for the given settings.
func ServingRelayParametersFromServingRelaySettings(
	settings *ttnpb.ServingRelaySettings,
) *ttnpb.ServingRelayParameters {
	if settings == nil {
		return nil
	}
	parameters := &ttnpb.ServingRelayParameters{
		SecondChannel:         settings.SecondChannel,
		CadPeriodicity:        settings.CadPeriodicity,
		UplinkForwardingRules: settings.UplinkForwardingRules,
		Limits:                settings.Limits,
	}
	if defaultChIdx := settings.DefaultChannelIndex; defaultChIdx != nil {
		parameters.DefaultChannelIndex = defaultChIdx.Value
	}
	return parameters
}

// RelayParametersFromRelaySettings returns the relay parameters for the given settings.
func RelayParametersFromRelaySettings(settings *ttnpb.RelaySettings, phy *band.Band) *ttnpb.RelayParameters {
	if settings == nil {
		return nil
	}
	parameters := &ttnpb.RelayParameters{}
	switch {
	case settings.GetServed() != nil:
		parameters.Mode = &ttnpb.RelayParameters_Served{
			Served: ServedRelayParametersFromServedRelaySettings(settings.GetServed(), phy),
		}
	case settings.GetServing() != nil:
		parameters.Mode = &ttnpb.RelayParameters_Serving{
			Serving: ServingRelayParametersFromServingRelaySettings(settings.GetServing()),
		}
	default:
		panic("unreachable")
	}
	return parameters
}
