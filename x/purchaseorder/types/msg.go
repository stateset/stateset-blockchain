package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// TypeMsgCreatePurchaseOrder represents the type of the message for creating new purchaseorder
	TypeMsgCreatePurchaseOrder = "create_purchaseorder"
	// TypeMsgUpdateAgreement represents the type of the message for updating an purchaseorder
	TypeMsgUpdatePurchaseOrder = "update_purchaseorder"
	// TypeMsgDeletePurchaseOrder represents the type of the message for activating an purchaseorder
	TypeMsgDeletePurchaseOrder = "delete_purchaseorder"
	// TypeMsgEditPurchaseOrder represents the type of the message for activating an purchaseorder
	TypeMsgCompletePurchaseOrder = "complete_purchaseorder"
	// TypeMsgCompletePurchaseOrder represents the type of the message for amending an purchaseorder
	TypeMsgCancelPurchaseOrder = "cancel_purchaseorder"
	// TypeMsgCancelPurchaseOrder represents the type of the message for amending an purchaseorder
	TypeMsgLockPurchaseOrder = "lock_purchaseorder"
	// TypeMsgFinancePurchaseOrder represents the type of the message for renewing an purchaseorder
	TypeMsgFinancePurchaseOrder = "finance_purchaseorder"
	// TypeMsgFinancePurchaseOrder represents the type of the message for renewing an purchaseorder
	TypeMsgAddAdmin = "add_admin"
	// TypeMsgRemoveAdmin represents the type of message for removing an admin
	TypeMsgRemoveAdmin = "remove_admin"
	// TypeMsgUpdateParams represents the type of
	TypeMsgUpdateParams = "update_params"
	// TypeMsgUpdateParams represents the type of
	TypeMsgSendIbcPurchaseOrder = "send_ibc_purchaseorder"
)

// verify interface at compile time
var _ sdk.Msg = &MsgCreatePurchaseOrder{}
var _ sdk.Msg = &MsgEditPurchaseOrder{}
var _ sdk.Msg = &MsgDeletePurchaseOrder{}
var _ sdk.Msg = &MsgCompletePurchaseOrder{}
var _ sdk.Msg = &MsgCancelPurchaseOrder{}
var _ sdk.Msg = &MsgFinancePurchaseOrder{}
var _ sdk.Msg = &MsgLockPurchaseOrder{}
var _ sdk.Msg = &MsgRemoveAdmin{}
var _ sdk.Msg = &MsgUpdateParams{}
var _ sdk.Msg = &MsgSendIbcPurchaseOrder{}

// MsgCreatePurchaseOrder defines a message to create an purchaseorder
type MsgCreatePurchaseOrder struct {
	PurchaseOrderID 	  string 	`json:"purchaseorder_id"`
	Creator       sdk.AccAddress     `json:"creator"`
}

// NewMsgCreatePurchaseOrder creates a new message to create an purchaseorder
func NewMsgCreatePurchaseOrder(purchaseOrderID string) MsgCreatePurchaseOrder {
	return MsgCreatePurchaseOrder {
		PurchaseOrderID:    purchaseOrderID,
		Creator: creator,
	}
}

// Route is the name of the route
func (msg MsgCreatePurchaseOrder) Route() string { return RouterKey }

// Type is the name for the Msg
func (msg MsgCreatePurchaseOrder) Type() string { return TypeMsgCreatePurchaseOrder }

// ValidateBasic validates basic fields of the Msg
func (msg MsgCreatePurchaseOrder) ValidateBasic() sdkerrors {

	if len(msg.PurchaseOrderID) == 0 {
		return ErrInvalidPurchaseOrderID(msg.PurchaseOrderID)
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgCreatePurchaseOrder) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgCreatePurchaseOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Counterparty)}
}


// Edit Purchase Order

// MsgAmendAgreement defines a message to amend an purchaseorder
type MsgUpdatePurchaseOrder struct {
	PurchaseOrderID      uint64         `json:"id"`
}


// Complete Purchase Order

// Msg Purchase Order defines a message to activate an Purchase Order
type MsgCompletePurchaseOrder struct {
	PurchaseOrderID            uint64         `json:"id"`
	PurchaseOrderStatus string         `json:"purcahseOrderStatus"`
}

// Route is the name of the route for an purchaseorder
func (msg MsgCompletePurchaseOrder) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgCompletePurchaseOrder) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgCompletePurchaseOrder) ValidateBasic() sdkerrors {
	if msg.PurchaseOrderID == 0 {
		return ErrUnknownPurchaseOrder(msg.PurchaseOrderID)
	}
	if len(msg.Counterparty) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Counterparty.String())
	}

	if msg.PurchaseOrderStatus != "COMPLETED" {
		return Error("The Purchase Order status must be Completed.")
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgCompletePurchaseOrder) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgCompletePurchaseOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Party)}
}

// Cancel Purchase Order

// Msg Purchase Order defines a message to cancel an Purchase Order
type MsgCancelPurchaseOrder struct {
	PurchaseOrderID              uint64         `json:"purchaseorder_id"`
	PurchaseOrderStatus string         `json:"purchaseOrderStatus"`
}

// Route is the name of the route for an purchaseorder
func (msg MsgCancelPurchaseOrder) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgCancelPurchaseOrder) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgCancelPurchaseOrder) ValidateBasic() sdk.Error {
	if msg.PurchaseOrderID == 0 {
		return ErrUnknownPurchaseOrder(msg.PurchaseOrderID)
	}

	if msg.PurchaseOrderStatus != "CANCELLED" {
		return Error("The Purchase Order status must be Cancelled.")
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgCancelPurchaseOrder) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgCancelPurchaseOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Party)}
}



// Finance Purchase Order

// Msg Purchase Order defines a message to activate an Purchase Order
type MsgFinancePurchaseOrder struct {
	PurchaseOrderID            uint64         `json:"purchaseorder_id"`
	PurchaseOrderStatus string         `json:"purcahseOrderStatus"`
}

// Route is the name of the route for an purchaseorder
func (msg MsgFinancePurchaseOrder) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgFinancePurchaseOrder) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgFinancePurchaseOrder) ValidateBasic() sdkerrors {
	if msg.PurchaseOrderID == 0 {
		return ErrUnknownPurchaseOrder(msg.PurchaseOrderID)
	}

	if msg.PurchaseOrderStatus != "FINANCED" {
		return Error("The Purchase Order status must be Financed.")
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgFinancePurchaseOrder) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgFinancePurchaseOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Party)}
}




// Send IBC Purchase Order

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
