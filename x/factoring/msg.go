package factoring

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// verify interface at compile time
var _ sdk.Msg = &MsgFactorInvoice{}
var _ sdk.Msg = &MsgAddAdmin{}
var _ sdk.Msg = &MsgRemoveAdmin{}
var _ sdk.Msg = &MsgUpdateParams{}

const (
	TypeMsgFactorInvoice = "factor_invoice"
	TypeMsgAddAdmin       = "add_admin"
	TypeMsgRemoveAdmin    = "remove_admin"
	TypeMsgUpdateParams   = "update_params"
)

// MsgFactorInvoice msg for factoring an invoice.
type MsgFactorInvoice struct {
	InvoiceID   uint64         `json:"invoice_id"`
	Amount      uint64         `json:"amount"`
	FactorType  FactorType      `json:"factor_type"`
	Factor      sdk.AccAddress `json:"factor"`
}

// NewMsgFactorInvoice returns a new factor invoice message.
func NewMsgFactorInvoice(factor sdk.AccAddress, invoiceID uint64, amount, factorType FactorType) MsgFactorInvoice {
	return MsgFactorInvoice{
		InvoiceID:   invoiceID,
		Amount:   amount,
		FactorType: factorType,
		Factor:   factor,
	}
}
func (MsgFactorInvoice) Route() string {
	return RouterKey
}

func (MsgFactorInvoice) Type() string {
	return TypeMsgFactorInvoice
}

func (msg MsgFactorInvoice) ValidateBasic() sdk.Error {
	if !msg.FactorType.ValidForArgument() {
		return ErrCodeInvalidFactorType(msg.FactorType)
	}

	if len(msg.Factor) == 0 {
		return sdk.ErrInvalidAddress("Must provide a valid address")
	}

	if len(msg.Amount) == 0 {
		return ErrCodeInvalidAmount()
	}
	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgFactorInvoice) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgFactorInvoice) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Factor}
}


// MsgAddAdmin defines the message to add a new admin
type MsgAddAdmin struct {
	Admin   sdk.AccAddress `json:"admin"`
	Creator sdk.AccAddress `json:"creator"`
}

// NewMsgAddAdmin returns the messages to add a new admin
func NewMsgAddAdmin(admin, creator sdk.AccAddress) MsgAddAdmin {
	return MsgAddAdmin{
		Admin:   admin,
		Creator: creator,
	}
}

// ValidateBasic implements Msg
func (msg MsgAddAdmin) ValidateBasic() sdk.Error {
	if len(msg.Admin) == 0 {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid address: %s", msg.Admin.String()))
	}

	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid address: %s", msg.Creator.String()))
	}

	return nil
}

// Route implements Msg
func (msg MsgAddAdmin) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgAddAdmin) Type() string { return TypeMsgAddAdmin }

// GetSignBytes implements Msg
func (msg MsgAddAdmin) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners implements Msg. Returns the creator as the signer.
func (msg MsgAddAdmin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// MsgRemoveAdmin defines the message to remove an admin
type MsgRemoveAdmin struct {
	Admin   sdk.AccAddress `json:"admin"`
	Remover sdk.AccAddress `json:"remover"`
}

// NewMsgRemoveAdmin returns the messages to remove an admin
func NewMsgRemoveAdmin(admin, remover sdk.AccAddress) MsgRemoveAdmin {
	return MsgRemoveAdmin{
		Admin:   admin,
		Remover: remover,
	}
}

// ValidateBasic implements Msg
func (msg MsgRemoveAdmin) ValidateBasic() sdk.Error {
	if len(msg.Admin) == 0 {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid address: %s", msg.Admin.String()))
	}

	if len(msg.Remover) == 0 {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid address: %s", msg.Remover.String()))
	}

	return nil
}

// Route implements Msg
func (msg MsgRemoveAdmin) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgRemoveAdmin) Type() string { return TypeMsgRemoveAdmin }

// GetSignBytes implements Msg
func (msg MsgRemoveAdmin) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners implements Msg. Returns the remover as the signer.
func (msg MsgRemoveAdmin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Remover)}
}

// MsgUpdateParams defines the message to remove an admin
type MsgUpdateParams struct {
	Updates       Params         `json:"updates"`
	UpdatedFields []string       `json:"updated_fields"`
	Updater       sdk.AccAddress `json:"updater"`
}

// NewMsgUpdateParams returns the message to update the params
func NewMsgUpdateParams(updates Params, updatedFields []string, updater sdk.AccAddress) MsgUpdateParams {
	return MsgUpdateParams{
		Updates:       updates,
		UpdatedFields: updatedFields,
		Updater:       updater,
	}
}

// ValidateBasic implements Msg
func (msg MsgUpdateParams) ValidateBasic() sdk.Error {
	return nil
}

// Route implements Msg
func (msg MsgUpdateParams) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgUpdateParams) Type() string { return TypeMsgUpdateParams }

// GetSignBytes implements Msg
func (msg MsgUpdateParams) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners implements Msg. Returns the remover as the signer.
func (msg MsgUpdateParams) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Updater)}
}
© 2020 GitHub, Inc.
Terms
Privacy
Security
Status
Help
Contact GitHub
Pricing
API
Training
Blog
About
