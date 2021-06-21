package rest

import (
	"github.com/gorilla/mux"

	
)

const (
	RestCodeID          = "code_id"
	RestContractAddress = "contract_address"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}