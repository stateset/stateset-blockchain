package marketplace

const (
	ModuleName   = types.ModuleName
	RouterKey    = types.RouterKey
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper    = keeper.NewKeeper
	BeginBlocker = keeper.BeginBlocker
)

type (
	Keeper = keeper.Keeper
)