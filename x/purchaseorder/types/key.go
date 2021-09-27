package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "purchaseorder"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName
)

var (
	// KeyNextGlobalPurchaseOrderNumber defines key to store the next Purchase Order ID to be used
	KeyNextGlobalPurchaseOrderNumber = []byte{0x01}
	// KeyPrefixPurchasOrders defines prefix to store purchaseorders
	KeyPrefixPurchaseOrders = []byte{0x02}
	// KeyTotalLiquidity defines key to store total liquidity
	KeyTotalLiquidity = []byte{0x03}
)

func GetPurchaseOrderShareDenom(purchaseOrderId uint64) string {
	return fmt.Sprintf("purchaseorder/%d", purchaseOrderId)
}

func GetKeyPrefixPurchaseOrders(purchaseOrderId uint64) []byte {
	return append(KeyPrefixPurchaseOrders, sdk.Uint64ToBigEndian(purchaseOrderId)...)
}