package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgSendState     = "send_state"
	TypeMsgUpdateParams = "update_params"
)

var (
	_ sdk.Msg = &MsgSendState{}
)

type MsgSendState struct {
	Sender    sdk.AccAddress
	Recipient sdk.AccAddress
	Amount    sdk.Coin
}

func NewMsgSendState(sender, recipient sdk.AccAddress, amount sdk.Coin) MsgSendState {
	return MsgSendGift{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
}
func (msg MsgSendState) Route() string { return RouterKey }

func (msg MsgSendState) Type() string {
	return TypeMsgSendState
}

func (msg MsgSendState) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress("invalid creator address")
	}
	if len(msg.Recipient) == 0 {
		return sdk.ErrInvalidAddress("invalid recipient address")
	}

	if msg.Amount.IsNegative() || msg.Reward.IsZero() {
		return sdk.ErrInvalidCoins("invalid coins")
	}
	return nil
}

func (msg MsgSendState) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgSendState) GetSignBytes() []byte {
	bz := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
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