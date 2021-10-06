package keeper

import (
	"github.com/stateset/stateset-blockchain/x/invoice/types"
)

var _ types.QueryServer = Keeper{}
