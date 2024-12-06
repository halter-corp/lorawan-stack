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

package cryptoutil

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/pem"
	"maps"
	"strings"
	"sync"

	"go.thethings.network/lorawan-stack/v3/pkg/crypto"
)

type memKeyVault struct {
	mu sync.Mutex
	m  map[string][]byte
}

// Key implements crypto.KeyVault.
func (kv *memKeyVault) Key(_ context.Context, label string) ([]byte, error) {
	kv.mu.Lock()
	key, ok := kv.m[label]
	kv.mu.Unlock()
	if !ok {
		return nil, errKeyNotFound.WithAttributes("label", label)
	}
	return key, nil
}

func (kv *memKeyVault) SetKey(_ context.Context, label string, key []byte) error {
	kv.mu.Lock()
	kv.m[label] = key
	kv.mu.Unlock()
	return nil
}

func (kv *memKeyVault) certificate(label string) (tls.Certificate, error) {
	kv.mu.Lock()
	raw, ok := kv.m[label]
	kv.mu.Unlock()
	if !ok {
		return tls.Certificate{}, errCertificateNotFound.WithAttributes("label", label)
	}
	certPEMBlock, keyPEMBlock := &bytes.Buffer{}, &bytes.Buffer{}
	for {
		block, rest := pem.Decode(raw)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			if err := pem.Encode(certPEMBlock, block); err != nil {
				return tls.Certificate{}, err
			}
		} else if block.Type == "PRIVATE KEY" || strings.HasSuffix(block.Type, " PRIVATE KEY") {
			if err := pem.Encode(keyPEMBlock, block); err != nil {
				return tls.Certificate{}, err
			}
		}
		raw = rest
	}
	return tls.X509KeyPair(certPEMBlock.Bytes(), keyPEMBlock.Bytes())
}

// ServerCertificate implements crypto.KeyVault.
func (kv *memKeyVault) ServerCertificate(_ context.Context, label string) (tls.Certificate, error) {
	return kv.certificate(label)
}

// ClientCertificate implements crypto.KeyVault.
func (kv *memKeyVault) ClientCertificate(_ context.Context, label string) (tls.Certificate, error) {
	return kv.certificate(label)
}

// NewMemKeyVault returns a crypto.KeyVault that stores keys in memory.
// Certificates must be PEM encoded.
func NewMemKeyVault(m map[string][]byte) crypto.KeyVault {
	if m == nil {
		m = make(map[string][]byte)
	} else {
		m = maps.Clone(m)
	}
	return &memKeyVault{
		m: m,
	}
}
