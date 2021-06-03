package keeper

import (
	"github.com/stateset/stateset-blockchain/x/purchaseorder/types"
)

var _ types.QueryServer = Keeper{}
