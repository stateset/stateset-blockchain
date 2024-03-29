package agreement

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	porttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/05-port/types"
	"github.com/stateset/stateset-blockchain/x/agreement/client/cli"
	"github.com/stateset/stateset-blockchain/x/agreement/keeper"
	"github.com/stateset/stateset-blockchain/x/agreement/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
	_ porttypes.IBCModule   = AppModule{}
)

// ModuleName is the name of this module
const ModuleName = "agreement"

// AppModuleBasic defines the internal data for the module
// ----------------------------------------------------------------------------
type AppModuleBasic struct{}

// Name define the name of the module
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the types needed for amino encoding/decoding
func (AppModuleBasic) RegisterCodec(cdc *codec.LegacyAmino) {
	RegisterCodec(cdc)
}

// DefaultGenesis creates the default genesis state for testing
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCodec.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis validates the genesis state
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCodec.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the supply module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {

}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModuleBasic) RegisterGRPCRoutes(ctx client.Context, *runtime.ServeMux) {

}

// GetTxCmd returns the root tx command for the supply module.
func (AppModuleBasic) GetTxCmd(_ *codec.LegacyAmino) *cobra.Command { return nil }

// GetQueryCmd returns no root query command for the supply module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.LegacyAmino) *cobra.Command {
	return nil
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule defines external data for the module
// ----------------------------------------------------------------------------
type AppModule struct {
	AppModuleBasic

	keeper Keeper
}

// NewAppModule creates a NewAppModule object
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

// RegisterInvariants enforces registering of invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
}

// Route defines the key for the route
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler creates the handler for the module
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute defines the querier route
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler creates a new querier handler
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(Configurator)

// InitGenesis enforces the creation of the genesis state for this module
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCodec.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis enforces exporting this module's data to a genesis file
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCodec.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the purchaseorder module.
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the purchaseorder module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}