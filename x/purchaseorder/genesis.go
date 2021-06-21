package purchaseorder

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	PurchaseOrders []PurchaseOrder `json:"purchaseOrders"`
	Params         Params          `json:"params"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{
		PurchaseOrders: nil,
		Params:         DefaultParams(),
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState { return NewGenesisState() }

// InitGenesis initializes stateset from genesis file
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	for _, c := range data.PurchaseOrders {
		k.setPurchaseOrder(ctx, c)
		k.setPurchaseOrderAgreement(ctx, c.AgreementID, c.PurchaseOrderID)
	}
	k.setPurchaseOrderID(ctx, uint64(len(data.PurchaseOrders)+1))
	k.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{
		PurchaseOrderss: k.PurchaseOrders(ctx),
		Params:          k.GetParams(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	if data.Params.MinPurchaseOrderLength < 1 {
		return fmt.Errorf("Param: MinPurchaseOrderLength must have a positive value")
	}
	if data.Params.MaxPurchaseOrderLength < 1 {
		return fmt.Errorf("Param: MaxPurchaseOrderLength must have a positive value")
	}

	return nil
}
