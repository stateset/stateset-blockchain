package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// TypeMsgCreateAgreement represents the type of the message for creating new agreement
	TypeMsgCreateAgreement = "create_agreement"
	// TypeMsgEditAgreement represents the type of the message for updating an agreement
	TypeMsgUpdateAgreement = "update_agreement"
	// TypeMsgActivateAgreement represents the type of the message for delete an agreement
	TypeMsgDeleteAgreement = "delete_agreement"
	// TypeMsgActivateAgreement represents the type of the message for activating an agreement
	TypeMsgActivateAgreement = "activate_agreement"
	// TypeMsgRenewAgreement represents the type of the message for amending an agreement
	TypeMsgAmendAgreement = "amend_agreement"
	// TypeMsgRenewAgreement represents the type of the message for renewing an agreement
	TypeMsgRenewAgreement = "renew_agreement"
	// TypeMsgTerminateAgreement represents the type of the message for renewing an agreement
	TypeMsgTerminateAgreement = "terminate_agreement"
	// TypeMsgPaybackLoan represents the type of the message for expiring an agreement
	TypeMsgExpireAgreement = "expire_agreement"
	// TypeMsgPaybackLoan represents the type of the message for expiring an agreement
	TypeMsgSendIbcAgreement = "send_ibc_agreement"
)

// verify interface at compile time
var _ sdk.Msg = &MsgCreateAgreement{}
var _ sdk.Msg = &MsgUpdateAgreement{}
var _ sdk.Msg = &MsgDeleteAgreement{}
var _ sdk.Msg = &MsgActivateAgreement{}
var _ sdk.Msg = &MsgAmendAgreement{}
var _ sdk.Msg = &MsgRenewAgreement{}
var _ sdk.Msg = &MsgTerminateAgreement{}
var _ sdk.Msg = &MsgExpireAgreement{}
var _ sdk.Msg = &MsgSendIbcAgreement{}

// MsgCreateAgreement defines a message to create an agreement
type MsgCreateAgreement struct {
	AgreementID         uint64 `json:"agreementid"`
	AgreementName       string `json:"agreementName"`
	AgreementNumber     string `json:"agreementNumber"`
	AgreementType       string `json:"agreementType"`
	AgreementStatus     string `json:"agreementStatus"`
	AgreementNumber     string `json:"agreementNumber"`
	Party               string `json:"party"`
	Counterparty        string `json:"counterparty"`
	AgreementStartBlock string `json:"AgreementStartBlock"`
	AgreementEndBlock   string `json:"AgreementEndBlock"`
}

// NewMsgCreateAgreement creates a new message to create an agreement
func NewMsgCreateAgreement(agreementID string, agreementNumber string, agreementName string, agreementType string, agreementStatus string, totalAgreementValue int, party sdk.AccAddress, counterparty sdk.AccAddress, agreementStartBlock string, agreementEndBlock string) MsgCreateAgreement {
	return MsgCreateAgreement{
		AgreementID:         agreementID,
		AgreementNumber:     agreementNumber,
		AgreementName:       agreementName,
		AgreementType:       agreementType,
		AgreementStatus:     agreementStatus,
		TotalAgreementValue: totalAgreementValue,
		Party:               party,
		Counterparty:        counterparty,
		AgreementStartBlock: agreementStartBlock,
		AgreementEndBlock:   agreementEndBlock,
	}
}

// Route is the name of the route for agreement
func (msg MsgCreateAgreement) Route() string { return RouterKey }

// Type is the name for the Msg
func (msg MsgCreateAgreement) Type() string { return TypeMsgCreateAgreement }

// ValidateBasic validates basic fields of the Msg
func (msg MsgCreateAgreement) ValidateBasic() sdk.Error {
	if len(msg.TotalAgreementValue) == 0 {
		return ErrInvalidAgreementTooSmall(msg.TotalAgreementValue)
	}
	if len(msg.Counterparty) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Counterparty.String())
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgCreateAgreement) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// Update Agreement
func NewMsgUpdateAgreement(creator string, agreementId string, agreementNumber string, agreementName string, agreementType string, agreementStatus string, totalAgreementValue string, party string, counterparty string, AgreementStartBlock string, AgreementEndBlock string) *MsgUpdateAgreement {
	return &MsgUpdateAgreement{
		AgreeemntId:         agreementId,
		AgreementNumber:     agreementNumber,
		AgreementName:       agreementName,
		AgreementType:       agreementType,
		AgreementStatus:     agreementStatus,
		TotalAgreementValue: totalAgreementValue,
		Party:               party,
		Counterparty:        counterparty,
		AgreementStartBlock: AgreementStartBlock,
		AgreementEndBlock:   AgreementEndBlock,
	}
}

func (msg *MsgUpdateAgreement) Route() string { return RouterKey }

func (msg *MsgUpdateAgreement) Type() string { return "UpdateAgreement" }

func (msg *MsgUpdateAgreement) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateAgreement) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateAgreement) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateAgreement{}

// Delete Agreement
func NewMsgDeleteAgreement(creator string, id string) *MsgDeleteAgreement {
	return &MsgDeleteAgreement{
		Id:      id,
		Creator: creator,
	}
}
func (msg *MsgDeleteAgreement) Route() string {
	return RouterKey
}

func (msg *MsgDeleteAgreement) Type() string {
	return "DeleteAgreement"
}

func (msg *MsgDeleteAgreement) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteAgreement) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteAgreement) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// Amend Agreement

// MsgAmendAgreement defines a message to amend an agreement
type MsgAmendAgreement struct {
	ID           uint64         `json:"id"`
	Counterparty sdk.AccAddress `json:"counterparty"`
}

// Route is the name of the route for loan
func (msg MsgAmendAgreement) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgAmendAgreement) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgAmendAgreement) ValidateBasic() sdk.Error {
	if msg.ID == 0 {
		return ErrUnknownAgreement(msg.ID)
	}
	if len(msg.Counterparty) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Counterparty.String())
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgAmendAgreement) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgAmendAgreement) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Party)}
}

// Activate Agreement

// MsgActivateAgreement defines a message to activate an agreement
type MsgActivateAgreement struct {
	ID              uint64         `json:"id"`
	AgreementStatus string         `json:"agreementStatus"`
	Counterparty    sdk.AccAddress `json:"counterparty"`
}

// Route is the name of the route for an agreement
func (msg MsgActivateAgreement) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgActivateAgreement) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgActivateAgreement) ValidateBasic() sdk.Error {
	if msg.AgreementID == 0 {
		return ErrUnknownAgreement(msg.AgreementID)
	}
	if len(msg.Counterparty) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Counterparty.String())
	}

	if msg.AgreementStatus != "ACTIVATED" {
		return Error("The Agreement status must be Activated.")
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgActivateAgreement) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgActivateAgreement) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Party)}
}

// Renew Agreement

// MsgRenewAgreement defines a message to renew an agreement
type MsgRenewAgreement struct {
	AgreementID     uint64         `json:"agreementid"`
	AgreementStatus string         `json:"agreementStatus"`
	Counterparty    sdk.AccAddress `json:"counterparty"`
}

// Route is the name of the route for an agreement
func (msg MsgRenewAgreement) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgRenewAgreement) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgRenewAgreement) ValidateBasic() sdk.Error {
	if msg.AgreementID == 0 {
		return ErrUnknownAgreement(msg.AgreementID)
	}
	if len(msg.Counterparty) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Counterparty.String())
	}

	if msg.AgreementStatus != "RENEWED" {
		return Error("The Agreement status must be Renewed.")
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgRenewAgreement) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgRenewAgreement) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Party)}
}

// Terminate Agreement

// MsgTerminateAgreement defines a message to terminate an agreement
type MsgTerminateAgreement struct {
	ID              uint64         `json:"id"`
	AgreementStatus string         `json:"agreementStatus"`
	Counterparty    sdk.AccAddress `json:"counterparty"`
}

// Route is the name of the route for an agreement
func (msg MsgTerminateAgreement) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgTerminateAgreement) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgTerminateAgreement) ValidateBasic() sdk.Error {
	if msg.AgreementID == 0 {
		return ErrUnknownAgreement(msg.AgreementID)
	}
	if len(msg.Counterparty) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Counterparty.String())
	}

	if msg.AgreementStatus != "TERMINATED" {
		return Error("The Agreement status must be Terminated.")
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgTerminateAgreement) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgTerminateAgreement) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Party)}
}

// Expire Agreement

// MsgExpireAgreement defines a message to expire an agreement
type MsgExpireAgreement struct {
	ID              uint64         `json:"id"`
	AgreementStatus string         `json:"agreementStatus"`
	Counterparty    sdk.AccAddress `json:"counterparty"`
}

// Route is the name of the route for an agreement
func (msg MsgExpireAgreement) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgExpireAgreement) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgExpireAgreement) ValidateBasic() sdk.Error {
	if msg.AgreementID == 0 {
		return ErrUnknownAgreement(msg.AgreementID)
	}
	if len(msg.Counterparty) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Counterparty.String())
	}

	if msg.AgreementStatus != "EXPIRED" {
		return Error("The Agreement status must be Expired.")
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgExpireAgreement) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgExpireAgreement) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Party)}
}


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

