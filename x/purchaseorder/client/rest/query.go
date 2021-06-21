package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/stateset/stateset-blockchain/x/purchaseorder"
)

const restPurchaseOrderID = "purchasorder-id"
const restDenom = "denom"

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/purchaseorder/{%s}", types.ModuleName, restSwapID), queryPurchaseOrderHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/purchaseorders", types.ModuleName), queryPurchaseOrdersHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/purchaseorderLineItem/{%s}", types.ModuleName, restDenom), queryPurchaseOrderLineItemHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/purchaseorderLineItems", types.ModuleName), queryPurchaseOrderLineItemsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/parameters", types.ModuleName), queryParamsHandlerFn(cliCtx)).Methods("GET")

}

func queryPurchaseOrderHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the query height
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// Prepare params for querier
		vars := mux.Vars(r)
		if len(vars[restPurchaseOrderID]) == 0 {
			err := fmt.Errorf("%s required but not specified", restPurchaseOrderID)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		purchaseorderID, err := hex.DecodeString(vars[restPurchaseOrderID])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		bz, err := cliCtx.Codec.MarshalJSON(types.QueryPurchaseOrderByID{PurchaseOrderID: purchaseorderID})
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Query
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("/custom/%s/%s", types.ModuleName, types.QueryGetPurchaseOrder), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Decode and return results
		cliCtx = cliCtx.WithHeight(height)

		var purchaseorder types.AugmentedPurchaseOrder
		err = cliCtx.Codec.UnmarshalJSON(res, &purchaseorder)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, cliCtx.Codec.MustMarshalJSON(purchaseorder))
	}
}

// HTTP request handler to query list of agremeents filtered by optional params
func queryPurchaseOrdersHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
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
			purchaseorderStatus    types.PurchaseOrderStatus
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
			purchaseorderStatus = types.NewPurchaseOrderStatusFromString(x)
			if !purchaseorderStatus.IsValid() {
				rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid purchase order status %s", purchaseorderStatus))
				return
			}
		}

		params := types.NewQueryPurchaseOrders(page, limit, involveAddr, expiration, purchaseorderStatus)
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