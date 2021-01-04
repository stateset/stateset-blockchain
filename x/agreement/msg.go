package agreement

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// TypeMsgCreateAgreement represents the type of the message for creating new agreement
	TypeMsgCreateAgreement = "create_agreement"
	// TypeMsgEditAgreement represents the type of the message for editing an agreement
	TypeMsgEditAgreement = "edit_agreement"
	// TypeMsgActivateAgreement represents the type of the message for activating an agreement
	TypeMsgActivateAgreement = "activate_agreement"
	// TypeMsgRenewAgreement represents the type of the message for amending an agreement
	TypeMsgAmendAgreement = "amend_agreement"
	// TypeMsgRenewAgreement represents the type of the message for amending an agreement
	TypeMsgRenewAgreement = "renew_agreement"
	// TypeMsgTerminateAgreement represents the type of the message for renewing an agreement
	TypeMsgTerminateAgreement = "terminate_agreement"
	// TypeMsgPaybackLoan represents the type of the message for creating new loan
	TypeMsgExpireAgreement = "expire_agreement"
	// TypeMsgAddAdmin represents the type of message for adding a new admin
	TypeMsgAddAdmin = "add_admin"
	// TypeMsgRemoveAdmin represents the type of message for removing an admin
	TypeMsgRemoveAdmin = "remove_admin"
	// TypeMsgUpdateParams represents the type of
	TypeMsgUpdateParams = "update_params"
)

// verify interface at compile time
var _ sdk.Msg = &MsgCreateAgreement{}
var _ sdk.Msg = &MsgEditAgreement{}
var _ sdk.Msg = &MsgActivateAgreement{}
var _ sdk.Msg = &MsgAmendAgreement{}
var _ sdk.Msg = &MsgRenewAgreement{}
var _ sdk.Msg = &MsgTerminateAgreement{}
var _ sdk.Msg = &MsgExpireAgreement{}
var _ sdk.Msg = &MsgRemoveAdmin{}
var _ sdk.Msg = &MsgUpdateParams{}

// MsgCreateAgreement defines a message to create an agreement
type MsgCreateAgreement struct {
	MarketID string             `json:"market_id"`
	AgreementID 	  string 			 `json:"agreement_id"`
	Body          string         	 `json:"body"`
	Lender        sdk.AccAddress     `json:"counterparty"`
	Source        string             `json:"source,omitempty"`
}

// NewMsgCreateAgreement creates a new message to create an agreement
func NewMsgCreateAgreement(marketID, agreementID, body string, lender sdk.AccAddress, source string) MsgCreateAgreement {
	return MsgCreateAgreement {
		MarketID: MarketID,
		AgreementID:    agreementID,
		Body:        body,
		Lender:     lender,
		Source:      source,
	}
}

// Route is the name of the route for agreement
func (msg MsgCreateAgreement) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgCreateAgreement) Type() string {
	return TypeMsgCreateAgreement
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgCreateAgreement) ValidateBasic() sdk.Error {
	if len(msg.Body) == 0 {
		return ErrInvalidBodyTooShort(msg.Body)
	}
	if len(msg.MarketID) == 0 {
		return ErrInvalidMarketID(msg.MarketID)
	}
	if len(msg.InvoiceID) == 0 {
		return ErrINvalidInvoiceID(msg.AgreementID)
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

// GetSigners gets the signs of the Msg
func (msg MsgCreateAgreement) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Counterparty}
}


// Amend Agreement

// MsgAmendAgreement defines a message to amend an agreement
type MsgAmendAgreement struct {
	ID      uint64         `json:"id"`
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
	AgreementID              uint64         `json:"agreementid"`
	AgreementStatus 		 string         `json:"agreementStatus"`
	Counterparty    		 sdk.AccAddress `json:"counterparty"`
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








// MsgEditAgreement defines a message to edit an agreement
type MsgEditAgreement struct {
	AgreementID     uint64         `json:"agreementid"`
	AgreementName string `json:"agreementName"`
	AgreementNumber string `json:"agreementNumber"`
	TotalAgreementValue string `json:"totalAgreementValue"`
	AgreementStartDate string `json:"agreementStartDate"`
	AgreementEndDate string `json:"agreementEndDate"`
	Editor sdk.AccAddress `json:"editor"`
}

// NewMsgEditAgreement creates a new message to edit a loan
func NewMsgEditAgreement(agreementId uint64, agreementName string, agreementNumber, totalAgreementValue sdk.Coin, agreementStartDate time.Time, agreementEndDate time.Time, editor sdk.AccAddress) MsgEditLoan {
	return MsgEditLoan{
		ID:     agreementid,
		AgreementName: agreementName,
		AgreementNumber: agreementNumber,
		TotalAgreementValue totalAgreementValue,
		AgreementStartDate: agreementStartDate,
		AgreementEndDate: agreementEndDate,
		Editor: editor,
	}
}

// Route is the name of the route for loan
func (msg MsgEditAgreement) Route() string {
	return RouterKey
}

// Type is the name for the Msg
func (msg MsgEditAgreement) Type() string {
	return ModuleName
}

// ValidateBasic validates basic fields of the Msg
func (msg MsgEditAgreement) ValidateBasic() sdk.Error {
	if msg.ID == 0 {
		return ErrUnknownLoan(msg.ID)
	}
	if len(msg.Editor) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Editor.String())
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgEditAgreement) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgEditAgreement) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Editor)}
}