package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/stateset/stateset-blockchain/x/invoice"
)

const restInvoiceID = "invoice-id"
const restDenom = "denom"

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/invoice/{%s}", types.ModuleName, restSwapID), queryInvoiceHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/invoices", types.ModuleName), queryInvoiceHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/invoiceLineItem/{%s}", types.ModuleName, restDenom), queryInvoiceLineItemHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/invoiceLineItems", types.ModuleName), queryInvoicesLineItemsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/parameters", types.ModuleName), queryParamsHandlerFn(cliCtx)).Methods("GET")

}

func queryInvoiceHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the query height
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// Prepare params for querier
		vars := mux.Vars(r)
		if len(vars[restInvoiceID]) == 0 {
			err := fmt.Errorf("%s required but not specified", restInvoiceID)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		invoiceID, err := hex.DecodeString(vars[restInvoiceID])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		bz, err := cliCtx.Codec.MarshalJSON(types.QueryInvoiceByID{InvoiceID: invoiceID})
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Query
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("/custom/%s/%s", types.ModuleName, types.QueryGetInvoice), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Decode and return results
		cliCtx = cliCtx.WithHeight(height)

		var invoice types.AugmentedInvoice
		err = cliCtx.Codec.UnmarshalJSON(res, &invoice)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, cliCtx.Codec.MustMarshalJSON(invoice))
	}
}

// HTTP request handler to query list of agremeents filtered by optional params
func queryInvoicesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
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
			invoiceStatus    types.InvoiceStatus
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
			invoiceStatus = types.NewInvoiceStatusFromString(x)
			if !invoiceStatus.IsValid() {
				rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid invoice status %s", agreementStatus))
				return
			}
		}

		params := types.NewQueryInvoices(page, limit, involveAddr, expiration, invoiceStatus)
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