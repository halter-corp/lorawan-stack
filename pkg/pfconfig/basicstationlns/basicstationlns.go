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

// Package basicstationlns implements the JSON configuration for the Basic Station `router_config` message.
package basicstationlns

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.thethings.network/lorawan-stack/pkg/band"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/frequencyplans"
	pfconfig "go.thethings.network/lorawan-stack/pkg/pfconfig/shared"
)

const (
	configHardwareSpecPrefix            = "sx1301"
	configHardwareSpecNoOfConcentrators = "1"
)

var errFrequencyPlan = errors.DefineInvalidArgument("frequency_plan", "invalid frequency plan")

// DataRates encodes the available datarates of the channel plan for the Station in the format below:
// [0] -> SF (Spreading Factor; Range: 7...12 for LoRa, 0 for FSK)
// [1] -> BW (Bandwidth; 125/250/500 for LoRa, ignored for FSK)
// [2] -> DNONLY (Downlink Only; 1 = true, 0 = false)
type DataRates [16][3]int

// RouterConfig contains the router configuration.
// This message is sent by the Gateway Server.
type RouterConfig struct {
	NetID          []int                   `json:"NetID"`
	JoinEUI        [][]int                 `json:"JoinEui"`
	Region         string                  `json:"region"`
	HardwareSpec   string                  `json:"hwspec"`
	FrequencyRange []int                   `json:"freq_range"`
	DataRates      DataRates               `json:"DRs"`
	SX1301Config   []pfconfig.SX1301Config `json:"sx1301_conf"`

	// These are debug options to be unset in production gateways.
	NoCCA       bool `json:"nocca"`
	NoDutyCycle bool `json:"nodc"`
	NoDwellTime bool `json:"nodwell"`

	MuxTime float64 `json:"MuxTime"`
}

// MarshalJSON implements json.Marshaler.
func (cfg RouterConfig) MarshalJSON() ([]byte, error) {
	type Alias RouterConfig
	return json.Marshal(struct {
		Type string `json:"msgtype"`
		Alias
	}{
		Type:  "router_config",
		Alias: Alias(cfg),
	})
}

// GetRouterConfig returns the routerconfig message to be sent to the gateway.
func GetRouterConfig(fp frequencyplans.FrequencyPlan, isProd bool, dlTime time.Time) (RouterConfig, error) {
	if err := fp.Validate(); err != nil {
		return RouterConfig{}, errFrequencyPlan
	}

	cfg := RouterConfig{}
	cfg.JoinEUI = nil
	cfg.NetID = nil

	band, err := band.GetByID(fp.BandID)
	if err != nil {
		return RouterConfig{}, errFrequencyPlan
	}

	s := strings.Split(band.ID, "_")
	if len(s) < 2 {
		return RouterConfig{}, errFrequencyPlan
	}
	cfg.Region = fmt.Sprintf("%s%s", s[0], s[1])
	if len(fp.Radios) == 0 {
		return RouterConfig{}, errFrequencyPlan
	}
	// TODO: Handle FP with multiple radios if necessary (https://github.com/TheThingsNetwork/lorawan-stack/issues/761).
	cfg.FrequencyRange = []int{
		int(fp.Radios[0].TxConfiguration.MinFrequency),
		int(fp.Radios[0].TxConfiguration.MaxFrequency),
	}

	// TODO: Dynamically fill this field based on no of SX1301_conf objects (https://github.com/TheThingsNetwork/lorawan-stack/issues/761).
	cfg.HardwareSpec = fmt.Sprintf("%s/%s", configHardwareSpecPrefix, configHardwareSpecNoOfConcentrators)

	drs, err := getDataRatesFromBandID(fp.BandID)
	if err != nil {
		return RouterConfig{}, errFrequencyPlan
	}
	cfg.DataRates = drs

	cfg.NoCCA = !isProd
	cfg.NoDutyCycle = !isProd
	cfg.NoDwellTime = !isProd

	sx1301Conf, err := pfconfig.BuildSX1301Config(&fp)
	if err != nil {
		return RouterConfig{}, err
	}

	// These fields are not defined in the v1.5 ref design https://doc.sm.tc/station/gw_v1.5.html#rfconf-object and would cause a parsing error.
	sx1301Conf.Radios[0].TxFreqMin = 0
	sx1301Conf.Radios[0].TxFreqMax = 0

	// TODO: Extend this for > 8ch gateways (https://github.com/TheThingsNetwork/lorawan-stack/issues/761).
	cfg.SX1301Config = append(cfg.SX1301Config, *sx1301Conf)

	// Add the MuxTime for RTT measurement.
	cfg.MuxTime = float64(dlTime.Unix()) + float64(dlTime.Nanosecond())/(1e9)

	return cfg, nil
}

// getDataRatesFromBandID parses the available data rates from the band into DataRates.
func getDataRatesFromBandID(id string) (DataRates, error) {
	band, err := band.GetByID(id)
	if err != nil {
		return DataRates{}, err
	}

	// Set the default values.
	drs := DataRates{}
	for _, dr := range drs {
		dr[0] = -1
		dr[1] = 0
		dr[2] = 0
	}

	for i, dr := range band.DataRates {
		if loraDR := dr.Rate.GetLoRa(); loraDR != nil {
			loraDR.GetSpreadingFactor()
			drs[i][0] = int(loraDR.GetSpreadingFactor())
			drs[i][1] = int(loraDR.GetBandwidth() / 1000)
		} else if fskDR := dr.Rate.GetFSK(); fskDR != nil {
			drs[i][0] = 0 // must be set to 0 for FSK, the BW field is ignored.
		}
	}
	return drs, nil
}
