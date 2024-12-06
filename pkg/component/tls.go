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

package component

import (
	"context"
	"crypto/tls"

	"go.thethings.network/lorawan-stack/v3/pkg/config/tlsconfig"
)

func (c *Component) getTLSConfig(ctx context.Context) tlsconfig.Config {
	conf := c.GetBaseConfig(ctx).TLS

	// TODO: Remove detection mechanism (https://github.com/TheThingsNetwork/lorawan-stack/issues/1450)
	if conf.ServerAuth.Source == "" {
		switch {
		case conf.ServerAuth.Certificate != "" && conf.ServerAuth.Key != "":
			conf.ServerAuth.Source = "file"
		case conf.ServerAuth.KeyVault.ID != "":
			conf.ServerAuth.Source = "key-vault"
		case conf.ACME.Enable:
			conf.ServerAuth.Source = "acme"
		}
	}

	switch conf.ServerAuth.Source {
	case "acme":
		conf.ServerAuth.ACME = &conf.ACME
	case "key-vault":
		conf.ServerAuth.KeyVault.CertificateProvider = c.keyService
	}

	return conf
}

// GetTLSConfig gets the component's TLS config.
func (c *Component) GetTLSConfig(ctx context.Context) tlsconfig.Config {
	return c.getTLSConfig(ctx)
}

// GetTLSServerConfig gets the component's server TLS config and applies the given options.
func (c *Component) GetTLSServerConfig(ctx context.Context, opts ...tlsconfig.Option) (*tls.Config, error) {
	return tlsconfig.ConfigurationProvider(c.getTLSConfig).GetTLSServerConfig(ctx, opts...)
}

// GetTLSClientConfig gets the component's client TLS config and applies the given options.
func (c *Component) GetTLSClientConfig(ctx context.Context, opts ...tlsconfig.Option) (*tls.Config, error) {
	return tlsconfig.ConfigurationProvider(c.getTLSConfig).GetTLSClientConfig(ctx, opts...)
}
