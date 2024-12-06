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

package crypto

import (
	"context"
	"crypto/tls"
)

// KeyVault provides access to private keys.
type KeyVault interface {
	Key(ctx context.Context, label string) ([]byte, error)
	ServerCertificate(ctx context.Context, label string) (tls.Certificate, error)
	ClientCertificate(ctx context.Context, label string) (tls.Certificate, error)
}

// KeyVaultKeyWriter is a KeyVault that can set keys.
type KeyVaultKeyWriter interface {
	KeyVault
	SetKey(ctx context.Context, label string, key []byte) error
}
