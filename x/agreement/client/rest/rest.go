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
	AgreementStartBlock	   time.Time 	  `json:"AgreementStartBlock"`
	AgreementEndBlock       time.Time      `json:"AgreementEndBlock`
	paid			  	   bool		      `json:"paid"`
	active 	          	   bool           `json:"active"`
	CreatedTime       	   time.Time      `json:"created_time"`
}

type PostActivateAgreementReq struct {
	BaseReq      rest.BaseReq     `json:"base_req" yaml:"base_req"`
	From         sdk.AccAddress   `json:"from" yaml:"from"`
	AgreementID       tmbytes.HexBytes `json:"agreement_id" yaml:"agreement_id"`
}

type PostAmendAgreementReq struct {
	BaseReq      rest.BaseReq     `json:"base_req" yaml:"base_req"`
	From         sdk.AccAddress   `json:"from" yaml:"from"`
	AgreementID       tmbytes.HexBytes `json:"agreement_id" yaml:"agreement_id"`
}

type PostRenewAgreementReq struct {
	BaseReq      rest.BaseReq     `json:"base_req" yaml:"base_req"`
	From         sdk.AccAddress   `json:"from" yaml:"from"`
	AgreementID       tmbytes.HexBytes `json:"agreement_id" yaml:"agreement_id"`
}

type PostTerminateAgreementReq struct {
	BaseReq      rest.BaseReq     `json:"base_req" yaml:"base_req"`
	From         sdk.AccAddress   `json:"from" yaml:"from"`
	AgreementID       tmbytes.HexBytes `json:"agreement_id" yaml:"agreement_id"`
}

type PostExpireAgreementReq struct {
	BaseReq      rest.BaseReq     `json:"base_req" yaml:"base_req"`
	From         sdk.AccAddress   `json:"from" yaml:"from"`
	AgreementID       tmbytes.HexBytes `json:"agreement_id" yaml:"agreement_id"`
}