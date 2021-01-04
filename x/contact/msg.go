package contact

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// TypeMsgCreateContact represents the type of the message for creating new contact
	TypeMsgCreateContact = "create_contact"
	// TypeMsgEditContact represents the type of the message for creating new contact
	TypeMsgEditContact = "edit_contact"
	// TypeMsgDeleteContact represents the type of the message for creating new contact
	TypeMsgDeleteContact = "delete_contact"
	// TypeMsgAddAdmin represents the type of message for adding a new admin
	TypeMsgAddAdmin = "add_admin"
	// TypeMsgRemoveAdmin represents the type of message for removeing an admin
	TypeMsgRemoveAdmin = "remove_admin"
	// TypeMsgUpdateParams represents the type of
	TypeMsgUpdateParams = "update_params"
)

// verify interface at compile time
var _ sdk.Msg = &MsgCreateContact{}
var _ sdk.Msg = &MsgEditContact{}
var _ sdk.Msg = &MsgDeleteContact{}
var _ sdk.Msg = &MsgAddAdmin{}
var _ sdk.Msg = &MsgRemoveAdmin{}
var _ sdk.Msg = &MsgUpdateParams{}

// MsgCreateContact defines a message to submit a contact
type MsgCreateContact struct {
	ContactID string         `json:"contact_id"`
	Body        string         `json:"body"`
	Creator     sdk.AccAddress `json:"creator"`
	Source      string         `json:"source,omitempty"`
}

// NewMsgCreateContact creates a new message to create a Contact
func NewMsgCreateContact(contactID, body string, creator sdk.AccAddress, source string) MsgCreateContact {
	return MsgCreateContact{
		ContactID:  ContactID,
		Body:        body,
		Creator:     creator,
		Source:      source,
	}
}

// Route is the name of the route for Contact
func (msg MsgCreateContact) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgCreateContact) Type() string {
	return TypeMsgCreateContact
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgCreateContact) ValidateBasic() sdk.Error {
	if len(msg.Body) == 0 {
		return ErrInvalidBodyTooShort(msg.Body)
	}
	if len(msg.MarketID) == 0 {
		return ErrInvalidMarketID(msg.MarketID)
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgCreateContact) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgCreateContact) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// MsgDeleteContact defines a message to submit a contact
type MsgDeleteContact struct {
	ID      uint64         `json:"id"`
	Creator sdk.AccAddress `json:"creator"`
}

// Route is the name of the route for Contact
func (msg MsgDeleteContact) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgDeleteContact) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgDeleteContact) ValidateBasic() sdk.Error {
	if msg.ID == 0 {
		return ErrUnknownContact(msg.ID)
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgDeleteContact) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgDeleteContact) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// MsgEditContact defines a message to submit a contact
type MsgEditContact struct {
	ContactID     uint64         `json:"contact_id"`
	Body   string         `json:"body"`
	Editor sdk.AccAddress `json:"editor"`
}

// NewMsgEditContact creates a new message to edit a Contact
func NewMsgEditContact(contact_id uint64, body string, editor sdk.AccAddress) MsgEditContact {
	return MsgEditContact{
		ID:     id,
		Body:   body,
		Editor: editor,
	}
}

// Route is the name of the route for Contact
func (msg MsgEditContact) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgEditContact) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgEditContact) ValidateBasic() sdk.Error {
	if msg.ID == 0 {
		return ErrUnknownContact(msg.ID)
	}
	if len(msg.Editor) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Editor.String())
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgEditContact) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgEditContact) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Editor)}
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