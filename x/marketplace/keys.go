package marketplace

// Keys for marketplace store
// Items are stored with the following key: values
//
// - 0x00<marketplaceID_Bytes>: Marketplace{} bytes
var (
	MarketplaceKeyPrefix = []byte{0x00}
)

// key for getting a specific marketplace from the store
func key(id string) []byte {
	return append(MarketplaceKeyPrefix, []byte(id)...)
}