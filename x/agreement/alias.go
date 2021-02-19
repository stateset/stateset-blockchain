package agreement

import (
	"github.com/stateset/stateset-blockchain/x/agreement/keeper"
	"github.com/stateset/stateset-blockchain/x/agreement/types"
)

const (
	// StoreKey represents storekey of agreement module
	StoreKey = types.StoreKey
	// ModuleName represents current module name
	ModuleName = types.ModuleName
)

type (
	// Keeper defines keeper of agreement module
	Keeper = keeper.Keeper
)

var (
	// NewKeeper creates new keeper instance of agreement module
	NewKeeper = keeper.NewKeeper
)