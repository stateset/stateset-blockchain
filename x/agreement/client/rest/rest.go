package rest

import (
	"github.com/gorilla/mux"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/cosmos-sdk/client/context"
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

// PostCreateAgreementReq defines the properties of a agreement create request's body
type PostCreateAgreementReq struct {
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

// PostAgreementSwapReq defines the properties of a activate agreement request's body
type PostActivateAgreemenyReq struct {
	BaseReq      rest.BaseReq     `json:"base_req" yaml:"base_req"`
	From         sdk.AccAddress   `json:"from" yaml:"from"`
	AgreementID       tmbytes.HexBytes `json:"agreement_id" yaml:"agreement_id"`
}