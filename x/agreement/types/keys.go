package types

const (
	ModuleName = "agreement"

	StoreKey = ModuleName

	DefaultParamspace = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName

	MemStoreKey = "mem_agreement"

	// Version defines the current version the IBC module supports
	Version = "agreement-1"

	// PortID is the default port id that module binds to
	PortID = "agreement"
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("agreement-port-")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
