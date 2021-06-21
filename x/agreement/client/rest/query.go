package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

)

const restAgreementID = "agreement-id"
const restDenom = "denom"

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/agreement/{%s}", types.ModuleName, restAgreementID), queryAgreementHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/agreements", types.ModuleName), queryAgreementsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/agreementLineItem/{%s}", types.ModuleName, restDenom), queryAgreementLineItemHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/agreementLineItems", types.ModuleName), queryAgreementLineItemsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/parameters", types.ModuleName), queryParamsHandlerFn(cliCtx)).Methods("GET")

}

func queryAgreementHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the query height
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// Prepare params for querier
		vars := mux.Vars(r)
		if len(vars[restAgreementID]) == 0 {
			err := fmt.Errorf("%s required but not specified", restAgreementID)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		agreementID, err := hex.DecodeString(vars[restAgreementID])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		bz, err := cliCtx.Codec.MarshalJSON(types.QueryAgreementByID{AgreementID: agreementID})
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Query
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("/custom/%s/%s", types.ModuleName, types.QueryGetAgreement), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Decode and return results
		cliCtx = cliCtx.WithHeight(height)

		var agreement types.AugmentedAgreement
		err = cliCtx.Codec.UnmarshalJSON(res, &agreement)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, cliCtx.Codec.MustMarshalJSON(agreement))
	}
}

// HTTP request handler to query list of agremeents filtered by optional params
func queryAgreementsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		var (
			involveAddr   sdk.AccAddress
			expiration    uint64
			agreementStatus    types.AgreementStatus
		)

		if x := r.URL.Query().Get(RestInvolve); len(x) != 0 {
			involveAddr, err = sdk.AccAddressFromBech32(x)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		if x := r.URL.Query().Get(RestExpiration); len(x) != 0 {
			expiration, err = strconv.ParseUint(x, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		if x := r.URL.Query().Get(RestStatus); len(x) != 0 {
			agreementStatus = types.NewAgreementStatusFromString(x)
			if !agreementStatus.IsValid() {
				rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid agremeent status %s", agreementStatus))
				return
			}
		}

		params := types.NewQueryAgreements(page, limit, involveAddr, expiration, agreementStatus)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetAtomicSwaps)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetParams)

		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}