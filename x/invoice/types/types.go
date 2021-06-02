package invoice

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

// Invoice stores data about an invoice
type Invoice struct {
	invoiceID         uint64         `json:"invoiceId"`
	invoiceNumber     string         `json:"invoiceNumber`
	invoiceName		  string		 `json:"invoiceName"`
	billingReason     string		 `json:"billingReason"`
	amountDue	 	  sdk.Coin	     `json:"amountDue"`
	amountPaid		  sdk.Coin		 `json:"amountPaid"`
	amountRemaining   sdk.Coin       `json:"amountRemaining"`
	subtotal	      sdk.Coin		 `json:"subtotal"`
	total			  sdk.Coin 	     `json:"total"`
	party			  sdk.AccAddress `json:"party"`
	counterparty      sdk.AccAddress `json:"counterparty"`
	dueDate			  time.Time 	 `json:"dueDate"`
	periodStartDate   time.Time	     `json:"periodStartDate"`
	periodEndDate	  time.Time 	 `json:"periodEndDate"`
	paid			  bool			 `json:"paid"`
	active 	          bool           `json:"active"`
	CreatedTime       time.Time      `json:"created_time"`
}

// Invoices is an array of invoices
type Invoices []Invoice

// NewInvoice creates a new invoice object
func NewInvoice(invoiceId uint64, invoiceNumber string, invoiceName string, billingReason string, amountDue sdk.Coin, amountPaid sdk.Coin, amountRemaining sdk.Coin, subtotal sdk.Coin, total sdk.Coin, party sdk.AccAddress, counterparty sdk.AccAddress, dueDate time.Time, periodStartDate time.Time, periodEndDate time.Time, paid bool, active bool, createdTime time.Time) Invoice {
	return Invoice{
		InvoiceID:       invoiceId,
		InvoiceNumber:   invoiceNumber,
		InvoiceName:     invoiceName,
		BillingReason:   billingReason,
		AmountDue:       amountDue,
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

func (i Invoice) String() string {
	return fmt.Sprintf(`Invoice %d:
		InvoiceID:       %s,
		InvoiceNumber:   %s,
		InvoiceName:     %s,
		BillingReason:   %s,
		AmountDue:       %s
		AmountPaid:		 %s,
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
		i.InvoiceID, i.InvoiceNumber, i.InvoiceName, i.BillingReason, i.AmountDue.String(), i.AmountPaid.String(), i.AmountRemaining.String(), i.Subtotal.String(), i.Total.String(), i.Party.String(), i.Counterparty.String(), i.DueDate.String(), i.PeriodStartDate.String(), i.PeriodEndDate.String(), i.Paid.String(), i.Active.String(), a.CreatedTime.String())
}