package loan

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	Loans []Loan `json:"loans"`
	Params Params  `json:"params"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{
		Loans: nil,
		Params: DefaultParams(),
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState { return NewGenesisState() }

// InitGenesis initializes stateset from genesis file
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	for _, c := range data.Loans {
		k.setLoan(ctx, c)
		k.setmarketLoan(ctx, c.MarketID, c.ID)
		k.setLenderLoan(ctx, c.Lender, c.ID)
		k.setCreatedTimeLoan(ctx, c.CreatedTime, c.ID)
	}
	k.setLoanID(ctx, uint64(len(data.Loans)+1))
	k.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{
		Loans: k.Loans(ctx),
		Params: k.GetParams(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	if data.Params.MinLoanLength < 1 {
		return fmt.Errorf("Param: MinLoanLength must have a positive value")
	}
	if data.Params.MaxLoanLength < 1 {
		return fmt.Errorf("Param: MaxLoanLength must have a positive value")
	}

	return nil
}