package market

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// TypeMsgNewMarketplace represents the type of the message for creating new marketplace
	TypeMsgNewMarket = "new_market"
	// TypeMsgAddAdmin represents the type of message for adding a new admin
	TypeMsgAddAdmin = "add_admin"
	// TypeMsgRemoveAdmin represents the type of message for removeing an admin
	TypeMsgRemoveAdmin = "remove_admin"
	// TypeMsgUpdateParams represents the type of
	TypeMsgUpdateParams = "update_params"
)

// MsgNewMarket defines the message to add a new admin
type MsgNewMarket struct {
	Name        string         `json:"name"`
	ID          string         `json:"id"`
	Description string         `json:"description"`
	Merchant     sdk.AccAddress `json:"merchant"`
}

// NewMsgNewMarketplace returns the messages to create a new marketplace
func NewMsgNewMarket(id, name, description string, merchant sdk.AccAddress) MsgNewMarket {
	return MsgNewMarket{
		Name:        name,
		ID:          id,
		Description: description,
		Merchant:     merchant,
	}
}

// ValidateBasic implements Msg
func (msg MsgNewMarket) ValidateBasic() sdk.Error {
	if len(msg.Merchant) == 0 {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid address: %s", msg.Merchant.String()))
	}

	return nil
}

// Route implements Msg
func (msg MsgNewMarket) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgNewMarket) Type() string { return TypeMsgNewMarket }

// GetSignBytes implements Msg
func (msg MsgNewMarket) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners implements Msg. Returns the merchant as the signer.
func (msg MsgNewMarket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Merchant)}
}

// MsgAddAdmin defines the message to add a new admin
type MsgAddAdmin struct {
	Admin   sdk.AccAddress `json:"admin"`
	Merchant sdk.AccAddress `json:"merchant"`
}

// NewMsgAddAdmin returns the messages to add a new admin
func NewMsgAddAdmin(admin, merchant sdk.AccAddress) MsgAddAdmin {
	return MsgAddAdmin{
		Admin:   admin,
		Merchant: merchant,
	}
}

// ValidateBasic implements Msg
func (msg MsgAddAdmin) ValidateBasic() sdk.Error {
	if len(msg.Admin) == 0 {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid address: %s", msg.Admin.String()))
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