package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSendIbcPurchaseOrder{}

func NewMsgSendIbcPurchaseOrder(
	sender string,
	port string,
	channelID string,
	timeoutTimestamp uint64,
	purchaseordername string,
	purchaseordernumber string,
	status string,
	total string,
) *MsgSendIbcPurchaseOrder {
	return &MsgSendIbcPurchaseOrder{
		Sender:              sender,
		Port:                port,
		ChannelID:           channelID,
		TimeoutTimestamp:    timeoutTimestamp,
		Purchaseordername:   purchaseordername,
		Purchaseordernumber: purchaseordernumber,
		Status:              status,
		Total:               total,
	}
}

func (msg *MsgSendIbcPurchaseOrder) Route() string {
	return RouterKey
}

func (msg *MsgSendIbcPurchaseOrder) Type() string {
	return "SendIbcPurchaseOrder"
}

func (msg *MsgSendIbcPurchaseOrder) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgSendIbcPurchaseOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSendIbcPurchaseOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
