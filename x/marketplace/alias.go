package marketplace

import (
	"github.com/stateset/stateset-blockchain/x/marketplace/keeper"
	"github.com/stateset/stateset-blockchain/x/market/types"
)

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