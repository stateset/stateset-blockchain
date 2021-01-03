package marketplace

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	Marketplaces []Marketplace `json:"marketplaces"`
	Params      Params      `json:"params"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(marketplaces []Marketplace, params Params) GenesisState {
	return GenesisState{
		Marketplaces: marketplaces,
		Params:      params,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	marketplace1 := NewMarketplace("crypto", "Cryptocurrency", "description string", time.Now())
	marketplace2 := NewMarketplace("meme", "Memes", "description string", time.Now())

	return NewGenesisState(Marketplaces{marketplace1, marketplace2}, DefaultParams())
}

// InitGenesis initializes marketplace state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, marketplace := range data.Marketplaces {
		keeper.setMarketplace(ctx, marketplace)
	}
	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Marketplaces: keeper.Marketplaces(ctx),
		Params:      keeper.GetParams(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	if data.Params.MinNameLength < 1 {
		return fmt.Errorf("Param: MinNameLength, must have a positive value")
	}

	if data.Params.MaxNameLength < 1 || data.Params.MaxNameLength < data.Params.MinNameLength {
		return fmt.Errorf("Param: MaxNameLength, must have a positive value and be larger than MinNameLength")
	}

	if data.Params.MinIDLength < 1 {
		return fmt.Errorf("Param: MinIDLength, must have a positive value")
	}

	if data.Params.MaxIDLength < 1 || data.Params.MaxIDLength < data.Params.MinIDLength {
		return fmt.Errorf("Param: MaxIDLength, must have a positive value and be larger than MinIDLength")
	}

	if data.Params.MaxDescriptionLength < 1 {
		return fmt.Errorf("Param: MaxDescriptionLength, must have a positive value")
	}

	if len(data.Params.MarketplaceAdmins) < 1 {
		return fmt.Errorf("Param: MarketplaceAdmins, must have atleast one admin")
	}

	return nil
}