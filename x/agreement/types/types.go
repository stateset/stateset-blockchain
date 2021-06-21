package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ID uint64

func (g ID) Uint64() uint64 {
	return uint64(g)
}

func (g ID) Empty() bool {
	return g == 0
}

func (g ID) Bytes() []byte {
	return orm.EncodeSequence(uint64(g))
}

type AgreementID uint64

func (a AgreementID) Bytes() []byte {
	return orm.EncodeSequence(uint64(a))
}

func (a AgreementID) Uint64() uint64 {
	return uint64(a)
}
func (a AgreementID) Empty() bool {
	return a == 0
}

// ID method returns AgreementID details of specific agreement
func (a Agreement) ID() AgreementID {
	return a.AgreementID
}

// String implements the Stringer interface for a Order object.
func (a Agreeement) String() string {
	out, _ := yaml.Marshal(a)
	return string(out)
}

// Agreements is an array of agreements
type Agreements []Agreement

// String implements the Stringer interface for a Orders object.
func (a Agreements) String() string {
	var out string
	for _, order := range a {
		out += order.String() + "\n"
	}

	return strings.TrimSpace(out)
}

// ValidateInactive method validates whether agreement is open or not and
// returns error if not
func (a Agreement) ValidateInactive() error {
	switch a.State {
	case AgreementInactive:
		return nil
	case AgreementActive:
		return ErrAgreementActive
	default:
		return ErrAgreementInactive
	}
}


// Total method returns total agreement value of specific agreement
func (a Agreement) Total() sdk.Coin {
	return a.Spec.Total()
}

// Agreement stores data about an agreement
type Agreement struct {
	agreementID            uint64         `json:"agreementId"`
	agreementNumber        string         `json:"agreementNumber`
	agreementName		   string		  `json:"agreementName"`
	agreementType          string         `json:"agreementType`
	agreementStatus 	   string 		  `json:"agreementStatus"`
	totalAgreementValue    sdk.Coin       `json:"totalAgreementValue"`
	party                  sdk.AccAddress `json:"party"`
	counterparty           sdk.AccAddress `json:"counterparty"`
	AgreementStartBlock	   time.Time 	  `json:"AgreementStartBlock"`
	AgreementEndBlock       time.Time      `json:"AgreementEndBlock`
	paid			  	   bool		      `json:"paid"`
	active 	          	   bool           `json:"active"`
	CreatedTime       	   time.Time      `json:"created_time"`
}

// NewAgreement creates a new agreement object
func NewAgreement(agreementId uint64, agreementNumber string, agreementName string, description string, loanAmount sdk.Coin, amountPaid sdk.Coin, amountRemaining sdk.Coin, subtotal sdk.Coin, total sdk.Coin, party sdk.AccAddress, counterparty sdk.AccAddress, dueDate time.Time, periodStartDate time.Time, periodEndDate time.Time, paid bool, active bool, createdTime time.Time) Loan {
	return Agreement{
		AgreementID:       	 agreementID,
		AgreementNumber:     agreementNumber,
		AgreementName:     	 agreementName,
		AgreementType: 		 agreementType,
		AgreementStatus: 	 agreementStatus,
		TotalAgreementValue: totalAgreementValue,
		Party:		     	 party,
		Counterparty:	 	 counterparty,
		AgreementStartBlock:  AgreementStartBlock,
		AgreementEndBlock:  	 AgreementEndBlock,
		Paid:			  	 paid,
		Active: 		 	 active,
		CreatedTime:     	 createdTime,
	}
}

func (a Agreement) String() string {
	return fmt.Sprintf(`Agreement %d:
		AgreementID:       	 %s,
		AgreementNumber:   	 %s,
		AgreementName:     	 %s,
		AgreementType:  	 %s,
		AgreementStatus:   	 %s,
		TotalAgreementValue: %s,
		Party:				 %s,
		Counterparty:	 	 %s,
		AgreementStartBlock:	 %s,
		AgreementEndBlock:    %s,
		PeriodStartDate:     %s,
		Paid:			     %s,
		Active: 		     %s,
		CreatedTime:         %s`,
		a.AgreementId, a.AgreementNumber, a.AgreementName, a.AgreementType, a.AgreementStatus, a.TotalAgreementValue, a.Party, a.Counterparty, a.AgreementStartBlock.String(), a.AgreementEndBlock.toString(), a.Paid.String(), a.Active.String(), a.CreatedTime.String())
}