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

var packet types.IbcPurchaseOrderPacketData

packet.PurchaseOrderID = msg.purchaseOrderID
packet.PurchaseOrderNumber = msg.purchaseOrderNumber
packet.PurchaseOrderName = msg.purchaseOrderName
packet.PurchaseOrderStatus = msg.purchaseOrderStatus
packet.Description = msg.description
packet.PurchaseDate = msg.purchaseDate
packet.DeliveryDate = msg.deliveryDate
packet.Subtotal = msg.subtotal
packet.Total = msg.total
packet.Purchaser = msg.purchaser
packet.Vendor = msg.vendor
packet.Fulfiller = msg.fulfiller
packet.Financer = msg.financer

// Transmit the packet
err := k.TransmitIbcPurchaseOrderPacket(
	ctx,
	packet,
	msg.Port,
	msg.ChannelID,
	clienttypes.ZeroHeight(),
	msg.TimeoutTimestamp,
)

func (k Keeper) OnRecvIbcPurchaseOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPurchaseOrderPacketData) (packetAck types.IbcPurchaseOrderPacketAck, err error) {
    // validate packet data upon receiving
    if err := data.ValidateBasic(); err != nil {
        return packetAck, err
    }

    id := k.AppendPurchaseOrder(
        ctx,
        types.PurchaseOrder{
            Creator: packet.SourcePort+"-"+packet.SourceChannel+"-"+data.Creator,
            PurchaseOrderID:       data.PurchaseOrderId,
            PurchaseOrderNumber:   data.PurchaseOrderNumber,
            PurchaseOrderName:     data.PurchaseOrderName,
            Description:   data.Description,
            PurchaseDate: data.PurchaseDate,
            DeliveryDate:   data.DeliveryDate,
            Subtotal:        data.Subtotal,
            Total: 			 data.Total,
            Purchaser:		     data.Purchaser,
            Vendor:	 data.Vendor,
            Financer:	 data.Financer,
        },
    )
    packetAck.PurchaseOrderID = strconv.FormatUint(id, 10)

    return packetAck, nil
}


// OnAcknowledgementIbcPurchaseOrderPacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementIbcPurchaseOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPurchaseOrderPacketData, ack channeltypes.Acknowledgement) error {
    switch dispatchedAck := ack.Response.(type) {
    case *channeltypes.Acknowledgement_Error:
        // We will not treat acknowledgment error in this tutorial
        return nil
    case *channeltypes.Acknowledgement_Result:
        // Decode the packet acknowledgment
        var packetAck types.IbcPurchaseOrderPacketAck
        
        if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
            // The counter-party module doesn't implement the correct acknowledgment format
            return errors.New("cannot unmarshal acknowledgment")
        }

        k.AppendSentPurchaseOrder(
            ctx,
            types.SentPurchaseOrder{
                Creator: data.Creator,
                PurchaseOrderID:       data.PurchaseOrderId,
                PurchaseOrderNumber:   data.PurchaseOrderNumber,
                PurchaseOrderName:     data.PurchaseOrderName,
                Description:   data.Description,
                PurchaseDate: data.PurchaseDate,
                DeliveryDate:   data.DeliveryDate,
                Subtotal:        data.Subtotal,
                Total: 			 data.Total,
                Purchaser:		     data.Purchaser,
                Vendor:	 data.Vendor,
                Financer:	 data.Financer,
                Chain: packet.DestinationPort+"-"+packet.DestinationChannel,
            }, 
        )


        return nil
    default:
        // The counter-party module doesn't implement the correct acknowledgment format
        return errors.New("invalid acknowledgment format")
    }
}


// OnTimeoutIbcPurchaseOrderPacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutIbcPurchaseOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPurchaseOrderPacketData) error {
    k.AppendTimedoutPost(
        ctx,
        types.TimedoutPost{
            Creator: data.Creator,
            PurchaseOrderID:       data.PurchaseOrderId,
            PurchaseOrderNumber:   data.PurchaseOrderNumber,
            PurchaseOrderName:     data.PurchaseOrderName,
            Description:   data.Description,
            PurchaseDate: data.PurchaseDate,
            DeliveryDate:   data.DeliveryDate,
            Subtotal:        data.Subtotal,
            Total: 			 data.Total,
            Purchaser:		     data.Purchaser,
            Vendor:	 data.Vendor,
            Financer:	 data.Financer,
            Chain: packet.DestinationPort+"-"+packet.DestinationChannel,
        },
    )

    return nil
}