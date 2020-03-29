package loan

import (
	"fmt"
	"net/url"
	"time"

	app "github.com/stateset/stateset-blockchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines module constants
const (
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName
	StoreKey          = ModuleName
	DefaultParamspace = ModuleName
)

// Loan stores data about an loan
type Loan struct {
	loanID            uint64         `json:"loanId"`
	loanNumber        string         `json:"loanNumber`
	loanName		  string		 `json:"loanName"`
	description    	  string		 `json:"description"`
	loanAmount        sdk.Coin       `json:"loanAmount"`
	paidAmount 		  sdk.Coin		 `json:"paidAmount"`
	amountRemaining   sdk.Coin       `json:"amountRemaining"`
	subtotal	      sdk.Coin		 `json:"subtotal"`
	total			  sdk.Coin 	     `json:"total"`
	invoice           app.Invoice    `json:"invoice"`
	party			  sdk.AccAddress `json:"party"`
	counterparty      sdk.AccAddress `json:"counterparty"`
	dueDate			  time.Time 	 `json:"dueDate"`
	paid			  bool			 `json:"paid"`
	active 	          bool           `json:"active"`
	CreatedTime       time.Time      `json:"created_time"`
}

// Loans is an array of loans
type Loans []Loan

// NewLoan creates a new loan object
func NewLoan(loanId uint64, loanNumber string, loanName string, description string, loanAmount sdk.Coin, amountPaid sdk.Coin, amountRemaining sdk.Coin, subtotal sdk.Coin, total sdk.Coin, party sdk.AccAddress, counterparty sdk.AccAddress, dueDate time.Time, periodStartDate time.Time, periodEndDate time.Time, paid bool, active bool, createdTime time.Time) Loan {
	return Loan{
		LoanID:       	 loanId,
		LoanNumber:      loanNumber,
		LoanName:     	 loanName,
		Description:     description,
		LoanAmount:      loanAmount,
		AmountPaid:		 amountPaid,
		AmountRemaining  amountRemaining,
		Subtotal:        subtotal,
		Total: 			 total,
		Party:		     party,
		Counterparty:	 counterparty,
		DueDate:		 dueDate,
		PeriodStartDate: periodStartDate,
		PeriodEndDate:   periodEndDate,
		Paid:			 paid,
		Active: 		 active,
		CreatedTime:     createdTime,
	}
}

func (l Loan) String() string {
	return fmt.Sprintf(`Loan %d:
		LoanID:       	 %s,
		LoanNumber:   	 %s,
		LoanName:     	 %s,
		Description:  	 %s,
		LoanAmount:   	 %s,
		AmountPaid:	  	 %s,
		AmountRemaining  %s,
		Subtotal:        %s,
		Total: 			 %s,
		Party:		     %s,
		Counterparty:	 %s,
		DueDate:		 %s,
		PeriodStartDate: %s,
		PeriodEndDate:   %s,
		Paid:			 %s,
		Active: 		 %s,
		CreatedTime:     %s`,
		l.LoanId, l.LoanNumber, l.LoanName, l.Description, l.LoanAmount, l.AmountPaid, l.AmountRemaining, l.Subtotal, l.Total, l.Counterparty.String(), l.DueDate.String(), l.PeriodStartDate.String(), l.PeriodEndDate.String(), l.Paid.String(), l.Active.String(), l.CreatedTime.String())
}