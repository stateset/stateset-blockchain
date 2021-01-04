package loan

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// TypeMsgCreateLoan represents the type of the message for creating new loan
	TypeMsgCreateLoan = "create_loan"
	// TypeMsgCreateLoan represents the type of the message for creating new loan
	TypeMsgEditLoan = "edit_loan"
	// TypeMsgDeleteLoan represents the type of the message for creating new loan
	TypeMsgDeleteLoan = "delete_loan"
	// TypeMsgPaybackLoan represents the type of the message for creating new loan
	TypeMsgCreateLoan = "payback_loan"
	// TypeMsgAddAdmin represents the type of message for adding a new admin
	TypeMsgAddAdmin = "add_admin"
	// TypeMsgRemoveAdmin represents the type of message for removing an admin
	TypeMsgRemoveAdmin = "remove_admin"
	// TypeMsgUpdateParams represents the type of
	TypeMsgUpdateParams = "update_params"
)

// verify interface at compile time
var _ sdk.Msg = &MsgCreateLoan{}
var _ sdk.Msg = &MsgEditLoan{}
var _ sdk.Msg = &MsgDeleteLoan{}
var _ sdk.Msg = &MsgPayBackLoan{}
var _ sdk.Msg = &MsgAddAdmin{}
var _ sdk.Msg = &MsgRemoveAdmin{}
var _ sdk.Msg = &MsgUpdateParams{}

// MsgCreateLoan defines a message to submit an loan
type MsgCreateLoan struct {
	MarketID string             `json:"market_id"`
	InvoiceID 	  string 			 `json:"invoice_id"`
	Body          string         	 `json:"body"`
	Lender        sdk.AccAddress     `json:"lender"`
	Source        string             `json:"source,omitempty"`
}

// NewMsgCreateloan creates a new message to create a oan
func NewMsgCreateLoan(marketID, invoiceID, body string, lender sdk.AccAddress, source string) MsgCreateLoan {
	return MsgCreateLoan {
		MarketID: marketID,
		InvoiceID:    invoiceID,
		Body:        body,
		Lender:     lender,
		Source:      source,
	}
}

// Route is the name of the route for loan
func (msg MsgCreateLoan) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgCreateLoan) Type() string {
	return TypeMsgCreateLoan
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgCreateLoan) ValidateBasic() sdk.Error {
	if len(msg.Body) == 0 {
		return ErrInvalidBodyTooShort(msg.Body)
	}
	if len(msg.MarketID) == 0 {
		return ErrInvalidMarketID(msg.MarketID)
	}
	if len(msg.InvoiceID) == 0 {
		return ErrINvalidInvoiceID(msg.InvoiceID)
	}
	if len(msg.Lender) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Lender.String())
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgCreateLoan) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgCreateLoan) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Lender)}
}

// MsgDeleteLoan defines a message to delete a loan
type MsgDeleteLoan struct {
	ID      uint64         `json:"id"`
	Lender sdk.AccAddress `json:"lender"`
}

// Route is the name of the route for loan
func (msg MsgDeleteLoan) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgDeleteLoan) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgDeleteLoan) ValidateBasic() sdk.Error {
	if msg.ID == 0 {
		return ErrUnknownLoan(msg.ID)
	}
	if len(msg.Lender) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Lender.String())
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgDeleteLoan) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgDeleteLoan) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// MsgEditLoan defines a message to submit an loan
type MsgEditLoan struct {
	ID     uint64         `json:"id"`
	Body   string         `json:"body"`
	Editor sdk.AccAddress `json:"editor"`
}

// NewMsgEditLoan creates a new message to edit a loan
func NewMsgEditLoan(id uint64, body string, editor sdk.AccAddress) MsgEditLoan {
	return MsgEditLoan{
		ID:     id,
		Body:   body,
		Editor: editor,
	}
}

// Route is the name of the route for loan
func (msg MsgEditLoan) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgEditLoan) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgEditAccount) ValidateBasic() sdk.Error {
	if msg.ID == 0 {
		return ErrUnknownLoan(msg.ID)
	}
	if len(msg.Editor) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Editor.String())
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgEditLoan) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgEditLoan) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Editor)}
}

// MsgAddAdmin defines the message to add a new admin
type MsgAddAdmin struct {
	Admin   sdk.AccAddress `json:"admin"`
	Lender sdk.AccAddress `json:"lender"`
}

// NewMsgAddAdmin returns the messages to add a new admin
func NewMsgAddAdmin(admin, lender sdk.AccAddress) MsgAddAdmin {
	return MsgAddAdmin{
		Admin:   admin,
		Lender: lender,
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