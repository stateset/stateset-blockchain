package purchaseorder

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stateset/stateset-blockchain/x/purchaseorder/keeper"
	"github.com/stateset/stateset-blockchain/x/purchaseorder/types"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PurchaseOrderList: []*PurchaseOrder{},
	}
}


func (gs GenesisState) Validate() error {

	purchaseOrderIdMap := make(map[uint64]bool)

	for _, elem := range gs.PurchaseOrderList {
		if _, ok := purchaseOrderIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for po")
		}
		purchaseOrderIdMap[elem.Id] = true
	}

	return nil
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export

	genesis.PortId = k.GetPort(ctx)

	return genesis
}