package agreement

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
func DefaultGenesisState() GenesisState { return NewGenesisState() }

// InitGenesis initializes stateset from genesis file
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	for _, c := range data.Agreements {
		k.setAgreement(ctx, c)
		k.setMarketplaceAgreement(ctx, c.MarketplaceID, c.AgreementID)
		k.setCounterpartyAgreement(ctx, c.Counterparty, c.AgreementID)
		k.setCreatedTimeAgreement(ctx, c.CreatedTime, c.AgreementID)
	}
	k.setAgreementID(ctx, uint64(len(data.Agreements)+1))
	k.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{
		Agreements: k.Agreements(ctx),
		Params: k.GetParams(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	if data.Params.MinAgreementLength < 1 {
		return fmt.Errorf("Param: MinAgreementLength must have a positive value")
	}
	if data.Params.MaxAgreementLength < 1 {
		return fmt.Errorf("Param: MaxAgreeemntLength must have a positive value")
	}

	return nil
}