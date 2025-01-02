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

package mac

import (
	"context"
	"math"

	"go.thethings.network/lorawan-stack/v3/pkg/band"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/frequencyplans"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/networkserver/internal"
	"go.thethings.network/lorawan-stack/v3/pkg/specification/macspec"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

var (
	EvtEnqueueLinkADRRequest = defineEnqueueMACRequestEvent(
		"link_adr", "link ADR",
		events.WithDataType(&ttnpb.MACCommand_LinkADRReq{}),
	)()
	EvtReceiveLinkADRAccept = defineReceiveMACAcceptEvent(
		"link_adr", "link ADR",
		events.WithDataType(&ttnpb.MACCommand_LinkADRAns{}),
	)()
	EvtReceiveLinkADRReject = defineReceiveMACRejectEvent(
		"link_adr", "link ADR",
		events.WithDataType(&ttnpb.MACCommand_LinkADRAns{}),
	)()
)

const (
	noChangeDataRateIndex = ttnpb.DataRateIndex_DATA_RATE_15
	noChangeTXPowerIndex  = 15
)

type linkADRReqParameters struct {
	Masks         []band.ChMaskCntlPair
	DataRateIndex ttnpb.DataRateIndex
	TxPowerIndex  uint32
	NbTrans       uint32
}

func generateLinkADRReq(ctx context.Context, dev *ttnpb.EndDevice, phy *band.Band) (linkADRReqParameters, bool, error) {
	if dev.GetMulticast() || dev.GetMacState() == nil {
		return linkADRReqParameters{}, false, nil
	}
	macState := dev.MacState
	desiredParameters, currentParameters := macState.DesiredParameters, macState.CurrentParameters
	if len(desiredParameters.Channels) > int(phy.MaxUplinkChannels) {
		return linkADRReqParameters{}, false, internal.ErrCorruptedMACState.
			WithAttributes(
				"desired_channels_len", len(desiredParameters.Channels),
				"phy_max_uplink_channels", phy.MaxUplinkChannels,
			).
			WithCause(internal.ErrUnknownChannel)
	}
	if len(currentParameters.Channels) > int(phy.MaxUplinkChannels) {
		return linkADRReqParameters{}, false, internal.ErrCorruptedMACState.
			WithAttributes(
				"current_channels_len", len(currentParameters.Channels),
				"phy_max_uplink_channels", phy.MaxUplinkChannels,
			).
			WithCause(internal.ErrUnknownChannel)
	}

	currentChs := make([]bool, phy.MaxUplinkChannels)
	for i, ch := range currentParameters.Channels {
		currentChs[i] = ch.GetEnableUplink()
	}
	pendingChs := make([]bool, phy.MaxUplinkChannels)
	iteratePendingNewChannelReq(dev, func(req *ttnpb.MACCommand_NewChannelReq) bool {
		pendingChs[req.ChannelIndex] = true
		// NewChannelReq will automatically enable the channel if the frequency is not 0.
		currentChs[req.ChannelIndex] = req.Frequency != 0
		return true
	})
	desiredChs := make([]bool, phy.MaxUplinkChannels)
	for i, ch := range desiredParameters.Channels {
		isEnabled := ch.GetEnableUplink()
		if isEnabled && ch.UplinkFrequency == 0 {
			return linkADRReqParameters{}, false, internal.ErrCorruptedMACState.
				WithAttributes(
					"i", i,
					"enabled", isEnabled,
					"uplink_frequency", ch.UplinkFrequency,
				).
				WithCause(internal.ErrDownlinkChannel)
		}
		if i >= len(currentParameters.Channels) && !pendingChs[i] {
			// The channel is not yet part of the end device channels list, and it is not pending
			// registration by a NewChannelReq command. As such, we avoid trying to enable it via
			// the channel mask.
			continue
		}
		desiredChs[i] = isEnabled
	}

	equalMasks := band.EqualChMasks(currentChs, desiredChs)
	switch {
	case !equalMasks:
		// NOTE: LinkADRReq is scheduled regardless of ADR settings if channel mask is required,
		// which often is the case with ABP devices or when ChMask CFList is not supported/used.
	case desiredParameters.AdrNbTrans != currentParameters.AdrNbTrans,
		desiredParameters.AdrDataRateIndex != currentParameters.AdrDataRateIndex,
		desiredParameters.AdrTxPowerIndex != currentParameters.AdrTxPowerIndex:
	default:
		return linkADRReqParameters{}, false, nil
	}
	desiredMasks, err := phy.GenerateChMasks(currentChs, desiredChs)
	if err != nil {
		return linkADRReqParameters{}, false, err
	}
	if len(desiredMasks) > math.MaxUint16 {
		// Something is really wrong.
		return linkADRReqParameters{}, false, internal.ErrCorruptedMACState.
			WithAttributes(
				"len", len(desiredMasks),
			).
			WithCause(internal.ErrChannelMask)
	}

	var (
		drIdx      ttnpb.DataRateIndex
		txPowerIdx uint32
		nbTrans    uint32
	)
	minDataRateIndex, _, allowedDataRateIndices, ok := channelDataRateRange(desiredParameters.Channels...)
	if !ok {
		return linkADRReqParameters{}, false, internal.ErrCorruptedMACState.
			WithCause(internal.ErrChannelDataRateRange)
	}

	// We need to check if the data rate index is valid with respect to the desired channels even in situations
	// in which the data rate index does not change, as it may be invalid with respect to the desired channel
	// mask.
	if _, ok := allowedDataRateIndices[desiredParameters.AdrDataRateIndex]; !ok {
		return linkADRReqParameters{}, false, internal.ErrCorruptedMACState.
			WithAttributes(
				"current_adr_data_rate_index", currentParameters.AdrDataRateIndex,
				"desired_adr_data_rate_index", desiredParameters.AdrDataRateIndex,
			)
	}
	if desiredParameters.AdrTxPowerIndex > uint32(phy.MaxTxPowerIndex()) {
		return linkADRReqParameters{}, false, internal.ErrCorruptedMACState.
			WithAttributes(
				"current_adr_tx_power_index", currentParameters.AdrTxPowerIndex,
				"desired_adr_tx_power_index", desiredParameters.AdrTxPowerIndex,
				"phy_max_tx_power_index", phy.MaxTxPowerIndex(),
			)
	}

	drIdx = desiredParameters.AdrDataRateIndex
	txPowerIdx = desiredParameters.AdrTxPowerIndex
	nbTrans = desiredParameters.AdrNbTrans
	resetDRTXToCurrent := func() {
		drIdx = currentParameters.AdrDataRateIndex
		txPowerIdx = currentParameters.AdrTxPowerIndex
	}
	hasNoChangeADRIndices := macspec.HasNoChangeADRIndices(macState.LorawanVersion)
	switch {
	case !deviceRejectedADRDataRateIndex(dev, drIdx) && !deviceRejectedADRTXPowerIndex(dev, txPowerIdx):
		// Only send the desired DataRateIndex and TXPowerIndex if neither of them were rejected.

	case equalMasks && desiredParameters.AdrNbTrans == currentParameters.AdrNbTrans:
		log.FromContext(ctx).Debug("Either desired data rate index or TX power output index have been rejected and there are no channel mask and NbTrans changes desired, avoid enqueueing LinkADRReq")
		return linkADRReqParameters{}, false, nil

	case hasNoChangeADRIndices &&
		!deviceRejectedADRDataRateIndex(dev, noChangeDataRateIndex) &&
		!deviceRejectedADRTXPowerIndex(dev, noChangeTXPowerIndex):

	default:
		logger := log.FromContext(ctx).WithFields(log.Fields(
			"current_adr_nb_trans", currentParameters.AdrNbTrans,
			"desired_adr_nb_trans", desiredParameters.AdrNbTrans,
			"desired_mask_count", len(desiredMasks),
		))
		for deviceRejectedADRDataRateIndex(dev, drIdx) || deviceRejectedADRTXPowerIndex(dev, txPowerIdx) {
			if drIdx < minDataRateIndex {
				logger.Warn("Device desired data rate is under the minimum data rate for ADR. Avoiding data rate and TX power changes")
				resetDRTXToCurrent()
				break
			}
			// Since either data rate or TX power index (or both) were rejected by the device, undo the
			// desired ADR adjustments step-by-step until possibly fitting index pair is found.
			if drIdx == minDataRateIndex && (deviceRejectedADRDataRateIndex(dev, drIdx) || txPowerIdx == 0) {
				logger.Warn("Device rejected either all available data rate indexes or all available TX power output indexes. Avoiding data rate and TX power changes")
				resetDRTXToCurrent()
				break
			}
			for drIdx > minDataRateIndex && (deviceRejectedADRDataRateIndex(dev, drIdx) || txPowerIdx == 0 && deviceRejectedADRTXPowerIndex(dev, txPowerIdx)) {
				// Increase data rate until a non-rejected index is found.
				// Set TX power to maximum possible value.
				drIdx--
				txPowerIdx = uint32(phy.MaxTxPowerIndex())
			}
			for txPowerIdx > 0 && deviceRejectedADRTXPowerIndex(dev, txPowerIdx) {
				// Increase TX output power until a non-rejected index is found.
				txPowerIdx--
			}
		}
	}
	if _, ok := allowedDataRateIndices[drIdx]; !ok {
		return linkADRReqParameters{}, false, internal.ErrCorruptedMACState.
			WithAttributes(
				"current_adr_data_rate_index", currentParameters.AdrDataRateIndex,
				"desired_adr_data_rate_index", desiredParameters.AdrDataRateIndex,
				"adr_data_rate_index", drIdx,
			)
	}
	if hasNoChangeADRIndices {
		if drIdx == currentParameters.AdrDataRateIndex && !deviceRejectedADRDataRateIndex(dev, noChangeDataRateIndex) {
			drIdx = noChangeDataRateIndex
		}
		if txPowerIdx == currentParameters.AdrTxPowerIndex && !deviceRejectedADRTXPowerIndex(dev, noChangeTXPowerIndex) {
			txPowerIdx = noChangeTXPowerIndex
		}
	}
	return linkADRReqParameters{
		Masks:         desiredMasks,
		DataRateIndex: drIdx,
		TxPowerIndex:  txPowerIdx,
		NbTrans:       nbTrans,
	}, true, nil
}

func DeviceNeedsLinkADRReq(ctx context.Context, dev *ttnpb.EndDevice, phy *band.Band) bool {
	_, required, err := generateLinkADRReq(ctx, dev, phy)
	return err == nil && required
}

func EnqueueLinkADRReq(
	ctx context.Context, dev *ttnpb.EndDevice, maxDownLen, maxUpLen uint16, phy *band.Band,
) (EnqueueState, error) {
	params, required, err := generateLinkADRReq(ctx, dev, phy)
	if err != nil {
		return EnqueueState{
			MaxDownLen: maxDownLen,
			MaxUpLen:   maxUpLen,
		}, err
	}
	if !required {
		return EnqueueState{
			MaxDownLen: maxDownLen,
			MaxUpLen:   maxUpLen,
			Ok:         true,
		}, nil
	}

	var st EnqueueState
	macState := dev.MacState
	f := func(nDown, nUp uint16) ([]*ttnpb.MACCommand, uint16, events.Builders, bool) {
		if int(nDown) < len(params.Masks) {
			return nil, 0, nil, false
		}

		uplinksNeeded := uint16(len(params.Masks))
		if macspec.SingularLinkADRAns(macState.LorawanVersion) {
			uplinksNeeded = 1
		}
		if nUp < uplinksNeeded {
			return nil, 0, nil, false
		}
		evs := make(events.Builders, 0, len(params.Masks))
		cmds := make([]*ttnpb.MACCommand, 0, len(params.Masks))
		for i, m := range params.Masks {
			req := &ttnpb.MACCommand_LinkADRReq{
				DataRateIndex:      params.DataRateIndex,
				TxPowerIndex:       params.TxPowerIndex,
				NbTrans:            params.NbTrans,
				ChannelMaskControl: uint32(m.Cntl),
				ChannelMask:        params.Masks[i].Mask[:],
			}
			cmds = append(cmds, req.MACCommand())
			evs = append(evs, EvtEnqueueLinkADRRequest.With(events.WithData(req)))
			log.FromContext(ctx).WithFields(log.Fields(
				"data_rate_index", req.DataRateIndex,
				"nb_trans", req.NbTrans,
				"tx_power_index", req.TxPowerIndex,
				"channel_mask_control", req.ChannelMaskControl,
				"channel_mask", req.ChannelMask,
			)).Debug("Enqueued LinkADRReq")
		}
		return cmds, uplinksNeeded, evs, true
	}
	macState.PendingRequests, st = enqueueMACCommand(
		ttnpb.MACCommandIdentifier_CID_LINK_ADR, maxDownLen, maxUpLen, f, macState.PendingRequests...,
	)
	return st, nil
}

func HandleLinkADRAns(
	ctx context.Context,
	dev *ttnpb.EndDevice,
	pld *ttnpb.MACCommand_LinkADRAns,
	dupCount uint,
	fCntUp uint32,
	fps *frequencyplans.Store,
	adrEnabled bool,
) (events.Builders, error) {
	if pld == nil {
		return nil, ErrNoPayload.New()
	}
	macState := dev.MacState
	allowDuplicateLinkADRAns := macspec.AllowDuplicateLinkADRAns(macState.LorawanVersion)
	if !allowDuplicateLinkADRAns && dupCount != 0 {
		return nil, internal.ErrInvalidPayload.New()
	}

	ev := EvtReceiveLinkADRAccept
	rejected := false

	// LoRaWAN 1.0.4 spec L534-538:
	// An end-device SHOULD accept the channel mask controls present in LinkADRReq, even
	// when the ADR bit is not set. The end-device SHALL respond to all LinkADRReq commands
	// with a LinkADRAns indicating which command elements were accepted and which were
	// rejected. This behavior differs from when the uplink ADR bit is set, in which case the end-
	// device accepts or rejects the entire command.
	if macspec.UseADRBit(macState.LorawanVersion) {
		rejected = !pld.ChannelMaskAck ||
			(adrEnabled && !pld.DataRateIndexAck) ||
			(adrEnabled && !pld.TxPowerIndexAck)
	} else {
		rejected = !pld.ChannelMaskAck || !pld.DataRateIndexAck || !pld.TxPowerIndexAck
	}

	if rejected {
		ev = EvtReceiveLinkADRReject

		// See "Table 6: LinkADRAns status bits signification" of LoRaWAN 1.1 specification
		if !pld.ChannelMaskAck {
			log.FromContext(ctx).Warn("Either Network Server sent a channel mask, which enables a yet undefined channel or requires all channels to be disabled, or device is malfunctioning.")
		}
	}
	evs := events.Builders{ev.With(events.WithData(pld))}

	phy, err := internal.DeviceBand(dev, fps)
	if err != nil {
		return evs, err
	}

	handler := handleMACResponseBlock
	if !allowDuplicateLinkADRAns && !macspec.SingularLinkADRAns(macState.LorawanVersion) {
		handler = handleMACResponse
	}
	var n uint
	var req *ttnpb.MACCommand_LinkADRReq
	currentParameters := macState.CurrentParameters
	macState.PendingRequests, err = handler(
		ttnpb.MACCommandIdentifier_CID_LINK_ADR,
		false,
		func(cmd *ttnpb.MACCommand) error {
			if allowDuplicateLinkADRAns && n > dupCount+1 {
				return internal.ErrInvalidPayload.New()
			}
			n++

			req = cmd.GetLinkAdrReq()
			if req.NbTrans > 15 || len(req.ChannelMask) != 16 || req.ChannelMaskControl > 7 {
				panic("Network Server scheduled an invalid LinkADR command")
			}
			if !pld.ChannelMaskAck || !pld.DataRateIndexAck || !pld.TxPowerIndexAck {
				return nil
			}
			var mask [16]bool
			copy(mask[:], req.ChannelMask)
			m, err := phy.ParseChMask(mask, uint8(req.ChannelMaskControl))
			if err != nil {
				return err
			}
			for i, masked := range m {
				if int(i) >= len(currentParameters.Channels) || currentParameters.Channels[i] == nil {
					if !masked {
						continue
					}
					return internal.ErrCorruptedMACState.
						WithAttributes(
							"i", i,
							"channels_len", len(currentParameters.Channels),
						).
						WithCause(internal.ErrUnknownChannel)
				}
				currentParameters.Channels[i].EnableUplink = masked
			}
			return nil
		},
		macState.PendingRequests...,
	)
	if err != nil || req == nil {
		return evs, err
	}

	if !pld.DataRateIndexAck {
		i := searchDataRateIndex(req.DataRateIndex, macState.RejectedAdrDataRateIndexes...)
		if i == len(macState.RejectedAdrDataRateIndexes) ||
			macState.RejectedAdrDataRateIndexes[i] != req.DataRateIndex {
			macState.RejectedAdrDataRateIndexes = append(
				macState.RejectedAdrDataRateIndexes, ttnpb.DataRateIndex_DATA_RATE_0,
			)
			copy(macState.RejectedAdrDataRateIndexes[i+1:], macState.RejectedAdrDataRateIndexes[i:])
			macState.RejectedAdrDataRateIndexes[i] = req.DataRateIndex
		}
	}
	if !pld.TxPowerIndexAck {
		i := searchUint32(req.TxPowerIndex, macState.RejectedAdrTxPowerIndexes...)
		if i == len(macState.RejectedAdrTxPowerIndexes) ||
			macState.RejectedAdrTxPowerIndexes[i] != req.TxPowerIndex {
			macState.RejectedAdrTxPowerIndexes = append(macState.RejectedAdrTxPowerIndexes, 0)
			copy(macState.RejectedAdrTxPowerIndexes[i+1:], macState.RejectedAdrTxPowerIndexes[i:])
			macState.RejectedAdrTxPowerIndexes[i] = req.TxPowerIndex
		}
	}
	if !pld.ChannelMaskAck || !pld.DataRateIndexAck || !pld.TxPowerIndexAck {
		return evs, nil
	}
	hasNoChangeADRIndices := macspec.HasNoChangeADRIndices(macState.LorawanVersion)
	if !hasNoChangeADRIndices || req.DataRateIndex != noChangeDataRateIndex {
		currentParameters.AdrDataRateIndex = req.DataRateIndex
		macState.LastAdrChangeFCntUp = fCntUp
	}
	if !hasNoChangeADRIndices || req.TxPowerIndex != noChangeTXPowerIndex {
		currentParameters.AdrTxPowerIndex = req.TxPowerIndex
		macState.LastAdrChangeFCntUp = fCntUp
	}
	if req.NbTrans > 0 && currentParameters.AdrNbTrans != req.NbTrans {
		currentParameters.AdrNbTrans = req.NbTrans
		macState.LastAdrChangeFCntUp = fCntUp
	}
	return evs, nil
}
