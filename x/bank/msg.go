package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgSend     = "send"
	TypeMsgUpdateParams = "update_params"
)

var (
	_ sdk.Msg = &MsgSend{}
)

type MsgSend struct {
	Sender    sdk.AccAddress
	Recipient sdk.AccAddress
	Amount    sdk.Coin
}

func NewMsgSend(sender, recipient sdk.AccAddress, amount sdk.Coin) MsgSend {
	return MsgSend{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
}
func (msg MsgSend) Route() string { return RouterKey }

func (msg MsgSend) Type() string {
	return TypeMsgSend
}

func (msg MsgSend) ValidateBasic() sdk.Error {
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

func (msg MsgSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// Implements Msg.
func (msg MsgSend) GetSignBytes() []byte {
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


//----------------------------------------
// MsgIssue

// MsgIssue - high level transaction of the coin module
type MsgIssue struct {
	Banker  sdk.AccAddress `json:"banker"`
	Outputs []Output       `json:"outputs"`
}

var _ sdk.Msg = MsgIssue{}

//----------------------------------------
// MsgBurn

// MsgBurn - high level transaction of the coin module
type MsgBurn struct {
	Owner sdk.AccAddress `json:"owner"`
	Coins sdk.Coins      `json:"coins"`
}

var _ sdk.Msg = MsgBurn{}






// NewMsgIssue - construct arbitrary multi-in, multi-out send msg.
func NewMsgBurn(owner sdk.AccAddress, coins sdk.Coins) MsgBurn {
	return MsgBurn{Owner: owner, Coins: coins}
}

func (msg MsgBurn) Route() string { return "bank" }
func (msg MsgBurn) Type() string  { return "burn" }

// Implements Msg.
func (msg MsgBurn) ValidateBasic() sdk.Error {
	if len(msg.Owner) == 0 {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Coins) == 0 {
		return ErrBurnEmptyCoins(DefaultCodespace).TraceSDK("")
	}
	if !msg.Coins.IsValidV0() {
		return sdk.ErrInvalidCoins(msg.Coins.String())
	}
	return nil
}

// Implements Msg.
func (msg MsgBurn) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgBurn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}