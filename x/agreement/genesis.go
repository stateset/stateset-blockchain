package agreement

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stateset/stateset-blockchain/x/agreement/keeper"
	"github.com/stateset/stateset-blockchain/x/agreement/types"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	Agreements []Agreement `json:"agreements"`
	Params Params  `json:"params"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{
		Agreements: nil,
		Params: DefaultParams(),
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState(cdc codec.JSONMarshaler) GenesisState { return NewGenesisState() }

// InitGenesis initializes stateset from genesis file
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {

	k.SetPort(ctx, genState.PortId)
	
	for _, c := range data.Agreements {
		k.setAgreement(ctx, c)
	}
	k.setAgreementID(ctx, uint64(len(data.Agreements)+1))
	k.SetParams(ctx, data.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export

	genesis.PortId = k.GetPort(ctx)

	return genesis
}


// ValidateGenesis validates the genesis state data
func ValidateGenesis(cdc codec.JSONMarshaler, config client.TxEncodingConfig, bz json.RawMessage) error {
	if data.Params.MinAgreementLength < 1 {
		return fmt.Errorf("Param: MinAgreementLength must have a positive value")
	}
	if data.Params.MaxAgreementLength < 1 {
		return fmt.Errorf("Param: MaxAgreeemntLength must have a positive value")
	}

	return nil
}