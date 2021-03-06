package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCodec encodes module codec
var ModuleCodec *codec.LegacyAmino

const (
	ModuleName        = "statesetbank"
	StoreKey          = ModuleName
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName
	DefaultParamspace = ModuleName

	AttributeRecipient = "recipient"
)
