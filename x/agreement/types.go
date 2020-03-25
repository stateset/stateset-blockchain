package agreement

import (
	"fmt"
	"net/url"
	"time"

	app "github.com/stateset/stateset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines module constants
const (
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName
	StoreKey          = ModuleName
	DefaultParamspace = ModuleName
)

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
	agreementStartDate	   time.Time 	  `json:"agreementStartDate"`
	agreementEndDate       time.Time      `json:"agreementEndDate`
	paid			  	   bool		      `json:"paid"`
	active 	          	   bool           `json:"active"`
	CreatedTime       	   time.Time      `json:"created_time"`
}

// Agreements is an array of agreements
type Agreements []Agreement

// NewAgreement creates a new agreement object
func NewLoan(agreementId uint64, agreementNumber string, agreementName string, description string, loanAmount sdk.Coin, amountPaid sdk.Coin, amountRemaining sdk.Coin, subtotal sdk.Coin, total sdk.Coin, party sdk.AccAddress, counterparty sdk.AccAddress, dueDate time.Time, periodStartDate time.Time, periodEndDate time.Time, paid bool, active bool, createdTime time.Time) Loan {
	return Agreement{
		AgreementID:       	 agreementID,
		AgreementNumber:     agreementNumber,
		AgreementName:     	 agreementName,
		AgreementType: 		 agreementType,
		AgreementStatus: 	 agreementStatus,
		TotalAgreementValue: totalAgreementValue,
		Party:		     	 party,
		Counterparty:	 	 counterparty,
		AgreementStartDate:  agreementStartDate,
		AgreementEndDate:  	 agreementEndDate,
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
		AgreementStartDate:	 %s,
		AgreementEndDate:    %s,
		PeriodStartDate:     %s,
		Paid:			     %s,
		Active: 		     %s,
		CreatedTime:         %s`,
		a.AgreementId, a.AgreementNumber, a.AgreementName, a.AgreementType, a.AgreementStatus, a.TotalAgreementValue, a.Party, a.Counterparty, a.AgreementStartDate.String(), a.AgreementEndDate.toString(), a.Paid.String(), a.Active.String(), a.CreatedTime.String())
}