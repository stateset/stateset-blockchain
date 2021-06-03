package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/stateset/stateset-blockchain/x/purchaseorder/types"
)

func (k msgServer) SendIbcPurchaseOrder(goCtx context.Context, msg *types.MsgSendIbcPurchaseOrder) (*types.MsgSendIbcPurchaseOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

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
	if err != nil {
		return nil, err
	}

	return &types.MsgSendIbcPurchaseOrderResponse{}, nil
}
