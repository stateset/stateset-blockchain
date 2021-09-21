package agreement

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stateset/stateset-blockchain/x/agreement/keeper"
	"github.com/stateset/stateset-blockchain/x/agreement/types"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		AgreementList: []*Agreement{},
	}
}


func (gs GenesisState) Validate() error {

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
	params, _ := k.GetParams(ctx)
	genesis := types.DefaultGenesis()
	genesis.Agreements = k.GetAgreements(ctx)
	genesis.PortId = k.GetPort(ctx)

	return genesis
}