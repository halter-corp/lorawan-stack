// Copyright © 2024 The Things Network Foundation, The Things Industries B.V.
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

// Package redis implements a Redis-backed MAC settings profile registry.
package redis

import (
	"context"
	"testing"

	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

var Timeout = 10 * test.Delay

func TestMACSettingsProfileRegistry(t *testing.T) {
	t.Parallel()
	a, ctx := test.New(t)
	cl, flush := test.NewRedis(ctx, "redis_test")
	t.Cleanup(func() {
		flush()
		cl.Close()
	})

	ids := &ttnpb.MACSettingsProfileIdentifiers{
		ApplicationIds: &ttnpb.ApplicationIdentifiers{
			ApplicationId: "myapp",
		},
		ProfileId: "prof-01",
	}
	ids2 := &ttnpb.MACSettingsProfileIdentifiers{
		ApplicationIds: &ttnpb.ApplicationIdentifiers{
			ApplicationId: "myapp-02",
		},
		ProfileId: "prof-02",
	}

	frequencies := []uint64{868100000, 868300000, 868500000}
	frequencies2 := []uint64{868100000, 868300000, 868500000, 868700000}

	registry := &MACSettingsProfileRegistry{
		Redis:   cl,
		LockTTL: test.Delay << 10,
	}
	if err := registry.Init(ctx); !a.So(err, should.BeNil) {
		t.FailNow()
	}

	createProfileFunc := func(_ context.Context, pb *ttnpb.MACSettingsProfile,
	) (*ttnpb.MACSettingsProfile, []string, error) { // nolint: unparam
		a.So(pb, should.BeNil)
		return &ttnpb.MACSettingsProfile{
			Ids: ids,
			MacSettings: &ttnpb.MACSettings{
				ResetsFCnt: &ttnpb.BoolValue{Value: true},
			},
		}, []string{"ids", "mac_settings"}, nil
	}

	updateProfileFunc := func(_ context.Context, pb *ttnpb.MACSettingsProfile,
	) (*ttnpb.MACSettingsProfile, []string, error) { // nolint: unparam
		a.So(pb, should.NotBeNil)
		return &ttnpb.MACSettingsProfile{
			Ids: ids,
			MacSettings: &ttnpb.MACSettings{
				ResetsFCnt:               &ttnpb.BoolValue{Value: false},
				FactoryPresetFrequencies: frequencies,
			},
		}, []string{"ids", "mac_settings"}, nil
	}

	updateFieldMaskProfileFunc := func(_ context.Context, pb *ttnpb.MACSettingsProfile,
	) (*ttnpb.MACSettingsProfile, []string, error) { // nolint: unparam
		a.So(pb, should.NotBeNil)
		return &ttnpb.MACSettingsProfile{
			Ids: ids,
			MacSettings: &ttnpb.MACSettings{
				ResetsFCnt:               &ttnpb.BoolValue{Value: true},
				FactoryPresetFrequencies: frequencies2,
			},
		}, []string{"ids", "mac_settings.factory_preset_frequencies"}, nil
	}

	deleteProfileFunc := func(_ context.Context, pb *ttnpb.MACSettingsProfile,
	) (*ttnpb.MACSettingsProfile, []string, error) { // nolint: unparam
		a.So(pb, should.NotBeNil)
		return nil, nil, nil
	}

	listProfileFunc := func(_ context.Context, pb *ttnpb.MACSettingsProfile,
	) (*ttnpb.MACSettingsProfile, []string, error) { // nolint: unparam
		a.So(pb, should.BeNil)
		return &ttnpb.MACSettingsProfile{
			Ids: ids2,
			MacSettings: &ttnpb.MACSettings{
				ResetsFCnt: &ttnpb.BoolValue{Value: true},
			},
		}, []string{"ids", "mac_settings"}, nil
	}

	t.Run("GetNonExisting", func(t *testing.T) {
		t.Parallel()
		a, ctx := test.New(t)
		profile, err := registry.Get(ctx, ids, []string{"ids"})
		a.So(profile, should.BeNil)
		a.So(errors.IsNotFound(err), should.BeTrue)
	})

	t.Run("CreateReadUpdateDelete", func(t *testing.T) {
		t.Parallel()
		a, ctx := test.New(t)
		profile, err := registry.Set(ctx, ids, []string{"ids", "mac_settings"}, createProfileFunc)
		a.So(err, should.BeNil)
		a.So(profile, should.NotBeNil)
		a.So(profile.Ids, should.Resemble, ids)
		a.So(profile.MacSettings, should.NotBeNil)
		a.So(profile.MacSettings.ResetsFCnt, should.NotBeNil)
		a.So(profile.MacSettings.ResetsFCnt.Value, should.BeTrue)

		retrieved, err := registry.Get(ctx, ids, []string{"ids", "mac_settings"})
		a.So(err, should.BeNil)
		a.So(retrieved, should.NotBeNil)
		a.So(retrieved.Ids, should.Resemble, ids)
		a.So(retrieved.MacSettings, should.NotBeNil)
		a.So(retrieved.MacSettings.ResetsFCnt, should.NotBeNil)
		a.So(retrieved.MacSettings.ResetsFCnt.Value, should.BeTrue)

		updated, err := registry.Set(ctx, ids, []string{"ids", "mac_settings"}, updateProfileFunc)
		a.So(err, should.BeNil)
		a.So(updated, should.NotBeNil)
		a.So(updated.Ids, should.Resemble, ids)
		a.So(updated.MacSettings, should.NotBeNil)
		a.So(updated.MacSettings.ResetsFCnt, should.NotBeNil)
		a.So(updated.MacSettings.ResetsFCnt.Value, should.BeFalse)
		a.So(updated.MacSettings.FactoryPresetFrequencies, should.Resemble, frequencies)

		updated2, err := registry.Set(ctx, ids, []string{"ids", "mac_settings"}, updateFieldMaskProfileFunc)
		a.So(err, should.BeNil)
		a.So(updated2, should.NotBeNil)
		a.So(updated2.Ids, should.Resemble, ids)
		a.So(updated2.MacSettings, should.NotBeNil)
		a.So(updated2.MacSettings.ResetsFCnt, should.NotBeNil)
		a.So(updated2.MacSettings.ResetsFCnt.Value, should.BeFalse)
		a.So(updated2.MacSettings.FactoryPresetFrequencies, should.Resemble, frequencies2)

		deleted, err := registry.Set(ctx, ids, []string{"ids", "mac_settings"}, deleteProfileFunc)
		a.So(err, should.BeNil)
		a.So(deleted, should.BeNil)
	})

	t.Run("List", func(t *testing.T) {
		t.Parallel()
		a, ctx := test.New(t)
		profile, err := registry.Set(ctx, ids2, []string{"ids", "mac_settings"}, listProfileFunc)
		a.So(err, should.BeNil)
		a.So(profile, should.NotBeNil)
		a.So(profile.Ids, should.Resemble, ids2)
		a.So(profile.MacSettings, should.NotBeNil)
		a.So(profile.MacSettings.ResetsFCnt, should.NotBeNil)
		a.So(profile.MacSettings.ResetsFCnt.Value, should.BeTrue)

		profiles, err := registry.List(ctx, ids2.ApplicationIds, []string{"ids", "mac_settings"})
		a.So(err, should.BeNil)
		a.So(profiles, should.HaveLength, 1)
		a.So(profiles[0], should.NotBeNil)
		a.So(profiles[0].Ids, should.Resemble, ids2)
		a.So(profiles[0].MacSettings, should.NotBeNil)
		a.So(profiles[0].MacSettings.ResetsFCnt, should.NotBeNil)
		a.So(profiles[0].MacSettings.ResetsFCnt.Value, should.BeTrue)
	})
}
