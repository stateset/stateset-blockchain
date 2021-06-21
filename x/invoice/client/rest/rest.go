package rest

import (
	"github.com/gorilla/mux"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// REST Variable names
// nolint
const (
	RestExpiration = "expiration"
	RestInvolve    = "involve"
	RestStatus     = "status"
	RestType  = "type"
)

// RegisterRoutes registers bep3-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

// PostCreateInvoiceReq defines the properties of a invoice create request's body
type PostCreateInvoiceReq struct {
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

// PostInvoiceReq defines the properties of a  request's body
type PostPayInvoiceReq struct {
	BaseReq      rest.BaseReq     `json:"base_req" yaml:"base_req"`
	From         sdk.AccAddress   `json:"from" yaml:"from"`
	InvoiceID       tmbytes.HexBytes `json:"invoice_id" yaml:"invoice_id"`
}