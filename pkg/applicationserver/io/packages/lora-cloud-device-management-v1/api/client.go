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

package api

import (
	"io"
	"net/http"
	"path"

	"go.thethings.network/lorawan-stack/pkg/version"
)

// Option is an option for the API client.
type Option interface {
	apply(*Client)
}

// OptionFunc is an Option implemented as a function.
type OptionFunc func(*Client)

func (f OptionFunc) apply(c *Client) { f(c) }

// Client is an API client for the LoRa Cloud Device Management v1 service.
type Client struct {
	token string
}

const (
	baseURL     = "/api/v1"
	contentType = "application/json"
)

var (
	userAgent = "ttn-lw-application-server/" + version.TTN
)

func (c *Client) newRequest(method, category, entity, operation string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, path.Join(baseURL, category, entity, operation), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", userAgent)
	if c.token != "" {
		req.Header.Set("Authorization", c.token)
	}
	return req, nil
}

// WithToken uses the given authentication token in the client.
func WithToken(token string) Option {
	return OptionFunc(func(c *Client) {
		c.token = token
	})
}

// New creates a new Client with the given options.
func New(opts ...Option) (*Client, error) {
	client := &Client{}
	for _, opt := range opts {
		opt.apply(client)
	}
	return client, nil
}
