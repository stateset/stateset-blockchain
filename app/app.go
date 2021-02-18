package app

import (
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/stateset/stateset-blockchain/x/auth"
	"github.com/stateset/stateset-blockchain/x/auth/ante"
	"github.com/stateset/stateset-blockchain/x/auth/banking"
	"github.com/stateset/stateset-blockchain/x/bank"
	"github.com/stateset/stateset-blockchain/x/crisis"
	"github.com/stateset/stateset-blockchain/x/market"
	"github.com/stateset/stateset-blockchain/x/mint"
	"github.com/stateset/stateset-blockchain/x/agreement"
	"github.com/stateset/stateset-blockchain/x/invoice"
	"github.com/stateset/stateset-blockchain/x/purchaseorder"
	"github.com/stateset/stateset-blockchain/x/loan"
	"github.com/stateset/stateset-blockchain/x/factoring"
	distr "github.com/stateset/stateset-blockchain/x/distribution"
	"github.com/stateset/stateset-blockchain/x/genutil"
	"github.com/stateset/stateset-blockchain/x/params"
	"github.com/stateset/stateset-blockchain/x/slashing"
	"github.com/stateset/stateset-blockchain/x/liquidity"
	"github.com/stateset/stateset-blockchain/x/staking"
	"github.com/stateset/stateset-blockchain/x/supply"
	"github.com/stateset/stateset-blockchain/x/wasm"

)

const appName = "StatesetApp"

var (
	// default home directories for stateset binaries
	DefaultCLIHome = os.ExpandEnv("$HOME/.statesetcli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.statesetd")

	// ModuleBasics The module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distr.ProposalHandler, upgradeclient.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		liquidity.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		//stateset modules
		market.AppModuleBasic{},
		agreement.AppModuleBasic{},
		purchasorder.AppModuleBasic{},
		invoice.AppModuleBasic{},
		loan.AppModuleBasic{},
		// ...
		capability.AppModuleBasic{},
		ibc.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},

	)


	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		slashing.ModuleName:          {supply.Minter},
		slashing.PenaltyAccount:      nil,
		gov.ModuleName:            {supply.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
	}
)

// MakeCodec creates the application codec. The codec is sealed before it is
// returned.
func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)
	authvesting.RegisterCodec(cdc)

	return cdc.Seal()
}

// StatesetApp extended ABCI application
type StatesetApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tKeys map[string]*sdk.TransientStoreKey

	// cosmos keepers
	AccountKeeper  auth.AccountKeeper
	BankKeeper     bank.Keeper
	SupplyKeeper   supply.Keeper
	StakingKeeper  staking.Keeper
	SlashingKeeper slashing.Keeper
	DistrKeeper    distr.Keeper
	ParamsKeeper   params.Keeper
	GovKeeper      gov.Keeper
	CrisisKeeper   crisis.
	LiquidityKeeper liquidity.Keeper
	AppAccountKeeper   account.Keeper
	MarketKeeper  		market.Keeper
	MintKeeper     mint.Keeper
	AgreementKeeper    agreement.Keeper
	PurchaseorderKeeper purchaseorder.Keeper
	InvoiceKeeper      invoice.Keeper
	LoanKeeper         loan.Keeper
	FactoringKeeper    factoring.Keeper
	WasmKeeper     wasm.Keeper

	// other keepers
	IBCKeeper        *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
    EvidenceKeeper   evidencekeeper.Keeper // required to set up the client misbehaviour route
	TransferKeeper   ibctransferkeeper.Keeper // for cross-chain fungible token transfers
	
	 // make scoped keepers public for test purposes
	 ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	 ScopedTransferKeeper capabilitykeeper.ScopedKeeper


	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

// NewStatesetApp returns a reference to an initialized Stateset.
func NewStatesetApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp), ) *StatesetApp {
	
	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)
	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey, auth.StoreKey, staking.StoreKey,
		supply.StoreKey, mint.StoreKey, distr.StoreKey, slashing.StoreKey,
		gov.StoreKey, params.StoreKey, evidence.StoreKey, upgrade.StoreKey,
	)
	tKeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	app := &StatesetApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tKeys:          tKeys,
		
	}

	// add capability keeper and ScopeToModule for ibc module
  	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])

  	// grant capabilities for the ibc and ibc-transfer modules
  	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
  	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)

	// init params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tKeys[params.TStoreKey], params.DefaultCodespace)
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
	mintSubspace := app.paramsKeeper.Subspace(mint.DefaultParamspace)
	distrSubspace := app.paramsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	govSubspace := app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	crisisSubspace := app.paramsKeeper.Subspace(crisis.DefaultParamspace)
	evidenceSubspace := app.paramsKeeper.Subspace(evidence.DefaultParamspace)
	wasmSubspace = app.paramsKeeper.Subspace(wasm.DefaultParamspace)


	// stateset subspaces
	marketSubspace := app.paramsKeeper.Subspace(market.DefaultParamspace)
	agreementSubspace := app.paramsKeeper.Subspace(agreement.DefaultParamspace)
	purchaseorderSubspace := app.paramsKeeper.Subspace(purchaseorder.DefaultParamspace)
	invoiceSubspace := app.paramsKeeper.Subspace(invoice.DefaultParamspace)
	loanSubspace := app.paramsKeeper.Subspace(loan.DefaultParamspace)

	// Create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
	appCodec, keys[ibchost.StoreKey], app.StakingKeeper, scopedIBCKeeper,
	)

	// add cosmos keepers

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		keys[auth.StoreKey],
		authSubspace,
		auth.ProtoBaseAccount,
	)

	// The BankKeeper allows you to perofrm sdk.Coins interactions
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		bankSubspace,
		bank.DefaultCodespace,
		app.ModuleAccountAddrs(),
	)


	// The SupplyKeeper collects transaction fees and renders them to the fee distribution module
	app.supplyKeeper = supply.NewKeeper(
		app.cdc,
		keys[supply.StoreKey], 
		app.accountKeeper,
		 app.bankKeeper,
		maccPerms,
	)

	// The staking keeper
	stakingKeeper := staking.NewKeeper(
		app.cdc,
		keys[staking.StoreKey],
		app.supplyKeeper,
		stakingSubspace,
		staking.DefaultCodespace,
	)

	// The mint keeper
	app.mintKeeper = mint.NewKeeper(
		app.cdc,
		keys[mint.StoreKey],
		mintSubspace,
		&stakingKeeper,
		app.supplyKeeper,
		auth.FeeCollectorName,
	)
	
	app.distrKeeper = distr.NewKeeper(
		app.cdc,
		keys[distr.StoreKey],
		distrSubspace,
		&stakingKeeper,
		app.supplyKeeper,
		distr.DefaultCodespace,
		auth.FeeCollectorName,
		app.ModuleAccountAddrs(),
	)


	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		keys[slashing.StoreKey],
		&stakingKeeper,
		slashingSubspace,
		slashing.DefaultCodespace,
	)

	app.crisisKeeper = crisis.NewKeeper(
		crisisSubspace,
		invCheckPeriod,
		app.supplyKeeper,
		auth.FeeCollectorName,
	)

	// Create Transfer Keepers
    app.TransferKeeper = ibctransferkeeper.NewKeeper(
    appCodec, keys[ibctransfertypes.StoreKey],
    app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
	app.AccountKeeper, app.BankKeeper, scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(app.TransferKeeper)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
  	evidenceKeeper := evidencekeeper.NewKeeper(
    appCodec, keys[evidencetypes.StoreKey], &app.StakingKeeper, app.SlashingKeeper,
  	)
	
	app.upgradeKeeper = upgrade.NewKeeper(
		keys[upgrade.StoreKey],
		 app.cdc,
	)

	// create evidence keeper with evidence router
	evidenceKeeper := evidence.NewKeeper(
		app.cdc,
		keys[evidence.StoreKey],
		evidenceSubspace,
		evidence.DefaultCodespace,
		&stakingKeeper,
		app.slashingKeeper,
	)

	evidenceRouter := evidence.NewRouter()

	// TODO: register evidence routes
	evidenceKeeper.SetRouter(evidenceRouter)

	app.evidenceKeeper = *evidenceKeeper

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewMarketPoolSpendProposalHandler(app.distrKeeper)).
		AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.upgradeKeeper))
	app.govKeeper = gov.NewKeeper(
		app.cdc, keys[gov.StoreKey], govSubspace,
		app.supplyKeeper, &stakingKeeper, gov.DefaultCodespace, govRouter,
	)

	ibcRouter := port.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferModule)
	// Setting Router will finalize all routes by sealing router
	// No more routes can be added
	app.IBCKeeper.SetRouter(ibcRouter)
  
	// create static Evidence routers
  
	evidenceRouter := evidencetypes.NewRouter().
	  // add IBC ClientMisbehaviour evidence handler
	  AddRoute(ibcclient.RouterKey, ibcclient.HandlerClientMisbehaviour(app.IBCKeeper.ClientKeeper))
  
	// Setting Router will finalize all routes by sealing router
	// No more routes can be added
	evidenceKeeper.SetRouter(evidenceRouter)
  
	// set the evidence keeper from the section above
	app.EvidenceKeeper = *evidenceKeeper

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)

	// add stateset keepers

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		crisis.NewAppModule(&app.crisisKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.supplyKeeper),
		mint.NewAppModule(app.mintKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		liquidity.NewAppModule(app.liquidityKeeper, app)
		upgrade.NewAppModule(app.upgradeKeeper),
		evidence.NewAppModule(app.evidenceKeeper),
		market.NewAppModule(app.marketKeeper),
		agreement.NewAppModule(app.agreementKeeper),
		purchaseorder.NewAppModule(app.purchaseorderKeeper),
		invoice.NewAppModule(app.invoiceKeeper),
		loan.NewAppModule(app.loanKeeper),
		wasm.NewAppModule(app.wasmKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		transferModule,
		
	)


	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.

	app.mm.SetOrderBeginBlockers(upgrade.ModuleName, mint.ModuleName, distr.ModuleName, slashing.ModuleName, evidencetypes.ModuleName, stakingtypes.ModuleName, ibchost.ModuleName,
	)

	app.mm.SetOrderEndBlockers(crisis.ModuleName, gov.ModuleName, staking.ModuleName)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		distr.ModuleName, staking.ModuleName, auth.ModuleName, bank.ModuleName,
		slashing.ModuleName, gov.ModuleName, mint.ModuleName, supply.ModuleName,
		crisis.ModuleName, genutil.ModuleName, evidence.ModuleName, ibchost.ModuleName,
		evidencetypes.ModuleName, ibctransfertypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: This is not required for apps that don't use the simulator for fuzz testing
	// transactions.
	app.sm = module.NewSimulationManager(
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.supplyKeeper),
		mint.NewAppModule(app.mintKeeper),
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),
		wasm.NewAppModule(app.wasmKeeper, app.accountKeeper, app.bankKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		transferModule
	)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.supplyKeeper, auth.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			cmn.Exit(err.Error())
		}
	}

	return app
}

// BeginBlocker application updates every begin block
func (app *StatesetApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *StatesetApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *StatesetApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return app.mm.InitGenesis(ctx, genesisState)
}

// LoadHeight loads a particular height
func (app *StatesetApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}


func (app *StatesetApp) replayToHeight(replayHeight int64, logger log.Logger) int64 {
	loadHeight := int64(0)
	logger.Info("Please make sure the replay height is smaller than the latest block height.")
	if replayHeight >= DefaultSyncableHeight {
		loadHeight = replayHeight - replayHeight%DefaultSyncableHeight
	} else {
		// version 1 will always be kept
		loadHeight = 1
	}
	return loadHeight
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *StatesetApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// Codec returns the application's sealed codec.
func (app *StatesetApp) Codec() *codec.Codec {
	return app.cdc
}

// GetMaccPerms returns a mapping of the application's module account permissions.
func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}
	return modAccPerms
}