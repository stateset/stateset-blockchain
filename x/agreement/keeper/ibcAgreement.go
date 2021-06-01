package keeper

import (
	"net/url"
	"time"

	app "github.com/stateset/stateset-blockchain/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/gaskv"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	log "github.com/tendermint/tendermint/libs/log"
)

var packet types.IbcAgeementPacketData

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

func (k Keeper) OnRecvIbcAgreementPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcAgreementPacketData) (packetAck types.IbcAgreementPacketAck, err error) {
    // validate packet data upon receiving
    if err := data.ValidateBasic(); err != nil {
        return packetAck, err
    }

    id := k.AppendAgreement(
        ctx,
        types.Agreement{
            Creator: packet.SourcePort+"-"+packet.SourceChannel+"-"+data.Creator,
            AgreementName: data.AgreementName,
			AgreementNumber: data.AgreementNumber,
			AgreementType: data.agreementType,
			AgreementStatus: data.agreementStatus,
			TotalAgreementValue: data.totalAgreementValue,
			Party: data.party,
			Counterparty: data.counterparty
			AgreementStartBlock: data.agreementStartBlock,
			AgreementEndBlock: data.agreementEndBlock
        },
    )
    packetAck.AgreementID = strconv.FormatUint(id, 10)

    return packetAck, nil
}


// OnAcknowledgementIbcAgreementPacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementIbcAgreementPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcAgreementPacketData, ack channeltypes.Acknowledgement) error {
    switch dispatchedAck := ack.Response.(type) {
    case *channeltypes.Acknowledgement_Error:
        // We will not treat acknowledgment error in this tutorial
        return nil
    case *channeltypes.Acknowledgement_Result:
        // Decode the packet acknowledgment
        var packetAck types.IbcAgreementPacketAck
        
        if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
            // The counter-party module doesn't implement the correct acknowledgment format
            return errors.New("cannot unmarshal acknowledgment")
        }

        k.AppendSentAgreement(
            ctx,
            types.SentAgreement{
                Creator: data.Creator,
                AgreementID: packetAck.AgreementID,
				AgreementName: data.AgreementName,
				AgreementNumber: data.AgreementNumber,
				AgreementType: data.agreementType,
				AgreementStatus: data.agreementStatus,
				TotalAgreementValue: data.totalAgreementValue,
				Party: data.party,
				Counterparty: data.counterparty
				AgreementStartBlock: data.agreementStartBlock,
				AgreementEndBlock: data.agreementEndBlock
                Chain: packet.DestinationPort+"-"+packet.DestinationChannel,
            }, 
        )


        return nil
    default:
        // The counter-party module doesn't implement the correct acknowledgment format
        return errors.New("invalid acknowledgment format")
    }
}


// OnTimeoutIbcAgreementPacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutIbcAgreementPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcAgreementPacketData) error {
    k.AppendTimedoutPost(
        ctx,
        types.TimedoutPost{
            Creator: data.Creator,
			AgreementName: data.AgreementName,
			AgreementNumber: data.AgreementNumber,
			AgreementType: data.agreementType,
			AgreementStatus: data.agreementStatus,
			TotalAgreementValue: data.totalAgreementValue,
			Party: data.party,
			Counterparty: data.counterparty
			AgreementStartBlock: data.agreementStartBlock,
			AgreementEndBlock: data.agreementEndBlock
            Chain: packet.DestinationPort+"-"+packet.DestinationChannel,
        },
    )

    return nil
}