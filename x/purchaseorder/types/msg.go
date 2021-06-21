package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// TypeMsgCreatePurchaseOrder represents the type of the message for creating new purchaseorder
	TypeMsgCreatePurchaseOrder = "create_purchaseorder"
	// TypeMsgEditAgreement represents the type of the message for editing an purchaseorder
	TypeMsgEditPurchaseOrder = "edit_purchasorder"
	// TypeMsgDeletePurchaseOrder represents the type of the message for activating an purchaseorder
	TypeMsgDeletePurchaseOrder = "delete_purchasorder"
	// TypeMsgEditPurchaseOrder represents the type of the message for activating an purchaseorder
	TypeMsgCompletePurchaseOrder = "complete_purchaseorder"
	// TypeMsgCompletePurchaseOrder represents the type of the message for amending an purchaseorder
	TypeMsgCancelPurchaseOrder = "cancel_purchaseorder"
	// TypeMsgCancelPurchaseOrder represents the type of the message for amending an purchaseorder
	TypeMsgFinancePurchaseOrder = "finance_purchaseorder"
	// TypeMsgFinancePurchaseOrder represents the type of the message for renewing an purchaseorder
	TypeMsgAddAdmin = "add_admin"
	// TypeMsgRemoveAdmin represents the type of message for removing an admin
	TypeMsgRemoveAdmin = "remove_admin"
	// TypeMsgUpdateParams represents the type of
	TypeMsgUpdateParams = "update_params"
)

// verify interface at compile time
var _ sdk.Msg = &MsgCreatePurchaseOrder{}
var _ sdk.Msg = &MsgEditPurchaseOrder{}
var _ sdk.Msg = &MsgDeletePurchaseOrder{}
var _ sdk.Msg = &MsgCompletePurchaseOrder{}
var _ sdk.Msg = &MsgCancelPurchaseOrder{}
var _ sdk.Msg = &MsgFinancePurchaseOrder{}
var _ sdk.Msg = &MsgRemoveAdmin{}
var _ sdk.Msg = &MsgUpdateParams{}

// MsgCreatePurchaseOrder defines a message to create an purchaseorder
type MsgCreatePurchaseOrder struct {
	PurchaseOrderID 	  string 	`json:"purchaseorder_id"`
	Lender        sdk.AccAddress     `json:"counterparty"`
}

// NewMsgCreatePurchaseOrder creates a new message to create an purchaseorder
func NewMsgCreatePurchaseOrder(purchaseOrderID string, lender sdk.AccAddress,) MsgCreatePurchaseOrder {
	return MsgCreatePurchaseOrder {
		PurchaseOrderID:    purchaseOrderID,
		Lender: lender,
	}
}

// Route is the name of the route
func (msg MsgCreatePurchaseOrder) Route() string { return RouterKey }

// Type is the name for the Msg
func (msg MsgCreatePurchaseOrder) Type() string { return TypeMsgCreatePurchaseOrder }

// ValidateBasic validates basic fields of the Msg
func (msg MsgCreatePurchaseOrder) ValidateBasic() sdkerrors {
	if len(msg.Description) == 0 {
		return ErrInvalidDescrtiptionTooShort(msg.Description)
	}
	if len(msg.PurchaseOrderID) == 0 {
		return ErrInvalidPurchaseOrderID(msg.PurchaseOrderID)
	}
	if len(msg.Counterparty) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Counterparty.String())
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
type MsgEditPurchaseOrder struct {
	ID      uint64         `json:"id"`
	Counterparty sdk.AccAddress `json:"counterparty"`
}


// Complete Purchase Order

// Msg Purchase Order defines a message to activate an Purchase Order
type MsgCompletePurchaseOrder struct {
	ID              uint64         `json:"id"`
	PurchaseOrderStatus string         `json:"purcahseOrderStatus"`
	Counterparty    sdk.AccAddress `json:"counterparty"`
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
	if msg.AgreementID == 0 {
		return ErrUnknownAgreement(msg.AgreementID)
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
	ID              uint64         `json:"id"`
	PurchaseOrderStatus string         `json:"purchaseOrderStatus"`
	Counterparty    sdk.AccAddress `json:"counterparty"`
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
	if msg.AgreementID == 0 {
		return ErrUnknownAgreement(msg.AgreementID)
	}
	if len(msg.Counterparty) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Counterparty.String())
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