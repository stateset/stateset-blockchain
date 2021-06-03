package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSendIbcAgreement{}

func NewMsgSendIbcAgreement(
	sender string,
	port string,
	channelID string,
	timeoutTimestamp uint64,
	name string,
	number string,
	status string,
	totalagreementvalue string,
) *MsgSendIbcAgreement {
	return &MsgSendIbcAgreement{
		Sender:              sender,
		Port:                port,
		ChannelID:           channelID,
		TimeoutTimestamp:    timeoutTimestamp,
		Name:                name,
		Number:              number,
		Status:              status,
		Totalagreementvalue: totalagreementvalue,
	}
}

func (msg *MsgSendIbcAgreement) Route() string {
	return RouterKey
}

func (msg *MsgSendIbcAgreement) Type() string {
	return "SendIbcAgreement"
}

func (msg *MsgSendIbcAgreement) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgSendIbcAgreement) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSendIbcAgreement) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
