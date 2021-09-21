package types

import (
	"fmt"

	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PortId: PortID,	
		PurchaseOrders:	[]*codectypes:Any{},
		Params:         DefaultParams(),
		TimedoutPurchaseOrderList: []*TimedoutPurchaseOrder{},
		SentPurchaseOrderList:     []*SentPurchaseOrder{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}

	// this line is used by starport scaffolding # genesis/types/validate
	// Check for duplicated ID in timedoutPurchaseOrder
	timedoutPurchaseOrderIdMap := make(map[uint64]bool)

	for _, elem := range gs.TimedoutPurchaseOrderList {
		if _, ok := timedoutPurchaseOrderIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for timedoutPurchaseOrder")
		}
		timedoutPurchaseOrderIdMap[elem.Id] = true
	}
	// Check for duplicated ID in sentPurchaseOrder
	sentPurchaseOrderIdMap := make(map[uint64]bool)

	for _, elem := range gs.SentPurchaseOrderList {
		if _, ok := sentPurchaseOrderIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for sentPurchaseOrder")
		}
		sentPurchaseOrderIdMap[elem.Id] = true
	}

	return nil
}
