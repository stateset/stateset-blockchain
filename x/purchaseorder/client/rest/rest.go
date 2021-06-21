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

type PostCreatePurchaseOrderReq struct {
	purchaseOrderID         uint64         `json:"purchaseOrderId"`
	purchaseOrderNumber     string         `json:"purchaseOrderNumber`
	purchaseOrderName		  string		 `json:"purchaseOrderName"`
	description     string		 `json:"description"`
	subtotal	      sdk.Coin		 `json:"subtotal"`
	total			  sdk.Coin 	     `json:"total"`
	purchaser			  sdk.AccAddress `json:"purchaser"`
	vendor      sdk.AccAddress `json:"vendor"`
	financer      sdk.AccAddress `json:"financer"`
	purchaseDate			  time.Time 	 `json:"purchaseDate"`
	deliveryDate			  time.Time 	 `json:"deliveryDate"`
	paid			  bool			 `json:"paid"`
	active 	          bool           `json:"active"`
	CreatedTime       time.Time      `json:"created_time"`
}


type PostCompletePurchaseOrderReq struct {
	BaseReq      rest.BaseReq     `json:"base_req" yaml:"base_req"`
	From         sdk.AccAddress   `json:"from" yaml:"from"`
	PurchaseOrderID       tmbytes.HexBytes `json:"purchaseorder_id" yaml:"purchaseorder_id"`
}

type PostCancelPurchaseOrderReq struct {
	BaseReq      rest.BaseReq     `json:"base_req" yaml:"base_req"`
	From         sdk.AccAddress   `json:"from" yaml:"from"`
	PurchaseOrderID       tmbytes.HexBytes `json:"purchaseorder_id" yaml:"purchaseorder_id"`
}

type PostFinancePurchaseOrderReq struct {
	BaseReq      rest.BaseReq     `json:"base_req" yaml:"base_req"`
	From         sdk.AccAddress   `json:"from" yaml:"from"`
	PurchaseOrderID       tmbytes.HexBytes `json:"purchaseorder_id" yaml:"purchaseorder_id"`
}