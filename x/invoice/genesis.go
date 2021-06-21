package invoice

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stateset/stateset-blockchain/x/invoice/keeper"
	"github.com/stateset/stateset-blockchain/x/invoice/types"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	Invoices []Invoice `json:"invoices"`
	Params         Params          `json:"params"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{
		Invoices: nil,
		Params:         DefaultParams(),
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState { return NewGenesisState() }

// InitGenesis initializes stateset from genesis file
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	for _, c := range data.Invoices {
		k.setInvoices(ctx, c)
		k.setInvoiceAgreement(ctx, c.AgreementID, c.InvoiceID)
	}
	k.setInvoiceID(ctx, uint64(len(data.Invoices)+1))
	k.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{
		Invoices: k.Invoices(ctx),
		Params:          k.GetParams(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	if data.Params.MinInvoiceLength < 1 {
		return fmt.Errorf("Param: MinInvoiceLength must have a positive value")
	}
	if data.Params.MaxInvoiceLength < 1 {
		return fmt.Errorf("Param: MaxInvoiceLength must have a positive value")
	}

	return nil
}
