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

package pubsub

import (
	"context"

	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

// Registry is a registry for pub/sub integrations.
type Registry interface {
	// Get returns the pub/sub integration by its identifiers.
	Get(ctx context.Context, ids *ttnpb.ApplicationPubSubIdentifiers, paths []string) (*ttnpb.ApplicationPubSub, error)
	// Range ranges over the pub/sub integrations and calls the callback function, until false is returned.
	Range(ctx context.Context, paths []string, f func(context.Context, *ttnpb.ApplicationIdentifiers, *ttnpb.ApplicationPubSub) bool) error
	// List returns all pub/sub integrations of the application.
	List(ctx context.Context, ids *ttnpb.ApplicationIdentifiers, paths []string) ([]*ttnpb.ApplicationPubSub, error)
	// Set creates, updates or deletes the pub/sub integration by its identifiers.
	Set(ctx context.Context, ids *ttnpb.ApplicationPubSubIdentifiers, paths []string, f func(*ttnpb.ApplicationPubSub) (*ttnpb.ApplicationPubSub, []string, error)) (*ttnpb.ApplicationPubSub, error)
	// WithPagination returns a new context with pagination parameters.
	WithPagination(ctx context.Context, limit uint32, page uint32, total *int64) context.Context
}
