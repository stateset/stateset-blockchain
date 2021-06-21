package agreement

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stateset/stateset-blockchain/x/agreement/keeper"
	"github.com/stateset/stateset-blockchain/x/agreement/types"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # ibc/genesistype/default
		// this line is used by starport scaffolding # genesis/types/default
		AgreementList: []*Agreement{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # ibc/genesistype/validate

	// this line is used by starport scaffolding # genesis/types/validate
	// Check for duplicated ID in agreement
	agreementIdMap := make(map[uint64]bool)

	for _, elem := range gs.AgreementList {
		if _, ok := agreementIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for agreement")
		}
		agreementIdMap[elem.Id] = true
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