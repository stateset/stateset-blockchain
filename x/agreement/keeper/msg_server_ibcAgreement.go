package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/stateset/stateset-blockchain/x/agreement/types"
)

func (k msgServer) SendIbcAgreement(goCtx context.Context, msg *types.MsgSendIbcAgreement) (*types.MsgSendIbcAgreementResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Construct the packet
	var packet types.IbcAgreementPacketData

	packet.AgreementID = msg.agreementID
	packet.AgreementNumber = msg.agreementNumber
	packet.AgreementName = msg.agreementName
	packet.AgreementType = msg.agreementType
	packet.AgreementStatus = msg.agreementStatus
	packet.TotalAgreementValue = msg.totalAgreementValue
	packet.Party = msg.party
	packet.Counterparty = msg.counterparty
	packet.AgreementStartBlock = msg.agreementStartBlock
	packet.AgreementEndBlock = msg.agreementEndBlock

	// Transmit the packet
	err := k.TransmitIbcAgreementPacket(
		ctx,
		packet,
		msg.Port,
		msg.ChannelID,
		clienttypes.ZeroHeight(),
		msg.TimeoutTimestamp,
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendIbcAgreementResponse{}, nil
}
