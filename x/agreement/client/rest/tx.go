package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/agreement/create", types.ModuleName), createHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/agreement/update", types.ModuleName), updateHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/agreement/delete", types.ModuleName), deleteHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/agreement/activate", types.ModuleName), activateHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/agreement/amend", types.ModuleName), amendHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/agreement/renew", types.ModuleName), renewHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/agreement/terminate", types.ModuleName), terminateHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/agreement/expire", types.ModuleName), expireHandlerFn(cliCtx)).Methods("POST")
}

type createAgreementRequest struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Creator string `json:"creator"`
	AgreementNumber string `json:"agreementNumber"`
	AgreementName string `json:"agreementName"`
	AgreementType string `json:"agreementType"`
	AgreementStatus string `json:"agreementStatus"`
	TotalAgreementValue string `json:"totalAgreementValue"`
	Party string `json:"party"`
	Counterparty string `json:"counterparty"`
	AgreementStartBlock string `json:"AgreementStartBlock"`
	AgreementEndBlock string `json:"AgreementEndBlock"`
	
}

func createAgreementHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createAgreementRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		_, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		parsedAgreementNumber := req.AgreementNumber
		
		parsedAgreementName := req.AgreementName
		
		parsedAgreementType := req.AgreementType
		
		parsedAgreementStatus := req.AgreementStatus
		
		parsedTotalAgreementValue := req.TotalAgreementValue
		
		parsedParty := req.Party
		
		parsedCounterparty := req.Counterparty
		
		parsedAgreementStartBlock := req.AgreementStartBlock
		
		parsedAgreementEndBlock := req.AgreementEndBlock

		msg := types.NewMsgCreateAgreement(
			req.Creator,
			parsedAgreementNumber,
			parsedAgreementName,
			parsedAgreementType,
			parsedAgreementStatus,
			parsedTotalAgreementValue,
			parsedParty,
			parsedCounterparty,
			parsedAgreementStartBlock,
			parsedAgreementEndBlock,
			
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func updateAgreementHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        id := mux.Vars(r)["id"]

		var req updateAgreementRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		_, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		
		parsedAgreementNumber := req.AgreementNumber
		
		parsedAgreementName := req.AgreementName
		
		parsedAgreementType := req.AgreementType
		
		parsedAgreementStatus := req.AgreementStatus
		
		parsedTotalAgreementValue := req.TotalAgreementValue
		
		parsedParty := req.Party
		
		parsedCounterparty := req.Counterparty
		
		parsedAgreementStartBlock := req.AgreementStartBlock
		
		parsedAgreementEndBlock := req.AgreementEndBlock
		

		msg := types.NewMsgUpdateAgreement(
			req.Creator,
            id,
			parsedAgreementNumber,
			parsedAgreementName,
			parsedAgreementType,
			parsedAgreementStatus,
			parsedTotalAgreementValue,
			parsedParty,
			parsedCounterparty,
			parsedAgreementStartBlock,
			parsedAgreementEndBlock,
			
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type deleteAgreementRequest struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Creator string `json:"creator"`
}

func deleteAgreementHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        id := mux.Vars(r)["id"]

		var req deleteAgreementRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		_, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgDeleteAgreement(
			req.Creator,
            id,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func postActivateHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode PUT request body
		var req PostActivateSwapReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// Create and return msg
		msg := types.NewMsgActivateSwap(
			req.From,
			req.AgreementID,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postAmendHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode PUT request body
		var req PostAmendSwapReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// Create and return msg
		msg := types.NewMsgAmendSwap(
			req.From,
			req.AgreementID,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postRenewHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode PUT request body
		var req PostRenewSwapReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// Create and return msg
		msg := types.NewMsgRenewSwap(
			req.From,
			req.AgreementID,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postTerminateHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode PUT request body
		var req PostTerminateSwapReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// Create and return msg
		msg := types.NewMsgTerminateSwap(
			req.From,
			req.AgreementID,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postExpireHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode PUT request body
		var req PostExpireSwapReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// Create and return msg
		msg := types.NewMsgExpireSwap(
			req.From,
			req.AgreementID,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}