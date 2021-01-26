package invoice

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// TypeMsgCreateInvoice represents the type of the message for creating new invoice
	TypeMsgCreateInvoice = "create_invoice"
	// TypeMsgEditInvoice represents the type of the message for editing an invoice
	TypeMsgCancelInvoice = "cancel_invoice"
	// TypeMsgEditInvoice represents the type of the message for editing an invoice
	TypeMsgEditInvoice = "edit_invoice"
	// TypeMsgFactorInvoice represents the type of the message for factoring an invoice
	TypeMsgFactorInvoice = "factor_invoice"
	// TypeMsgPayInvoice represents the type of the message for paying an invoice
	TypeMsgPayInvoice = "pay_invoice"
	// TypeMsgAddAdmin represents the type of message for adding a new admin
	TypeMsgAddAdmin = "add_admin"
	// TypeMsgRemoveAdmin represents the type of message for removeing an admin
	TypeMsgRemoveAdmin = "remove_admin"
	// TypeMsgUpdateParams represents the type of
	TypeMsgUpdateParams = "update_params"
)

// verify interface at compile time
var _ sdk.Msg = &MsgCreateInvoice{}
var _ sdk.Msg = &MsgCancelInvoice{}
var _ sdk.Msg = &MsgEditInvoice{}
var _ sdk.Msg = &MsgFactorInvoice{}
var _ sdk.Msg = &MsgPayInvoice{}
var _ sdk.Msg = &MsgAddAdmin{}
var _ sdk.Msg = &MsgRemoveAdmin{}
var _ sdk.Msg = &MsgUpdateParams{}

// MsgCreateInvoice defines a message to submit an invoice
type MsgCreateInvoice struct {
	InvoiceID         uint64         `json:"invoiceId"`
	InvoiceNumber     string         `json:"invoiceNumber`
	InvoiceName		  string		 `json:"invoiceName"`
	BillingReason     string		 `json:"billingReason"`
	AmountDue	 	  sdk.Coin	     `json:"amountDue"`
	AmountPaid		  sdk.Coin		 `json:"amountPaid"`
	AmountRemaining   sdk.Coin       `json:"amountRemaining"`
	Subtotal	      sdk.Coin		 `json:"subtotal"`
	Total			  sdk.Coin 	     `json:"total"`
	Party			  sdk.AccAddress `json:"party"`
	Counterparty      sdk.AccAddress `json:"counterparty"`
	DueDate			  time.Time 	 `json:"dueDate"`
	PeriodStartDate   time.Time	     `json:"periodStartDate"`
	PeriodEndDate	  time.Time 	 `json:"periodEndDate"`
	Paid			  bool			 `json:"paid"`
	Active 	          bool           `json:"active"`
	CreatedTime       time.Time      `json:"created_time"`
}

// NewMsgCreateInvoice creates a new message to create a invoice
func NewMsgCreateInvoice(statesetID, invoiceId string, invoiceNumber string, billingReason string, amountDue sdk.Coin, amountpaid sdk.Coin, amountRemaining sdk.Coin, subtotal Int, total Int, dueDate Date, periodStartDate Date, periodEndDate Date, merchant sdk.AccAddress, source string) MsgCreateInvoice {
	return MsgCreateInvoice{
		StatesetID: statesetID,
		Invoice: invoiceID,
		InvoiceNumber: invoiceNumber,
		InvoiceName: invoiceName,
		BillingReason: billingReason,
		AmountDue: amountDue,
		AmountPaid: amountPaid,
		AmountRemaining: amountRemaining,
		Subtotal: subtotal,
		Total: total,
		Party: party,
		Counterparty: counterparty,
		DueDate: dueDate,
		PeriodStartDate: periodStartDate
		PeriodEndDate: periodEndDate
		Paid: paid,
		Active: active,
		CreatedTime: CreatedTime
	}
}

// Route is the name of the route for invoice
func (msg MsgCreateInvoice) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgCreateInvoice) Type() string {
	return TypeMsgCreateInvoice
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgCreateInvoice) ValidateBasic() sdk.Error {
	if len(msg.Body) == 0 {
		return ErrInvalidBodyTooShort(msg.Body)
	}
	if len(msg.StatesetID) == 0 {
		return ErrInvalidStatesetID(msg.StatesetID)
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgCreateInvoice) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgCreateInvoice) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// UpdateInvoice
var _ sdk.Msg = &MsgUpdateInvoice{}

func NewMsgUpdateInvoice(creator string, id string, number string, name string, status string, amountDue string, amountPaid string, amontRemaining string, dueDate string) *MsgUpdateInvoice {
  return &MsgUpdateInvoice{
        Id: id,
		Creator: creator,
    Number: number,
    Name: name,
    Status: status,
    AmountDue: amountDue,
    AmountPaid: amountPaid,
    AmontRemaining: amontRemaining,
    DueDate: dueDate,
	}
}

func (msg *MsgUpdateInvoice) Route() string {
  return RouterKey
}

func (msg *MsgUpdateInvoice) Type() string {
  return "UpdateInvoice"
}

func (msg *MsgUpdateInvoice) GetSigners() []sdk.AccAddress {
  creator, err := sdk.AccAddressFromBech32(msg.Creator)
  if err != nil {
    panic(err)
  }
  return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateInvoice) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateInvoice) ValidateBasic() error {
  _, err := sdk.AccAddressFromBech32(msg.Creator)
  if err != nil {
    return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
  }
   return nil
}

var _ sdk.Msg = &MsgCreateInvoice{}

func NewMsgDeleteInvoice(creator string, id string) *MsgDeleteInvoice {
  return &MsgDeleteInvoice{
        Id: id,
		Creator: creator,
	}
} 
func (msg *MsgDeleteInvoice) Route() string {
  return RouterKey
}

func (msg *MsgDeleteInvoice) Type() string {
  return "DeleteInvoice"
}

func (msg *MsgDeleteInvoice) GetSigners() []sdk.AccAddress {
  creator, err := sdk.AccAddressFromBech32(msg.Creator)
  if err != nil {
    panic(err)
  }
  return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteInvoice) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteInvoice) ValidateBasic() error {
  _, err := sdk.AccAddressFromBech32(msg.Creator)
  if err != nil {
    return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
  }
  return nil
}


// MsgCancelInvoice defines a message to submit a invoice
type MsgCancelInvoice struct {
	ID      uint64         `json:"id"`
	Merchant sdk.AccAddress `json:"merchant"`
}

// Route is the name of the route for invoice
func (msg MsgCancelInvoice) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgCancelInvoice) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgCancelInvoice) ValidateBasic() sdk.Error {
	if msg.ID == 0 {
		return ErrUnknownInvoice(msg.ID)
	}
	if len(msg.Merchant) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Merchant.String())
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgCancelInvoice) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgCancelInvoice) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Merhcant)}
}

// MsgDeleteInvoice defines a message to submit a invoice
type MsgDeleteInvoice struct {
	ID      uint64         `json:"id"`
	Creator sdk.AccAddress `json:"creator"`
}

// Route is the name of the route for invoice
func (msg MsgDeleteInvoice) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgDeleteInvoice) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgDeleteInvoice) ValidateBasic() sdk.Error {
	if msg.ID == 0 {
		return ErrUnknownInvoice(msg.ID)
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgDeleteInvoice) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgDeleteInvoice) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// MsgEditInvoice defines a message to submit a invoice
type MsgEditInvoice struct {
	ID     uint64         `json:"id"`
	Body   string         `json:"body"`
	Editor sdk.AccAddress `json:"editor"`
}

// NewMsgEditInvoice creates a new message to edit a invoice
func NewMsgEditInvoice(id uint64, body string, editor sdk.AccAddress) MsgEditInvoice {
	return MsgEditInvoice{
		ID:     id,
		Body:   body,
		Editor: editor,
	}
}

// Route is the name of the route for invoice
func (msg MsgEditInvoice) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgEditInvoice) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgEditInvoice) ValidateBasic() sdk.Error {
	if msg.ID == 0 {
		return ErrUnknownInvoice(msg.ID)
	}
	if len(msg.Editor) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Editor.String())
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgEditInvoice) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgEditInvoice) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Editor)}
}

// MsgFactorInvoice defines a message to factor and invoice
type MsgFactorInvoice struct {
	ID     uint64         `json:"id"`
	Body   string         `json:"body"`
	Factor sdk.AccAddress `json:"factor"`
}

// NewMsgFactorInvoice creates a new message to factor an invoice
func NewMsgFactorInvoice(id uint64, body string, factor sdk.AccAddress) MsgFactorInvoice {
	return MsgEditInvoice{
		ID:     id,
		Body:   body,
		Factor: factor,
	}
}

// Route is the name of the route for invoice
func (msg MsgfactorInvoice) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgFactorInvoice) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgFactorInvoice) ValidateBasic() sdk.Error {
	if msg.ID == 0 {
		return ErrUnknownInvoice(msg.ID)
	}
	if len(msg.Factor) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Factor.String())
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
	return []sdk.AccAddress{sdk.AccAddress(msg.Factor)}
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