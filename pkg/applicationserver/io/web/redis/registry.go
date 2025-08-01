// Copyright © 2022 The Things Network Foundation, The Things Industries B.V.
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

package redis

import (
	"context"
	"regexp"
	"runtime/trace"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	ttnredis "go.thethings.network/lorawan-stack/v3/pkg/redis"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	errInvalidFieldmask   = errors.DefineInvalidArgument("invalid_fieldmask", "invalid fieldmask")
	errInvalidIdentifiers = errors.DefineInvalidArgument("invalid_identifiers", "invalid identifiers")
	errReadOnlyField      = errors.DefineInvalidArgument("read_only_field", "read-only field `{field}`")
)

// appendImplicitWebhookGetPaths appends implicit ttnpb.ApplicationWebhook get paths to paths.
func appendImplicitWebhookGetPaths(paths ...string) []string {
	return append(append(make([]string, 0, 3+len(paths)),
		"created_at",
		"ids",
		"updated_at",
	), paths...)
}

func applyWebhookFieldMask(dst, src *ttnpb.ApplicationWebhook, paths ...string) (*ttnpb.ApplicationWebhook, error) {
	if dst == nil {
		dst = &ttnpb.ApplicationWebhook{}
	}
	return dst, dst.SetFields(src, paths...)
}

// WebhookRegistry is a Redis webhook registry.
type WebhookRegistry struct {
	Redis   *ttnredis.Client
	LockTTL time.Duration
}

// Init initializes the WebhookRegistry.
func (r *WebhookRegistry) Init(ctx context.Context) error {
	if err := ttnredis.InitMutex(ctx, r.Redis); err != nil {
		return err
	}
	return nil
}

func (r *WebhookRegistry) appKey(uid string) string {
	return r.Redis.Key("uid", uid)
}

func (r *WebhookRegistry) idKey(appUID, id string) string {
	return r.Redis.Key("uid", appUID, id)
}

func (r *WebhookRegistry) makeIDKeyFunc(appUID string) func(id string) string {
	return func(id string) string {
		return r.idKey(appUID, id)
	}
}

func webhookRegex(key string) (*regexp.Regexp, error) {
	keyRegex := strings.ReplaceAll(key, ":", "\\:")
	keyRegex = strings.ReplaceAll(keyRegex, "*", ".[^\\:]*")
	keyRegex = keyRegex + "$"
	return regexp.Compile(keyRegex)
}

// Get implements WebhookRegistry.
func (r WebhookRegistry) Get(ctx context.Context, ids *ttnpb.ApplicationWebhookIdentifiers, paths []string) (*ttnpb.ApplicationWebhook, error) {
	pb := &ttnpb.ApplicationWebhook{}
	if err := ttnredis.GetProto(ctx, r.Redis, r.idKey(unique.ID(ctx, ids.ApplicationIds), ids.WebhookId)).ScanProto(pb); err != nil {
		return nil, err
	}
	return applyWebhookFieldMask(nil, pb, appendImplicitWebhookGetPaths(paths...)...)
}

// List implements WebhookRegistry.
func (r WebhookRegistry) List(ctx context.Context, ids *ttnpb.ApplicationIdentifiers, paths []string) ([]*ttnpb.ApplicationWebhook, error) {
	var pbs []*ttnpb.ApplicationWebhook
	appUID := unique.ID(ctx, ids)
	uidKey := r.appKey(appUID)

	opts := []ttnredis.FindProtosOption{}
	limit, offset := ttnredis.PaginationLimitAndOffsetFromContext(ctx)
	if limit != 0 {
		opts = append(opts,
			ttnredis.FindProtosSorted(true),
			ttnredis.FindProtosWithOffsetAndCount(offset, limit),
		)
	}

	rangeProtos := func(c redis.Cmdable) error {
		return ttnredis.FindProtos(ctx, c, uidKey, r.makeIDKeyFunc(appUID), opts...).Range(
			func() (proto.Message, func() (bool, error)) {
				pb := &ttnpb.ApplicationWebhook{}
				return pb, func() (bool, error) {
					pb, err := applyWebhookFieldMask(nil, pb, appendImplicitWebhookGetPaths(paths...)...)
					if err != nil {
						return false, err
					}
					pbs = append(pbs, pb)
					return true, nil
				}
			})
	}

	defer trace.StartRegion(ctx, "list webhooks by application id").End()

	var err error
	if limit != 0 {
		var lockerID string
		lockerID, err = ttnredis.GenerateLockerID()
		if err != nil {
			return nil, err
		}
		err = ttnredis.LockedWatch(ctx, r.Redis, uidKey, lockerID, r.LockTTL, func(tx *redis.Tx) (err error) {
			total, err := tx.SCard(ctx, uidKey).Result()
			if err != nil {
				return err
			}
			ttnredis.SetPaginationTotal(ctx, total)
			return rangeProtos(tx)
		})
	} else {
		err = rangeProtos(r.Redis)
	}

	if err != nil {
		return nil, ttnredis.ConvertError(err)
	}
	return pbs, nil
}

// Set implements WebhookRegistry.
func (r WebhookRegistry) Set(ctx context.Context, ids *ttnpb.ApplicationWebhookIdentifiers, gets []string, f func(*ttnpb.ApplicationWebhook) (*ttnpb.ApplicationWebhook, []string, error)) (*ttnpb.ApplicationWebhook, error) {
	appUID := unique.ID(ctx, ids.ApplicationIds)
	ik := r.idKey(appUID, ids.WebhookId)

	lockerID, err := ttnredis.GenerateLockerID()
	if err != nil {
		return nil, err
	}

	var pb *ttnpb.ApplicationWebhook
	err = ttnredis.LockedWatch(ctx, r.Redis, ik, lockerID, r.LockTTL, func(tx *redis.Tx) error {
		cmd := ttnredis.GetProto(ctx, tx, ik)
		stored := &ttnpb.ApplicationWebhook{}
		if err := cmd.ScanProto(stored); errors.IsNotFound(err) {
			stored = nil
		} else if err != nil {
			return err
		}

		gets = appendImplicitWebhookGetPaths(gets...)

		var err error
		if stored != nil {
			pb = &ttnpb.ApplicationWebhook{}
			if err := cmd.ScanProto(pb); err != nil {
				return err
			}
			pb, err = applyWebhookFieldMask(nil, pb, gets...)
			if err != nil {
				return err
			}
		}

		var sets []string
		pb, sets, err = f(pb)
		if err != nil {
			return err
		}
		if stored == nil && pb == nil {
			return nil
		}
		if pb != nil && len(sets) == 0 {
			pb, err = applyWebhookFieldMask(nil, stored, gets...)
			return err
		}

		var pipelined func(redis.Pipeliner) error
		if pb == nil && len(sets) == 0 {
			pipelined = func(p redis.Pipeliner) error {
				p.Del(ctx, ik)
				p.SRem(ctx, r.appKey(appUID), stored.Ids.WebhookId)
				return nil
			}
		} else {
			if pb == nil {
				pb = &ttnpb.ApplicationWebhook{}
			}

			pb.UpdatedAt = timestamppb.Now()
			sets = append(append(sets[:0:0], sets...),
				"updated_at",
			)

			updated := &ttnpb.ApplicationWebhook{}
			if stored == nil {
				if err := ttnpb.RequireFields(sets,
					"ids.application_ids",
					"ids.webhook_id",
				); err != nil {
					return errInvalidFieldmask.WithCause(err)
				}

				pb.CreatedAt = pb.UpdatedAt
				sets = append(sets, "created_at")

				updated, err = applyWebhookFieldMask(updated, pb, sets...)
				if err != nil {
					return err
				}
				if updated.Ids.ApplicationIds.ApplicationId != ids.ApplicationIds.ApplicationId || updated.Ids.WebhookId != ids.WebhookId {
					return errInvalidIdentifiers.New()
				}
			} else {
				if ttnpb.HasAnyField(sets, "ids.application_ids.application_id") && pb.Ids.ApplicationIds.ApplicationId != stored.Ids.ApplicationIds.ApplicationId {
					return errReadOnlyField.WithAttributes("field", "ids.application_ids.application_id")
				}
				if ttnpb.HasAnyField(sets, "ids.webhook_id") && pb.Ids.WebhookId != stored.Ids.WebhookId {
					return errReadOnlyField.WithAttributes("field", "ids.webhook_id")
				}
				if err := cmd.ScanProto(updated); err != nil {
					return err
				}
				updated, err = applyWebhookFieldMask(updated, pb, sets...)
				if err != nil {
					return err
				}
			}
			if err := updated.ValidateFields(); err != nil {
				return err
			}

			pipelined = func(p redis.Pipeliner) error {
				if _, err := ttnredis.SetProto(ctx, p, ik, updated, 0); err != nil {
					return err
				}
				p.SAdd(ctx, r.appKey(appUID), updated.Ids.WebhookId)
				return nil
			}

			pb, err = applyWebhookFieldMask(nil, updated, gets...)
			if err != nil {
				return err
			}
		}
		_, err = tx.TxPipelined(ctx, pipelined)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, ttnredis.ConvertError(err)
	}
	return pb, nil
}

// Range implements WebhookRegistry.
func (r WebhookRegistry) Range(ctx context.Context, paths []string, f func(context.Context, *ttnpb.ApplicationIdentifiers, *ttnpb.ApplicationWebhook) bool) error {
	webhookEntityRegex, err := webhookRegex(r.idKey(unique.GenericID(ctx, "*"), "*"))
	if err != nil {
		return err
	}
	return ttnredis.RangeRedisKeys(ctx, r.Redis, r.idKey(unique.GenericID(ctx, "*"), "*"), ttnredis.DefaultRangeCount, func(key string) (bool, error) {
		if !webhookEntityRegex.MatchString(key) {
			return true, nil
		}
		wh := &ttnpb.ApplicationWebhook{}
		if err := ttnredis.GetProto(ctx, r.Redis, key).ScanProto(wh); err != nil {
			return false, err
		}
		wh, err := applyWebhookFieldMask(nil, wh, paths...)
		if err != nil {
			return false, err
		}
		if !f(ctx, wh.GetIds().GetApplicationIds(), wh) {
			return false, nil
		}
		return true, nil
	})
}

// WithPagination returns a new context with pagination parameters.
func (WebhookRegistry) WithPagination(
	ctx context.Context,
	limit uint32,
	page uint32,
	total *int64,
) context.Context {
	return ttnredis.NewContextWithPagination(ctx, int64(limit), int64(page), total)
}
