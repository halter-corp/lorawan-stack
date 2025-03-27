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

package networkserver_test

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smarty/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/v3/pkg/component"
	"go.thethings.network/lorawan-stack/v3/pkg/config"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	. "go.thethings.network/lorawan-stack/v3/pkg/networkserver"
	. "go.thethings.network/lorawan-stack/v3/pkg/networkserver/internal/test"
	"go.thethings.network/lorawan-stack/v3/pkg/networkserver/mac"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestDeviceRegistryGet(t *testing.T) {
	for _, tc := range []struct {
		Name           string
		ContextFunc    func(context.Context) context.Context
		GetByIDFunc    func(context.Context, *ttnpb.ApplicationIdentifiers, string, []string) (*ttnpb.EndDevice, context.Context, error)
		KeyVault       map[string][]byte
		Request        *ttnpb.GetEndDeviceRequest
		Device         *ttnpb.EndDevice
		ErrorAssertion func(*testing.T, error) bool
		GetByIDCalls   uint64
	}{
		{
			Name: "No device read rights",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, &rights.Rights{
					ApplicationRights: *rights.NewMap(map[string]*ttnpb.Rights{
						unique.ID(test.Context(), &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.Right_RIGHT_GATEWAY_SETTINGS_BASIC,
							},
						},
					}),
				})
			},
			GetByIDFunc: func(ctx context.Context, appID *ttnpb.ApplicationIdentifiers, devID string, gets []string) (*ttnpb.EndDevice, context.Context, error) {
				err := errors.New("GetByIDFunc must not be called")
				test.MustTFromContext(ctx).Error(err)
				return nil, ctx, err
			},
			Request: &ttnpb.GetEndDeviceRequest{
				EndDeviceIds: &ttnpb.EndDeviceIdentifiers{
					DeviceId:       "test-dev-id",
					ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"},
				},
				FieldMask: ttnpb.FieldMask("frequency_plan_id"),
			},
			ErrorAssertion: func(t *testing.T, err error) bool {
				if !assertions.New(t).So(errors.IsPermissionDenied(err), should.BeTrue) {
					t.Errorf("Received error: %s", err)
					return false
				}
				return true
			},
		},

		{
			Name: "no keys",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, &rights.Rights{
					ApplicationRights: *rights.NewMap(map[string]*ttnpb.Rights{
						unique.ID(test.Context(), &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.Right_RIGHT_APPLICATION_DEVICES_READ,
							},
						},
					}),
				})
			},
			GetByIDFunc: func(ctx context.Context, appID *ttnpb.ApplicationIdentifiers, devID string, gets []string) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.HaveSameElementsDeep, []string{
					"frequency_plan_id",
				})
				return &ttnpb.EndDevice{
					Ids: &ttnpb.EndDeviceIdentifiers{
						DeviceId:       "test-dev-id",
						ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"},
					},
					FrequencyPlanId: test.EUFrequencyPlanID,
				}, ctx, nil
			},
			Request: &ttnpb.GetEndDeviceRequest{
				EndDeviceIds: &ttnpb.EndDeviceIdentifiers{
					DeviceId:       "test-dev-id",
					ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"},
				},
				FieldMask: ttnpb.FieldMask("frequency_plan_id"),
			},
			Device: &ttnpb.EndDevice{
				Ids: &ttnpb.EndDeviceIdentifiers{
					DeviceId:       "test-dev-id",
					ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"},
				},
				FrequencyPlanId: test.EUFrequencyPlanID,
			},
			GetByIDCalls: 1,
		},

		{
			Name: "with keys",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, &rights.Rights{
					ApplicationRights: *rights.NewMap(map[string]*ttnpb.Rights{
						unique.ID(test.Context(), &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.Right_RIGHT_APPLICATION_DEVICES_READ,
								ttnpb.Right_RIGHT_APPLICATION_DEVICES_READ_KEYS,
								ttnpb.Right_RIGHT_APPLICATION_TRAFFIC_READ,
							},
						},
					}),
				})
			},
			GetByIDFunc: func(ctx context.Context, appID *ttnpb.ApplicationIdentifiers, devID string, gets []string) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.HaveSameElementsDeep, []string{
					"frequency_plan_id",
					"session",
					"queued_application_downlinks",
				})
				return &ttnpb.EndDevice{
					Ids: &ttnpb.EndDeviceIdentifiers{
						DeviceId:       "test-dev-id",
						ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"},
					},
					FrequencyPlanId: test.EUFrequencyPlanID,
					Session: &ttnpb.Session{
						Keys: &ttnpb.SessionKeys{
							FNwkSIntKey: &ttnpb.KeyEnvelope{
								KekLabel:     "test",
								EncryptedKey: []byte{0x96, 0x77, 0x8b, 0x25, 0xae, 0x6c, 0xa4, 0x35, 0xf9, 0x2b, 0x5b, 0x97, 0xc0, 0x50, 0xae, 0xd2, 0x46, 0x8a, 0xb8, 0xa1, 0x7a, 0xd8, 0x4e, 0x5d},
							},
						},
					},
				}, ctx, nil
			},
			KeyVault: map[string][]byte{
				"test": {0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17},
			},
			Request: &ttnpb.GetEndDeviceRequest{
				EndDeviceIds: &ttnpb.EndDeviceIdentifiers{
					DeviceId:       "test-dev-id",
					ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"},
				},
				FieldMask: ttnpb.FieldMask("frequency_plan_id", "session"),
			},
			Device: &ttnpb.EndDevice{
				Ids: &ttnpb.EndDeviceIdentifiers{
					DeviceId:       "test-dev-id",
					ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"},
				},
				FrequencyPlanId: test.EUFrequencyPlanID,
				Session: &ttnpb.Session{
					Keys: &ttnpb.SessionKeys{
						FNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: types.AES128Key{0x0, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}.Bytes(),
						},
					},
				},
			},
			GetByIDCalls: 1,
		},

		{
			Name: "with specific key envelope",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, &rights.Rights{
					ApplicationRights: *rights.NewMap(map[string]*ttnpb.Rights{
						unique.ID(test.Context(), &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.Right_RIGHT_APPLICATION_DEVICES_READ,
								ttnpb.Right_RIGHT_APPLICATION_DEVICES_READ_KEYS,
							},
						},
					}),
				})
			},
			GetByIDFunc: func(ctx context.Context, appID *ttnpb.ApplicationIdentifiers, devID string, gets []string) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.HaveSameElementsDeep, []string{
					"pending_session.keys.f_nwk_s_int_key",
				})
				return &ttnpb.EndDevice{
					Ids: &ttnpb.EndDeviceIdentifiers{
						DeviceId:       "test-dev-id",
						ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"},
					},
					PendingSession: &ttnpb.Session{
						Keys: &ttnpb.SessionKeys{
							FNwkSIntKey: &ttnpb.KeyEnvelope{
								KekLabel:     "test",
								EncryptedKey: []byte{0x96, 0x77, 0x8b, 0x25, 0xae, 0x6c, 0xa4, 0x35, 0xf9, 0x2b, 0x5b, 0x97, 0xc0, 0x50, 0xae, 0xd2, 0x46, 0x8a, 0xb8, 0xa1, 0x7a, 0xd8, 0x4e, 0x5d},
							},
						},
					},
				}, ctx, nil
			},
			KeyVault: map[string][]byte{
				"test": {0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17},
			},
			Request: &ttnpb.GetEndDeviceRequest{
				EndDeviceIds: &ttnpb.EndDeviceIdentifiers{
					DeviceId:       "test-dev-id",
					ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"},
				},
				FieldMask: ttnpb.FieldMask("pending_session.keys.f_nwk_s_int_key"),
			},
			Device: &ttnpb.EndDevice{
				Ids: &ttnpb.EndDeviceIdentifiers{
					DeviceId:       "test-dev-id",
					ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"},
				},
				PendingSession: &ttnpb.Session{
					Keys: &ttnpb.SessionKeys{
						FNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: types.AES128Key{0x0, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}.Bytes(),
						},
					},
				},
			},
			GetByIDCalls: 1,
		},

		{
			Name: "with specific key",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, &rights.Rights{
					ApplicationRights: *rights.NewMap(map[string]*ttnpb.Rights{
						unique.ID(test.Context(), &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.Right_RIGHT_APPLICATION_DEVICES_READ,
								ttnpb.Right_RIGHT_APPLICATION_DEVICES_READ_KEYS,
							},
						},
					}),
				})
			},
			GetByIDFunc: func(ctx context.Context, appID *ttnpb.ApplicationIdentifiers, devID string, gets []string) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.HaveSameElementsDeep, []string{
					"pending_session.keys.f_nwk_s_int_key.encrypted_key",
					"pending_session.keys.f_nwk_s_int_key.kek_label",
					"pending_session.keys.f_nwk_s_int_key.key",
				})
				return &ttnpb.EndDevice{
					Ids: &ttnpb.EndDeviceIdentifiers{
						DeviceId:       "test-dev-id",
						ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"},
					},
					PendingSession: &ttnpb.Session{
						Keys: &ttnpb.SessionKeys{
							FNwkSIntKey: &ttnpb.KeyEnvelope{
								KekLabel:     "test",
								EncryptedKey: []byte{0x96, 0x77, 0x8b, 0x25, 0xae, 0x6c, 0xa4, 0x35, 0xf9, 0x2b, 0x5b, 0x97, 0xc0, 0x50, 0xae, 0xd2, 0x46, 0x8a, 0xb8, 0xa1, 0x7a, 0xd8, 0x4e, 0x5d},
							},
						},
					},
				}, ctx, nil
			},
			KeyVault: map[string][]byte{
				"test": {0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17},
			},
			Request: &ttnpb.GetEndDeviceRequest{
				EndDeviceIds: &ttnpb.EndDeviceIdentifiers{
					DeviceId:       "test-dev-id",
					ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"},
				},
				FieldMask: ttnpb.FieldMask("pending_session.keys.f_nwk_s_int_key.key"),
			},
			Device: &ttnpb.EndDevice{
				Ids: &ttnpb.EndDeviceIdentifiers{
					DeviceId:       "test-dev-id",
					ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"},
				},
				PendingSession: &ttnpb.Session{
					Keys: &ttnpb.SessionKeys{
						FNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: types.AES128Key{0x0, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}.Bytes(),
						},
					},
				},
			},
			GetByIDCalls: 1,
		},
	} {
		tc := tc
		test.RunSubtest(t, test.SubtestConfig{
			Name:     tc.Name,
			Parallel: true,
			Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
				var getByIDCalls uint64

				ns, ctx, _, stop := StartTest(
					ctx,
					TestConfig{
						Component: component.Config{
							ServiceBase: config.ServiceBase{
								KeyVault: config.KeyVault{
									Provider: "static",
									Static:   tc.KeyVault,
								},
							},
						},
						NetworkServer: Config{
							Devices: &MockDeviceRegistry{
								GetByIDFunc: func(ctx context.Context, appID *ttnpb.ApplicationIdentifiers, devID string, gets []string) (*ttnpb.EndDevice, context.Context, error) {
									atomic.AddUint64(&getByIDCalls, 1)
									return tc.GetByIDFunc(ctx, appID, devID, gets)
								},
							},
						},
						TaskStarter: StartTaskExclude(
							DownlinkProcessTaskName,
							DownlinkDispatchTaskName,
						),
					},
				)
				defer stop()

				ns.AddContextFiller(tc.ContextFunc)
				ns.AddContextFiller(func(ctx context.Context) context.Context {
					return test.ContextWithTB(ctx, t)
				})

				req := ttnpb.Clone(tc.Request)
				dev, err := ttnpb.NewNsEndDeviceRegistryClient(ns.LoopbackConn()).Get(ctx, req)
				if tc.ErrorAssertion != nil && a.So(tc.ErrorAssertion(t, err), should.BeTrue) {
					a.So(dev, should.BeNil)
				} else if a.So(err, should.BeNil) {
					a.So(dev, should.Resemble, tc.Device)
				}
				a.So(req, should.Resemble, tc.Request)
				a.So(getByIDCalls, should.Equal, tc.GetByIDCalls)
			},
		})
	}
}

func TestDeviceRegistrySet(t *testing.T) {
	defaultMACSettings := test.Must(DefaultConfig.DefaultMACSettings.Parse())

	customMACSettings := test.Must(DefaultConfig.DefaultMACSettings.Parse())
	customMACSettings.Rx1Delay = &ttnpb.RxDelayValue{Value: ttnpb.RxDelay_RX_DELAY_2}
	customMACSettings.Rx1DataRateOffset = nil

	customMACSettingsOpt := EndDeviceOptions.WithMacSettings(customMACSettings)

	multicastClassBMACSettings := test.Must(DefaultConfig.DefaultMACSettings.Parse())
	multicastClassBMACSettings.PingSlotPeriodicity = &ttnpb.PingSlotPeriodValue{
		Value: ttnpb.PingSlotPeriod_PING_EVERY_16S,
	}

	multicastClassBMACSettingsOpt := EndDeviceOptions.WithMacSettings(multicastClassBMACSettings)

	currentMACStateOverrideOpt := func(macState *ttnpb.MACState) *ttnpb.MACState {
		macState = ttnpb.Clone(macState)
		macState.CurrentParameters.Rx1Delay = ttnpb.RxDelay_RX_DELAY_3
		macState.CurrentParameters.Rx1DataRateOffset = ttnpb.DataRateOffset_DATA_RATE_OFFSET_1
		return macState
	}
	desiredMACStateOverrideOpt := func(macState *ttnpb.MACState) *ttnpb.MACState {
		macState = ttnpb.Clone(macState)
		macState.DesiredParameters.Rx1Delay = ttnpb.RxDelay_RX_DELAY_4
		macState.DesiredParameters.Rx1DataRateOffset = ttnpb.DataRateOffset_DATA_RATE_OFFSET_2
		return macState
	}
	activeMACStateOpts := []test.MACStateOption{
		currentMACStateOverrideOpt,
		desiredMACStateOverrideOpt,
	}

	activeSessionOpts := []test.SessionOption{
		SessionOptions.WithLastNFCntDown(0x24),
	}
	activeSessionOptsWithStartedAt := append(activeSessionOpts,
		SessionOptions.WithStartedAt(timestamppb.New(time.Unix(0, 42))),
	)

	activateOpt := EndDeviceOptions.Activate(customMACSettings, false, activeSessionOpts, activeMACStateOpts...)

	macStateWithoutRX1DelayOpt := func(dev *ttnpb.EndDevice) *ttnpb.EndDevice {
		dev = ttnpb.Clone(dev)
		dev.MacState.CurrentParameters.Rx1Delay = 0
		return dev
	}

	makeUpdateDeviceRequest := func(deviceOpts []test.EndDeviceOption, paths ...string) *SetDeviceRequest {
		return &SetDeviceRequest{
			EndDevice: test.MakeEndDevice(deviceOpts...),
			Paths:     paths,
		}
	}

	macSettingsProfileID := &ttnpb.MACSettingsProfileIdentifiers{
		ApplicationIds: &ttnpb.ApplicationIdentifiers{
			ApplicationId: "test-app-id",
		},
		ProfileId: "test-mac-settings-profile-id",
	}

	macSettingsProfileOpt := EndDeviceOptions.WithMacSettingsProfileIds(macSettingsProfileID)
	emptyMacSettingsProfileOpt := EndDeviceOptions.WithMacSettingsProfileIds(nil)

	for createDevice, tcs := range map[*ttnpb.EndDevice][]struct {
		SetDevice      SetDeviceRequest
		RequiredRights []ttnpb.Right

		ReturnedDevice *ttnpb.EndDevice
		StoredDevice   *ttnpb.EndDevice
	}{
		nil: {
			// OTAA Create
			{
				SetDevice: *MakeOTAASetDeviceRequest(nil),

				ReturnedDevice: MakeOTAAEndDevice(),
				StoredDevice:   MakeOTAAEndDevice(),
			},
			{
				SetDevice: *MakeOTAASetDeviceRequest([]test.EndDeviceOption{
					EndDeviceOptions.SendJoinRequest(customMACSettings, false),
				},
					"pending_mac_state",
				),
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE_KEYS,
				},

				ReturnedDevice: MakeOTAAEndDevice(
					EndDeviceOptions.SendJoinRequest(customMACSettings, false),
				),
				StoredDevice: MakeOTAAEndDevice(
					EndDeviceOptions.SendJoinRequest(customMACSettings, true),
				),
			},
			{
				SetDevice: *MakeOTAASetDeviceRequest([]test.EndDeviceOption{
					EndDeviceOptions.SendJoinRequest(customMACSettings, false),
					EndDeviceOptions.SendJoinAccept(ttnpb.TxSchedulePriority_HIGHEST),
				},
					"pending_mac_state",
					"pending_session",
				),
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE_KEYS,
				},

				ReturnedDevice: MakeOTAAEndDevice(
					EndDeviceOptions.SendJoinRequest(customMACSettings, false),
					EndDeviceOptions.SendJoinAccept(ttnpb.TxSchedulePriority_HIGHEST),
				),
				StoredDevice: MakeOTAAEndDevice(
					EndDeviceOptions.SendJoinRequest(customMACSettings, true),
					EndDeviceOptions.SendJoinAccept(ttnpb.TxSchedulePriority_HIGHEST),
				),
			},
			{
				SetDevice: *MakeOTAASetDeviceRequest([]test.EndDeviceOption{
					customMACSettingsOpt,
				},
					"mac_settings",
				),

				ReturnedDevice: MakeOTAAEndDevice(
					customMACSettingsOpt,
				),
				StoredDevice: MakeOTAAEndDevice(
					customMACSettingsOpt,
				),
			},
			{
				SetDevice: *MakeOTAASetDeviceRequest([]test.EndDeviceOption{
					macSettingsProfileOpt,
				},
					"mac_settings_profile_ids",
				),

				ReturnedDevice: MakeOTAAEndDevice(
					macSettingsProfileOpt,
				),
				StoredDevice: MakeOTAAEndDevice(
					macSettingsProfileOpt,
				),
			},
			{
				SetDevice: *MakeOTAASetDeviceRequest([]test.EndDeviceOption{
					activateOpt,
				},
					"mac_state",
					"session",
				),
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE_KEYS,
				},

				ReturnedDevice: MakeOTAAEndDevice(
					EndDeviceOptions.Activate(customMACSettings, false, activeSessionOpts, activeMACStateOpts...),
				),
				StoredDevice: MakeOTAAEndDevice(
					EndDeviceOptions.Activate(customMACSettings, true, activeSessionOpts, activeMACStateOpts...),
				),
			},
			{
				SetDevice: *MakeOTAASetDeviceRequest([]test.EndDeviceOption{
					activateOpt,
				},
					"mac_state.current_parameters",
					"mac_state.lorawan_version",
					"session",
				),
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE_KEYS,
				},

				ReturnedDevice: MakeOTAAEndDevice(
					EndDeviceOptions.Activate(customMACSettings, false, activeSessionOpts, currentMACStateOverrideOpt),
				),
				StoredDevice: MakeOTAAEndDevice(
					EndDeviceOptions.Activate(customMACSettings, true, activeSessionOpts, currentMACStateOverrideOpt),
					EndDeviceOptions.WithMACStateOptions(
						MACStateOptions.WithRecentUplinks(),
						MACStateOptions.WithRecentDownlinks(),
					),
				),
			},
			{
				SetDevice: *MakeOTAASetDeviceRequest([]test.EndDeviceOption{
					activateOpt,
				},
					"mac_state.desired_parameters",
					"mac_state.lorawan_version",
					"session.dev_addr",
					"session.keys.f_nwk_s_int_key.key",
					"session.keys.nwk_s_enc_key.key",
					"session.keys.s_nwk_s_int_key.key",
					"session.keys.session_key_id",
				),
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE_KEYS,
				},

				ReturnedDevice: MakeOTAAEndDevice(
					EndDeviceOptions.Activate(defaultMACSettings, false, nil, desiredMACStateOverrideOpt),
				),
				StoredDevice: MakeOTAAEndDevice(
					EndDeviceOptions.Activate(defaultMACSettings, true, nil, desiredMACStateOverrideOpt),
					EndDeviceOptions.WithMACStateOptions(
						MACStateOptions.WithRecentUplinks(),
						MACStateOptions.WithRecentDownlinks(),
					),
				),
			},

			// OTAA Create 1.0.3
			{
				SetDevice: *MakeOTAASetDeviceRequest([]test.EndDeviceOption{
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					activateOpt,
				},
					"mac_state",
					"session",
				),
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE_KEYS,
				},

				ReturnedDevice: MakeOTAAEndDevice(
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					activateOpt,
				),
				StoredDevice: MakeOTAAEndDevice(
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					EndDeviceOptions.Activate(customMACSettings, true, activeSessionOpts, activeMACStateOpts...),
				),
			},
			{
				SetDevice: *MakeOTAASetDeviceRequest([]test.EndDeviceOption{
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					EndDeviceOptions.SendJoinRequest(customMACSettings, false),
				},
					"pending_mac_state",
				),
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE_KEYS,
				},

				ReturnedDevice: MakeOTAAEndDevice(
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					EndDeviceOptions.SendJoinRequest(customMACSettings, false),
				),
				StoredDevice: MakeOTAAEndDevice(
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					EndDeviceOptions.SendJoinRequest(customMACSettings, true),
				),
			},
			{
				SetDevice: *MakeOTAASetDeviceRequest([]test.EndDeviceOption{
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					EndDeviceOptions.SendJoinRequest(customMACSettings, false),
					EndDeviceOptions.SendJoinAccept(ttnpb.TxSchedulePriority_HIGHEST),
				},
					"pending_mac_state",
					"pending_session",
				),
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE_KEYS,
				},

				ReturnedDevice: MakeOTAAEndDevice(
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					EndDeviceOptions.SendJoinRequest(customMACSettings, false),
					EndDeviceOptions.SendJoinAccept(ttnpb.TxSchedulePriority_HIGHEST),
				),
				StoredDevice: MakeOTAAEndDevice(
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					EndDeviceOptions.SendJoinRequest(customMACSettings, true),
					EndDeviceOptions.SendJoinAccept(ttnpb.TxSchedulePriority_HIGHEST),
				),
			},
			{
				SetDevice: *MakeOTAASetDeviceRequest([]test.EndDeviceOption{
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					activateOpt,
				},
					"mac_state.current_parameters",
					"mac_state.lorawan_version",
					"session",
				),
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE_KEYS,
				},

				ReturnedDevice: MakeOTAAEndDevice(
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					EndDeviceOptions.Activate(customMACSettings, false, activeSessionOpts, currentMACStateOverrideOpt),
				),
				StoredDevice: MakeOTAAEndDevice(
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					EndDeviceOptions.Activate(customMACSettings, true, activeSessionOpts, currentMACStateOverrideOpt),
					EndDeviceOptions.WithMACStateOptions(
						MACStateOptions.WithRecentUplinks(),
						MACStateOptions.WithRecentDownlinks(),
					),
				),
			},
			{
				SetDevice: *MakeOTAASetDeviceRequest([]test.EndDeviceOption{
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					activateOpt,
				},
					"mac_state.desired_parameters",
					"mac_state.lorawan_version",
					"session.dev_addr",
					"session.keys.f_nwk_s_int_key.key",
					"session.keys.session_key_id",
				),
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE_KEYS,
				},

				ReturnedDevice: MakeOTAAEndDevice(
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					EndDeviceOptions.Activate(defaultMACSettings, false, nil, desiredMACStateOverrideOpt),
				),
				StoredDevice: MakeOTAAEndDevice(
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					EndDeviceOptions.Activate(defaultMACSettings, true, nil, desiredMACStateOverrideOpt),
					EndDeviceOptions.WithMACStateOptions(
						MACStateOptions.WithRecentUplinks(),
						MACStateOptions.WithRecentDownlinks(),
					),
				),
			},

			// ABP Create
			{
				SetDevice: *MakeABPSetDeviceRequest(customMACSettings, activeSessionOpts, nil, nil,
					"mac_state.current_parameters.rx1_delay",
					"session",
				),
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE_KEYS,
				},

				ReturnedDevice: MakeABPEndDevice(customMACSettings, false, activeSessionOpts, nil),
				StoredDevice:   MakeABPEndDevice(customMACSettings, true, activeSessionOpts, nil),
			},

			// Multicast Create
			{
				SetDevice: *MakeMulticastSetDeviceRequest(ttnpb.Class_CLASS_C, defaultMACSettings, activeSessionOpts, nil, nil,
					"session.last_n_f_cnt_down",
				),
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE_KEYS,
				},

				ReturnedDevice: MakeMulticastEndDevice(ttnpb.Class_CLASS_C, defaultMACSettings, false, activeSessionOpts, nil),
				StoredDevice:   MakeMulticastEndDevice(ttnpb.Class_CLASS_C, defaultMACSettings, true, activeSessionOpts, nil),
			},
		},

		// OTAA Update
		MakeOTAAEndDevice(): {
			{
				SetDevice: *makeUpdateDeviceRequest([]test.EndDeviceOption{
					customMACSettingsOpt,
				},
					"mac_settings",
				),

				ReturnedDevice: MakeOTAAEndDevice(
					customMACSettingsOpt,
				),
				StoredDevice: MakeOTAAEndDevice(
					customMACSettingsOpt,
				),
			},
		},

		// ABP Update
		MakeABPEndDevice(defaultMACSettings, true, activeSessionOptsWithStartedAt, nil): {
			{
				SetDevice: *makeUpdateDeviceRequest([]test.EndDeviceOption{
					customMACSettingsOpt,
				},
					"mac_settings",
				),

				ReturnedDevice: customMACSettingsOpt(MakeABPEndDevice(defaultMACSettings, false, activeSessionOptsWithStartedAt, nil)),
				StoredDevice:   customMACSettingsOpt(MakeABPEndDevice(defaultMACSettings, true, activeSessionOptsWithStartedAt, nil)),
			},

			{
				SetDevice: *makeUpdateDeviceRequest(nil,
					"mac_settings.rx2_data_rate_index",
					"mac_state.current_parameters.rx1_delay",
					"pending_mac_state",
				),
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE_KEYS, // `pending_mac_state` requires key write rights
				},

				ReturnedDevice: macStateWithoutRX1DelayOpt(MakeABPEndDevice(defaultMACSettings, false, activeSessionOptsWithStartedAt, nil)),
				StoredDevice:   macStateWithoutRX1DelayOpt(MakeABPEndDevice(defaultMACSettings, true, activeSessionOptsWithStartedAt, nil)),
			},
		},

		// Multicast Update
		MakeMulticastEndDevice(ttnpb.Class_CLASS_B, defaultMACSettings, true, activeSessionOptsWithStartedAt, nil): {
			{
				SetDevice: *makeUpdateDeviceRequest([]test.EndDeviceOption{
					multicastClassBMACSettingsOpt,
				},
					"mac_settings",
				),

				ReturnedDevice: multicastClassBMACSettingsOpt(MakeMulticastEndDevice(ttnpb.Class_CLASS_B, defaultMACSettings, false, activeSessionOptsWithStartedAt, nil)),
				StoredDevice:   multicastClassBMACSettingsOpt(MakeMulticastEndDevice(ttnpb.Class_CLASS_B, defaultMACSettings, true, activeSessionOptsWithStartedAt, nil)),
			},
		},
		// Update with MAC settings profile
		MakeOTAAEndDevice(): {
			{
				SetDevice: *makeUpdateDeviceRequest([]test.EndDeviceOption{
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					EndDeviceOptions.WithDefaultFrequencyPlanID(),
					macSettingsProfileOpt,
				},
					"frequency_plan_id",
					"lorawan_version",
					"lorawan_phy_version",
					"mac_settings_profile_ids",
				),

				ReturnedDevice: MakeOTAAEndDevice(
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					EndDeviceOptions.WithDefaultFrequencyPlanID(),
					macSettingsProfileOpt,
				),
				StoredDevice: MakeOTAAEndDevice(
					EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
					EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
					EndDeviceOptions.WithDefaultFrequencyPlanID(),
					macSettingsProfileOpt,
				),
			},
		},
		// Update with empty MAC settings profile
		MakeOTAAEndDevice(macSettingsProfileOpt): {
			{
				SetDevice: *makeUpdateDeviceRequest([]test.EndDeviceOption{
					emptyMacSettingsProfileOpt,
				},
					"mac_settings_profile_ids",
				),

				ReturnedDevice: MakeOTAAEndDevice(),
				StoredDevice: MakeOTAAEndDevice(
					customMACSettingsOpt,
				),
			},
		},
	} {
		for _, tc := range tcs {
			createDevice := createDevice
			tc := tc
			test.RunSubtest(t, test.SubtestConfig{
				Name: MakeTestCaseName(func() []string {
					dev := createDevice
					typ := "Update"
					if createDevice == nil {
						dev = tc.SetDevice.EndDevice
						typ = "Create"
					}
					return []string{
						typ,
						fmt.Sprintf("mode:%s", func() string {
							switch {
							case dev.SupportsJoin:
								return "OTAA"
							case dev.Multicast:
								return "Multicast"
							default:
								return "ABP"
							}
						}()),
						fmt.Sprintf("MAC:%s", dev.LorawanVersion.String()),
						fmt.Sprintf("PHY:%s", dev.LorawanPhyVersion.String()),
						fmt.Sprintf("fp:%s", dev.FrequencyPlanId),
						fmt.Sprintf("paths:[%s]", strings.Join(tc.SetDevice.Paths, ",")),
					}
				}()...),
				Parallel: true,
				Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
					nsConf := DefaultConfig
					nsConf.DeviceKEKLabel = test.DefaultKEKLabel

					_, ctx, env, stop := StartTest(ctx, TestConfig{
						Component: component.Config{
							ServiceBase: config.ServiceBase{
								GRPC: config.GRPC{
									LogIgnoreMethods: []string{
										"/ttn.lorawan.v3.ApplicationAccess/ListRights",
										"/ttn.lorawan.v3.NsEndDeviceRegistry/Set",
									},
								},
								KeyVault: test.DefaultKeyVault,
								FrequencyPlans: config.FrequencyPlansConfig{
									ConfigSource: "static",
									Static:       test.StaticFrequencyPlans,
								},
							},
						},
						NetworkServer: nsConf,
						TaskStarter: StartTaskExclude(
							DownlinkProcessTaskName,
							DownlinkDispatchTaskName,
						),
					})
					defer stop()

					clock := test.NewMockClock(time.Now())
					defer SetMockClock(clock)()

					withCreatedAt := test.EndDeviceOptions.WithCreatedAt(timestamppb.New(clock.Now()))
					if createDevice != nil {
						_, ctx = MustCreateDevice(ctx, env.Devices, createDevice)
						clock.Add(time.Nanosecond)
					}

					macSettingsProfile := &ttnpb.MACSettingsProfile{
						Ids:         macSettingsProfileID,
						MacSettings: customMACSettings,
					}
					profile, err := env.MACSettingsProfileRegistry.Set(ctx, macSettingsProfileID, []string{"ids", "mac_settings"},
						func(context.Context, *ttnpb.MACSettingsProfile) (*ttnpb.MACSettingsProfile, []string, error) {
							return macSettingsProfile, []string{"ids", "mac_settings"}, nil
						})
					a.So(err, should.BeNil)
					a.So(profile, should.Resemble, macSettingsProfile)

					now := clock.Now()
					withTimestamps := withCreatedAt.Compose(
						test.EndDeviceOptions.WithUpdatedAt(timestamppb.New(now)),
						func(dev *ttnpb.EndDevice) *ttnpb.EndDevice {
							dev = ttnpb.Clone(dev)
							if dev.Session != nil && dev.Session.StartedAt == nil {
								dev.Session.StartedAt = timestamppb.New(now)
							}
							return dev
						},
					)

					req := &ttnpb.SetEndDeviceRequest{
						EndDevice: tc.SetDevice.EndDevice,
						FieldMask: ttnpb.FieldMask(tc.SetDevice.Paths...),
					}

					dev, err, ok := env.AssertSetDevice(ctx, createDevice == nil, req)
					if !a.So(ok, should.BeTrue) || !a.So(err, should.BeError) || !a.So(errors.IsPermissionDenied(err), should.BeTrue) {
						if err != nil {
							t.Errorf("Expected 'permission denied' error, got: %s", test.FormatError(err))
						}
						return
					}
					a.So(dev, should.BeNil)
					if len(tc.RequiredRights) > 0 {
						dev, err, ok = env.AssertSetDevice(ctx, createDevice == nil, req,
							ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE,
						)
						if !a.So(ok, should.BeTrue) || !a.So(err, should.BeError) || !a.So(errors.IsPermissionDenied(err), should.BeTrue) {
							if err != nil {
								t.Errorf("Expected 'permission denied' error, got: %s", test.FormatError(err))
							}
							return
						}
						a.So(dev, should.BeNil)
					}

					rights := append([]ttnpb.Right{
						ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE,
					}, tc.RequiredRights...)
					expectedReturn := test.Must(ttnpb.ApplyEndDeviceFieldMask(nil, withTimestamps(tc.ReturnedDevice), ttnpb.AddImplicitEndDeviceGetFields(tc.SetDevice.Paths...)...))

					dev, err, ok = env.AssertSetDevice(ctx, createDevice == nil, req, rights...)
					if !a.So(ok, should.BeTrue) || !a.So(err, should.BeNil) || !a.So(dev, should.NotBeNil) {
						if err != nil {
							t.Errorf("Expected no error, got: %s", test.FormatError(err))
						}
						return
					}
					a.So(dev, should.Resemble, expectedReturn)

					dev, _, err = env.Devices.GetByID(ctx, tc.SetDevice.Ids.ApplicationIds, tc.SetDevice.Ids.DeviceId, ttnpb.EndDeviceFieldPathsTopLevel)
					if !a.So(err, should.BeNil) || !a.So(dev, should.NotBeNil) {
						if err != nil {
							t.Errorf("Expected no error, got: %s", test.FormatError(err))
						}
						return
					}
					a.So(dev, should.Resemble, withTimestamps(tc.StoredDevice))

					now = clock.Add(time.Nanosecond)
					dev, err, ok = env.AssertSetDevice(ctx, false, &ttnpb.SetEndDeviceRequest{
						EndDevice: expectedReturn,
						FieldMask: ttnpb.FieldMask(tc.SetDevice.Paths...),
					}, rights...)
					if !a.So(ok, should.BeTrue) || !a.So(err, should.BeNil) || !a.So(dev, should.NotBeNil) {
						return
					}
					a.So(dev, should.Resemble, EndDeviceOptions.WithUpdatedAt(timestamppb.New(now))(expectedReturn))
				},
			})
		}
	}
}

func TestDeviceRegistryResetFactoryDefaults(t *testing.T) {
	activeSessionOpts := []test.SessionOption{
		SessionOptions.WithLastFCntUp(0x42),
		SessionOptions.WithLastNFCntDown(0x24),
		SessionOptions.WithDefaultQueuedApplicationDownlinks(),
	}
	macSettings := test.Must(DefaultConfig.DefaultMACSettings.Parse())
	activateOpt := EndDeviceOptions.Activate(macSettings, true, activeSessionOpts)

	// TODO: Refactor into same structure as Set
	for _, tc := range []struct {
		CreateDevice *SetDeviceRequest
	}{
		{},

		{
			CreateDevice: MakeOTAASetDeviceRequest(nil),
		},
		{
			CreateDevice: MakeOTAASetDeviceRequest([]test.EndDeviceOption{
				activateOpt,
			},
				"mac_state",
				"session",
			),
		},
		{
			CreateDevice: MakeOTAASetDeviceRequest([]test.EndDeviceOption{
				EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
				EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
				activateOpt,
			},
				"mac_state",
				"session",
			),
		},

		{
			CreateDevice: MakeABPSetDeviceRequest(macSettings, nil, nil, nil),
		},
		{
			CreateDevice: MakeABPSetDeviceRequest(macSettings, activeSessionOpts, nil, nil),
		},
		{
			CreateDevice: MakeABPSetDeviceRequest(macSettings, activeSessionOpts, nil, []test.EndDeviceOption{
				EndDeviceOptions.WithLorawanVersion(ttnpb.MACVersion_MAC_V1_0_3),
				EndDeviceOptions.WithLorawanPhyVersion(ttnpb.PHYVersion_RP001_V1_0_3_REV_A),
			}),
		},
	} {
		for _, conf := range []struct {
			Paths          []string
			RequiredRights []ttnpb.Right
		}{
			{},
			{
				Paths: []string{
					"battery_percentage",
					"downlink_margin",
					"last_dev_status_received_at",
					"mac_state.current_parameters",
					"session.last_f_cnt_up",
				},
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_READ,
				},
			},
			{
				Paths: []string{
					"battery_percentage",
					"session.last_f_cnt_up",
				},
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_READ,
				},
			},
			{
				Paths: []string{
					"session.keys",
				},
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_READ,
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_READ_KEYS,
				},
			},
			{
				Paths: []string{
					"battery_percentage",
					"downlink_margin",
					"last_dev_status_received_at",
					"pending_mac_state",
					"pending_session",
				},
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_READ,
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_READ_KEYS,
					ttnpb.Right_RIGHT_APPLICATION_TRAFFIC_READ,
				},
			},
			{
				Paths: []string{
					"battery_percentage",
					"downlink_margin",
					"last_dev_status_received_at",
					"mac_state",
					"pending_mac_state",
					"pending_session",
					"session",
					"supports_join",
				},
				RequiredRights: []ttnpb.Right{
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_READ,
					ttnpb.Right_RIGHT_APPLICATION_DEVICES_READ_KEYS,
					ttnpb.Right_RIGHT_APPLICATION_TRAFFIC_READ,
				},
			},
		} {
			tc := tc
			conf := conf
			test.RunSubtest(t, test.SubtestConfig{
				Name: func() string {
					if tc.CreateDevice == nil {
						return "no device"
					}
					return MakeTestCaseName(
						fmt.Sprintf("paths:[%s]", strings.Join(conf.Paths, ",")),
						func() string {
							if tc.CreateDevice.EndDevice.SupportsJoin {
								return "OTAA"
							}
							if tc.CreateDevice.EndDevice.Session == nil {
								return MakeTestCaseName("ABP", "no session")
							}
							return fmt.Sprintf(MakeTestCaseName("ABP", "dev_addr:%s", "queue_len:%d", "session_keys:%v"),
								types.MustDevAddr(tc.CreateDevice.Session.DevAddr).OrZero(),
								len(tc.CreateDevice.EndDevice.Session.QueuedApplicationDownlinks),
								tc.CreateDevice.Session.Keys,
							)
						}(),
					)
				}(),
				Parallel: true,
				Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
					nsConf := DefaultConfig
					nsConf.DeviceKEKLabel = test.DefaultKEKLabel

					ns, ctx, env, stop := StartTest(ctx, TestConfig{
						Component: component.Config{
							ServiceBase: config.ServiceBase{
								GRPC: config.GRPC{
									LogIgnoreMethods: []string{
										"/ttn.lorawan.v3.ApplicationAccess/ListRights",
										"/ttn.lorawan.v3.NsEndDeviceRegistry/ResetFactoryDefaults",
									},
								},
								KeyVault: test.DefaultKeyVault,
								FrequencyPlans: config.FrequencyPlansConfig{
									ConfigSource: "static",
									Static:       test.StaticFrequencyPlans,
								},
							},
						},
						NetworkServer: nsConf,
						TaskStarter: StartTaskExclude(
							DownlinkProcessTaskName,
							DownlinkDispatchTaskName,
						),
					})
					defer stop()

					clock := test.NewMockClock(time.Now().UTC())
					defer SetMockClock(clock)()

					req := &ttnpb.ResetAndGetEndDeviceRequest{
						EndDeviceIds: test.MakeEndDeviceIdentifiers(),
						FieldMask:    ttnpb.FieldMask(conf.Paths...),
					}

					var created *ttnpb.EndDevice
					if tc.CreateDevice != nil {
						created, ctx = MustCreateDevice(ctx, env.Devices, tc.CreateDevice.EndDevice)

						req.EndDeviceIds.ApplicationIds = tc.CreateDevice.Ids.ApplicationIds
						req.EndDeviceIds.DeviceId = tc.CreateDevice.Ids.DeviceId

						clock.Add(time.Nanosecond)
					}

					dev, err, ok := env.AssertResetFactoryDefaults(ctx, req)
					if !a.So(ok, should.BeTrue) {
						return
					}
					a.So(dev, should.BeNil)
					if !a.So(err, should.BeError) || !a.So(errors.IsPermissionDenied(err), should.BeTrue) {
						t.Errorf("Expected 'permission denied' error, got: %s", test.FormatError(err))
						return
					}

					now := clock.Now().UTC()

					dev, err, ok = env.AssertResetFactoryDefaults(ctx, req, append([]ttnpb.Right{
						ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE,
					}, conf.RequiredRights...)...)
					if !a.So(ok, should.BeTrue) {
						return
					}
					if created == nil {
						a.So(err, should.NotBeNil)
						if !a.So(errors.IsNotFound(err), should.BeTrue) {
							t.Errorf("Expected 'not found' error, got: %s", test.FormatError(err))
						}
						return
					}

					var (
						macState *ttnpb.MACState
						session  *ttnpb.Session
					)
					if !created.SupportsJoin {
						if created.Session == nil {
							a.So(err, should.NotBeNil)
							if !a.So(errors.IsDataLoss(err), should.BeTrue) {
								t.Errorf("Expected 'data loss' error, got: %s", test.FormatError(err))
							}
							return
						}

						fps, err := ns.FrequencyPlansStore(ctx)
						if !a.So(err, should.BeNil) {
							t.Fail()
							return
						}
						var newErr error
						defaultMACSettings := test.Must(DefaultConfig.DefaultMACSettings.Parse())
						macState, newErr = mac.NewState(created, fps, defaultMACSettings, nil)
						if newErr != nil {
							a.So(err, should.NotBeNil)
							a.So(err, should.HaveSameErrorDefinitionAs, newErr)
							return
						}
						session = &ttnpb.Session{
							DevAddr:                    created.Session.DevAddr,
							QueuedApplicationDownlinks: created.Session.QueuedApplicationDownlinks,
							Keys:                       created.Session.Keys,
							StartedAt:                  timestamppb.New(now),
						}
					}
					if !a.So(err, should.BeNil) {
						t.Errorf("Expected no error, got: %s", test.FormatError(err))
						return
					}

					expected := ttnpb.Clone(created)
					expected.BatteryPercentage = nil
					expected.DownlinkMargin = 0
					expected.LastDevStatusReceivedAt = nil
					expected.MacState = macState
					expected.PendingMacState = nil
					expected.PendingSession = nil
					expected.PowerState = ttnpb.PowerState_POWER_UNKNOWN
					expected.Session = session
					expected.UpdatedAt = timestamppb.New(clock.Now())
					if !a.So(dev, should.Resemble, test.Must(ttnpb.ApplyEndDeviceFieldMask(nil, expected, ttnpb.AddImplicitEndDeviceGetFields(conf.Paths...)...))) {
						return
					}
					updated, _, err := env.Devices.GetByID(ctx, tc.CreateDevice.Ids.ApplicationIds, tc.CreateDevice.Ids.DeviceId, ttnpb.EndDeviceFieldPathsTopLevel)
					if a.So(err, should.BeNil) {
						a.So(updated, should.Resemble, expected)
					}
				},
			})
		}
	}
}

func TestDeviceRegistryDelete(t *testing.T) {
	for _, tc := range []struct {
		Name           string
		ContextFunc    func(context.Context) context.Context
		SetByIDFunc    func(context.Context, *ttnpb.ApplicationIdentifiers, string, []string, func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error)
		Request        *ttnpb.EndDeviceIdentifiers
		ErrorAssertion func(*testing.T, error) bool
		SetByIDCalls   uint64
	}{
		{
			Name: "No device write rights",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, &rights.Rights{
					ApplicationRights: *rights.NewMap(map[string]*ttnpb.Rights{
						unique.ID(test.Context(), &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.Right_RIGHT_GATEWAY_SETTINGS_BASIC,
							},
						},
					}),
				})
			},
			SetByIDFunc: func(ctx context.Context, appID *ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
				err := errors.New("SetByIDFunc must not be called")
				test.MustTFromContext(ctx).Error(err)
				return nil, ctx, err
			},
			Request: &ttnpb.EndDeviceIdentifiers{
				DeviceId:       "test-dev-id",
				ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"},
			},
			ErrorAssertion: func(t *testing.T, err error) bool {
				if !assertions.New(t).So(errors.IsPermissionDenied(err), should.BeTrue) {
					t.Errorf("Received error: %s", err)
					return false
				}
				return true
			},
		},

		{
			Name: "Non-existing device",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, &rights.Rights{
					ApplicationRights: *rights.NewMap(map[string]*ttnpb.Rights{
						unique.ID(test.Context(), &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE,
							},
						},
					}),
				})
			},
			SetByIDFunc: func(ctx context.Context, appID *ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)
				a.So(appID, should.Resemble, &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.BeNil)

				dev, sets, err := f(ctx, nil)
				if !a.So(errors.IsNotFound(err), should.BeTrue) {
					return nil, ctx, err
				}
				a.So(sets, should.BeNil)
				a.So(dev, should.BeNil)
				return nil, ctx, nil
			},
			Request: &ttnpb.EndDeviceIdentifiers{
				DeviceId:       "test-dev-id",
				ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"},
			},
			SetByIDCalls: 1,
		},

		{
			Name: "Existing device",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, &rights.Rights{
					ApplicationRights: *rights.NewMap(map[string]*ttnpb.Rights{
						unique.ID(test.Context(), &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE,
							},
						},
					}),
				})
			},
			SetByIDFunc: func(ctx context.Context, appID *ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.BeNil)

				dev, sets, err := f(ctx, &ttnpb.EndDevice{
					Ids: &ttnpb.EndDeviceIdentifiers{
						DeviceId:       "test-dev-id",
						ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"},
					},
				})
				if !a.So(err, should.BeNil) {
					return nil, ctx, err
				}
				a.So(sets, should.BeNil)
				a.So(dev, should.BeNil)
				return nil, ctx, nil
			},
			Request: &ttnpb.EndDeviceIdentifiers{
				DeviceId:       "test-dev-id",
				ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-app-id"},
			},
			SetByIDCalls: 1,
		},
	} {
		tc := tc
		test.RunSubtest(t, test.SubtestConfig{
			Name:     tc.Name,
			Parallel: true,
			Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
				var setByIDCalls uint64

				ns, ctx, env, stop := StartTest(
					ctx,
					TestConfig{
						NetworkServer: Config{
							Devices: &MockDeviceRegistry{
								SetByIDFunc: func(ctx context.Context, appID *ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
									atomic.AddUint64(&setByIDCalls, 1)
									return tc.SetByIDFunc(ctx, appID, devID, gets, f)
								},
							},
						},
						TaskStarter: StartTaskExclude(
							DownlinkProcessTaskName,
							DownlinkDispatchTaskName,
						),
					},
				)
				defer stop()

				go LogEvents(t, env.Events)

				ns.AddContextFiller(tc.ContextFunc)
				ns.AddContextFiller(func(ctx context.Context) context.Context {
					return test.ContextWithTB(ctx, t)
				})

				req := ttnpb.Clone(tc.Request)
				res, err := ttnpb.NewNsEndDeviceRegistryClient(ns.LoopbackConn()).Delete(ctx, req)
				a.So(setByIDCalls, should.Equal, tc.SetByIDCalls)
				if tc.ErrorAssertion != nil && a.So(tc.ErrorAssertion(t, err), should.BeTrue) {
					a.So(res, should.BeNil)
				} else if a.So(err, should.BeNil) {
					a.So(res, should.Resemble, ttnpb.Empty)
				}
				a.So(req, should.Resemble, tc.Request)
			},
		})
	}
}

func TestDeviceRegistryBatchDelete(t *testing.T) { // nolint:paralleltest
	registeredApplicationID := "test-app"
	registeredApplicationIDs := &ttnpb.ApplicationIdentifiers{
		ApplicationId: registeredApplicationID,
	}
	dev1 := &ttnpb.EndDevice{
		Ids: &ttnpb.EndDeviceIdentifiers{
			ApplicationIds: &ttnpb.ApplicationIdentifiers{
				ApplicationId: registeredApplicationID,
			},
			DeviceId: "test-device-1",
			JoinEui:  types.EUI64{0x42, 0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}.Bytes(),
			DevEui:   types.EUI64{0x42, 0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}.Bytes(),
		},
		RootKeys: &ttnpb.RootKeys{
			RootKeyId: "testKey",
			NwkKey: &ttnpb.KeyEnvelope{
				Key: types.AES128Key{
					0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x0,
				}.Bytes(),
				KekLabel: "test",
			},
			AppKey: &ttnpb.KeyEnvelope{
				Key: types.AES128Key{
					0x0, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
				}.Bytes(),
				KekLabel: "test",
			},
		},
	}
	dev2 := ttnpb.Clone(dev1)
	dev2.Ids.DeviceId = "test-device-2"
	dev2.Ids.JoinEui = types.EUI64{0x42, 0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}.Bytes()
	dev2.Ids.DevEui = types.EUI64{0x42, 0x43, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}.Bytes()

	dev3 := ttnpb.Clone(dev1)
	dev3.Ids.DeviceId = "test-device-3"
	dev3.Ids.JoinEui = types.EUI64{0x42, 0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}.Bytes()
	dev3.Ids.DevEui = types.EUI64{0x42, 0x44, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}.Bytes()

	for _, tc := range []struct {
		Name            string
		ContextFunc     func(context.Context) context.Context
		BatchDeleteFunc func(
			ctx context.Context,
			appIDs *ttnpb.ApplicationIdentifiers,
			deviceIDs []string,
		) ([]*ttnpb.EndDeviceIdentifiers, error)
		Request          *ttnpb.BatchDeleteEndDevicesRequest
		ErrorAssertion   func(*testing.T, error) bool
		BatchDeleteCalls uint64
	}{
		{
			Name: "No device write rights",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, &rights.Rights{
					ApplicationRights: *rights.NewMap(map[string]*ttnpb.Rights{
						unique.ID(test.Context(), registeredApplicationIDs): {
							Rights: []ttnpb.Right{
								ttnpb.Right_RIGHT_GATEWAY_SETTINGS_BASIC,
							},
						},
					}),
				})
			},
			BatchDeleteFunc: func(
				ctx context.Context,
				appIDs *ttnpb.ApplicationIdentifiers,
				deviceIDs []string,
			) ([]*ttnpb.EndDeviceIdentifiers, error) {
				err := errors.New("BatchDeleteFunc must not be called")
				test.MustTFromContext(ctx).Error(err)
				return nil, err
			},
			Request: &ttnpb.BatchDeleteEndDevicesRequest{
				ApplicationIds: registeredApplicationIDs,
				DeviceIds: []string{
					dev1.Ids.DeviceId,
					dev2.Ids.DeviceId,
					dev3.Ids.DeviceId,
				},
			},
			ErrorAssertion: func(t *testing.T, err error) bool {
				t.Helper()
				if !assertions.New(t).So(errors.IsPermissionDenied(err), should.BeTrue) {
					t.Errorf("Received error: %s", err)
					return false
				}
				return true
			},
			BatchDeleteCalls: 0,
		},
		{
			Name: "Non-existing device",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, &rights.Rights{
					ApplicationRights: *rights.NewMap(map[string]*ttnpb.Rights{
						unique.ID(test.Context(), registeredApplicationIDs): {
							Rights: []ttnpb.Right{
								ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE,
							},
						},
					}),
				})
			},
			BatchDeleteFunc: func(ctx context.Context,
				appIDs *ttnpb.ApplicationIdentifiers,
				deviceIDs []string,
			) ([]*ttnpb.EndDeviceIdentifiers, error) {
				// Devices not found are skipped.
				return nil, nil
			},
			Request: &ttnpb.BatchDeleteEndDevicesRequest{
				ApplicationIds: registeredApplicationIDs,
				DeviceIds: []string{
					dev1.Ids.DeviceId,
					dev2.Ids.DeviceId,
					dev3.Ids.DeviceId,
				},
			},
			BatchDeleteCalls: 1,
		},
		{
			Name: "Wrong application ID",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, &rights.Rights{
					ApplicationRights: *rights.NewMap(map[string]*ttnpb.Rights{
						unique.ID(test.Context(), registeredApplicationIDs): {
							Rights: []ttnpb.Right{
								ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE,
							},
						},
					}),
				})
			},
			Request: &ttnpb.BatchDeleteEndDevicesRequest{
				ApplicationIds: &ttnpb.ApplicationIdentifiers{ApplicationId: "test-unknown-app-id"},
				DeviceIds: []string{
					dev1.Ids.DeviceId,
					dev2.Ids.DeviceId,
					dev3.Ids.DeviceId,
				},
			},
			BatchDeleteFunc: func(
				ctx context.Context,
				appIDs *ttnpb.ApplicationIdentifiers,
				deviceIDs []string,
			) ([]*ttnpb.EndDeviceIdentifiers, error) {
				err := errors.New("BatchDeleteFunc must not be called")
				test.MustTFromContext(ctx).Error(err)
				return nil, err
			},
			ErrorAssertion: func(t *testing.T, err error) bool {
				t.Helper()
				if !assertions.New(t).So(errors.IsPermissionDenied(err), should.BeTrue) {
					t.Errorf("Received error: %s", err)
					return false
				}
				return true
			},
			BatchDeleteCalls: 0,
		},
		{
			Name: "Invalid Device",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, &rights.Rights{
					ApplicationRights: *rights.NewMap(map[string]*ttnpb.Rights{
						unique.ID(test.Context(), registeredApplicationIDs): {
							Rights: []ttnpb.Right{
								ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE,
							},
						},
					}),
				})
			},
			Request: &ttnpb.BatchDeleteEndDevicesRequest{
				ApplicationIds: registeredApplicationIDs,
				DeviceIds: []string{
					"test-dev-&*@(#)",
				},
			},
			BatchDeleteFunc: func(
				ctx context.Context,
				appIDs *ttnpb.ApplicationIdentifiers,
				deviceIDs []string,
			) ([]*ttnpb.EndDeviceIdentifiers, error) {
				err := errors.New("BatchDeleteFunc must not be called")
				test.MustTFromContext(ctx).Error(err)
				return nil, err
			},
			ErrorAssertion: func(t *testing.T, err error) bool {
				t.Helper()
				if !assertions.New(t).So(errors.IsInvalidArgument(err), should.BeTrue) {
					t.Errorf("Received error: %s", err)
					return false
				}
				return true
			},
			BatchDeleteCalls: 0,
		},
		{
			Name: "Existing device",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, &rights.Rights{
					ApplicationRights: *rights.NewMap(map[string]*ttnpb.Rights{
						unique.ID(test.Context(), registeredApplicationIDs): {
							Rights: []ttnpb.Right{
								ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE,
							},
						},
					}),
				})
			},
			BatchDeleteFunc: func(
				ctx context.Context,
				appIDs *ttnpb.ApplicationIdentifiers,
				deviceIDs []string,
			) ([]*ttnpb.EndDeviceIdentifiers, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(deviceIDs, should.HaveLength, 1)
				a.So(appIDs, should.Resemble, registeredApplicationIDs)
				a.So(deviceIDs[0], should.Equal, dev1.GetIds().DeviceId)
				return []*ttnpb.EndDeviceIdentifiers{
					dev1.Ids,
				}, nil
			},
			Request: &ttnpb.BatchDeleteEndDevicesRequest{
				ApplicationIds: registeredApplicationIDs,
				DeviceIds: []string{
					dev1.Ids.DeviceId,
				},
			},
			BatchDeleteCalls: 1,
		},
		{
			Name: "One invalid device in batch",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, &rights.Rights{
					ApplicationRights: *rights.NewMap(map[string]*ttnpb.Rights{
						unique.ID(test.Context(), registeredApplicationIDs): {
							Rights: []ttnpb.Right{
								ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE,
							},
						},
					}),
				})
			},
			BatchDeleteFunc: func(
				ctx context.Context,
				appIDs *ttnpb.ApplicationIdentifiers,
				deviceIDs []string,
			) ([]*ttnpb.EndDeviceIdentifiers, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(deviceIDs, should.HaveLength, 3)
				a.So(appIDs, should.Resemble, registeredApplicationIDs)

				for _, devID := range deviceIDs {
					switch devID {
					case dev1.GetIds().DeviceId:
					case dev2.GetIds().DeviceId:
						t.Log("Known device ID")
					case "test-dev-unknown-id":
						t.Log("Ignore expected unknown device ID")
					default:
						t.Log("Unexpected device ID")
					}
				}
				return []*ttnpb.EndDeviceIdentifiers{
					dev1.Ids,
					dev2.Ids,
				}, nil
			},
			Request: &ttnpb.BatchDeleteEndDevicesRequest{
				ApplicationIds: registeredApplicationIDs,
				DeviceIds: []string{
					dev1.Ids.DeviceId,
					dev2.Ids.DeviceId,
					"test-dev-unknown-id",
				},
			},
			BatchDeleteCalls: 1,
		},
		{
			Name: "Valid Batch",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, &rights.Rights{
					ApplicationRights: *rights.NewMap(map[string]*ttnpb.Rights{
						unique.ID(test.Context(), registeredApplicationIDs): {
							Rights: []ttnpb.Right{
								ttnpb.Right_RIGHT_APPLICATION_DEVICES_WRITE,
							},
						},
					}),
				})
			},
			BatchDeleteFunc: func(
				ctx context.Context,
				appIDs *ttnpb.ApplicationIdentifiers,
				deviceIDs []string,
			) ([]*ttnpb.EndDeviceIdentifiers, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appIDs, should.Resemble, registeredApplicationIDs)
				a.So(deviceIDs, should.HaveLength, 3)
				for _, devID := range deviceIDs {
					switch devID {
					case dev1.GetIds().DeviceId:
					case dev2.GetIds().DeviceId:
					case dev3.GetIds().DeviceId:
						// Known device ID
					default:
						t.Error("Unknown device ID: ", devID)
					}
				}
				return []*ttnpb.EndDeviceIdentifiers{
					dev1.Ids,
					dev2.Ids,
					dev3.Ids,
				}, nil
			},
			Request: &ttnpb.BatchDeleteEndDevicesRequest{
				ApplicationIds: registeredApplicationIDs,
				DeviceIds: []string{
					dev1.Ids.DeviceId,
					dev2.Ids.DeviceId,
					dev3.Ids.DeviceId,
				},
			},
			BatchDeleteCalls: 1,
		},
	} {
		tc := tc
		test.RunSubtest(t, test.SubtestConfig{
			Name:     tc.Name,
			Parallel: true,
			Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
				t.Helper()
				var batchDeleteCalls uint64
				ns, ctx, env, stop := StartTest(
					ctx,
					TestConfig{
						NetworkServer: Config{
							Devices: &MockDeviceRegistry{
								BatchDeleteFunc: func(
									ctx context.Context,
									appIDs *ttnpb.ApplicationIdentifiers,
									deviceIDs []string,
								) ([]*ttnpb.EndDeviceIdentifiers, error) {
									atomic.AddUint64(&batchDeleteCalls, 1)
									return tc.BatchDeleteFunc(ctx, appIDs, deviceIDs)
								},
							},
						},
						TaskStarter: StartTaskExclude(
							DownlinkProcessTaskName,
							DownlinkDispatchTaskName,
						),
					},
				)
				defer stop()

				go LogEvents(t, env.Events)

				ns.AddContextFiller(tc.ContextFunc)
				ns.AddContextFiller(func(ctx context.Context) context.Context {
					return test.ContextWithTB(ctx, t)
				})

				req := ttnpb.Clone(tc.Request)
				a.So(req, should.Resemble, tc.Request)

				_, err := ttnpb.NewNsEndDeviceBatchRegistryClient(ns.LoopbackConn()).Delete(ctx, req)
				a.So(batchDeleteCalls, should.Equal, tc.BatchDeleteCalls)
				if tc.ErrorAssertion != nil {
					a.So(tc.ErrorAssertion(t, err), should.BeTrue)
				} else {
					a.So(err, should.BeNil)
				}
			},
		})
	}
}
