package keeper

import (
	"github.com/stateset/stateset-blockchain/x/ibc/applications/ibcdex/types"
)

var _ types.QueryServer = Keeper{}
