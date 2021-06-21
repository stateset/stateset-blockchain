package rest

// DONTCOVER
// client is excluded from test coverage in the poc phase milestone 1 and will be included in milestone 2 with completeness

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/stateset/stateset-blockchain/x/liquidity/types"
)

// TODO: Plans to increase completeness on Milestone 2
// using grpc
func registerTxRoutes(clientCtx client.Context, r *mux.Router) {
	//// create pool
	//r.HandleFunc(fmt.Sprintf("/liquidity/pool"), newPoolHandlerFn(clientCtx)).Methods("POST")
	//// deposit to pool
	//r.HandleFunc(fmt.Sprintf("/liquidity/pool/{%s}/deposit", RestPoolId), newDepositPoolHandlerFn(clientCtx)).Methods("POST")
	//// withdraw from pool
	//r.HandleFunc(fmt.Sprintf("/liquidity/pool/{%s}/withdraw", RestPoolId), newWithdrawPoolHandlerFn(clientCtx)).Methods("POST")
}

func newPoolHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MsgCreatePoolRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		//baseReq := req.BaseReq.Sanitize()
		//if !baseReq.ValidateBasic(w) {
		//	return
		//}
		//
		//poolCreator, e := sdk.AccAddressFromBech32(req.PoolCreator)
		//if e != nil {
		//	rest.WriteErrorResponse(w, http.StatusBadRequest, e.Error())
		//	return
		//}
		//
		//depositCoin, ok := sdk.NewIntFromString(req.DepositCoins)
		//if !ok || !depositCoin.IsPositive() {
		//	rest.WriteErrorResponse(w, http.StatusBadRequest, "coin amount: "+req.DepositCoins)
		//	return
		//}
		//
		//msg := types.NewMsgCreatePool()
		//if err := msg.ValidateBasic(); err != nil {
		//	rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		//	return
		//}
		//
		//tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// HTTP request handler to add liquidity.
func newDepositPoolHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//vars := mux.Vars(r)
		//poolID := vars[RestPoolId]
		//
		//var req DepositPoolReq
		//if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
		//	return
		//}
		//
		//baseReq := req.BaseReq.Sanitize()
		//if !baseReq.ValidateBasic(w) {
		//	return
		//}
		//
		//msg := types.NewMsgDepositToPool()
		//if err := msg.ValidateBasic(); err != nil {
		//	rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		//	return
		//}
		//
		//tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// HTTP request handler to remove liquidity.
func newWithdrawPoolHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//vars := mux.Vars(r)
		//poolID := vars[RestPoolId]
		//
		//var req WithdrawPoolReq
		//if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
		//	return
		//}
		//
		//baseReq := req.BaseReq.Sanitize()
		//if !baseReq.ValidateBasic(w) {
		//	return
		//}
		//
		//withdrawer, err := sdk.AccAddressFromBech32(req.Withdrawer)
		//if err != nil {
		//	return
		//}
		//poolId, err := strconv.ParseUint(req.PoolId, 10, 64)
		//sdk.NewCoin
		//msg := types.NewMsgWithdrawFromPool(withdrawer, poolId, req.PoolCoin)
		//if err := msg.ValidateBasic(); err != nil {
		//	rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		//	return
		//}
		//
		//tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

//// WithdrawPoolReq defines the properties of a Deposit from liquidity Pool request's body
//type CreatePoolReq struct {
//	BaseReq           rest.BaseReq `json:"base_req" yaml:"base_req"`
//	PoolCreator       string       `json:"pool_creator" yaml:"pool_creator"`
//	PoolTypeIndex     string       `json:"pool_type_index" yaml:"pool_type_index"`
//	ReserveCoinDenoms string       `json:"reserve_coin_denoms" yaml:"reserve_coin_denoms"`
//	DepositCoins      string       `json:"deposit_coins" yaml:"deposit_coins"`
//}
//
//// WithdrawPoolReq defines the properties of a Deposit from liquidity Pool request's body
//type WithdrawPoolReq struct {
//	BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"`
//	Withdrawer string       `json:"withdrawer" yaml:"withdrawer"`
//	PoolId     string       `json:"pool_id" yaml:"pool_id"`
//	PoolCoin   string       `json:"pool_coin_amount" yaml:"pool_coin"`
//}
//
//// DepositPoolReq defines the properties of a Deposit liquidity request's body
//type DepositPoolReq struct {
//	BaseReq      rest.BaseReq `json:"base_req" yaml:"base_req"`
//	Depositor    string       `json:"depositor" yaml:"depositor"`
//	PoolId       string       `json:"pool_id" yaml:"pool_id"`
//	DepositCoins string       `json:"deposit_coins_amount" yaml:"deposit_coins"`
//}
//
//// DepositPoolReq defines the properties of a Deposit liquidity request's body
//type SwapReq struct {
//	BaseReq         rest.BaseReq `json:"base_req" yaml:"base_req"`
//	SwapRequester   string       `json:"swap_requester" yaml:"swap_requester"`
//	PoolId          string       `json:"pool_id" yaml:"pool_id"`
//	PoolTypeIndex   string       `json:"pool_type_index" yaml:"pool_type_index"`
//	SwapType        string       `json:"swap_type" yaml:"swap_type"`
//	OfferCoin       string       `json:"offer_coin" yaml:"offer_coin"`
//	DemandCoinDenom string       `json:"demand_coin" yaml:"demand_coin"`
//	OrderPrice      string       `json:"order_price" yaml:"order_price"`
//}
