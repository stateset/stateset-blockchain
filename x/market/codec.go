package market

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers messages into the codec
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(MsgNewMarket{}, "market/MsgNewMarket", nil)
	c.RegisterConcrete(MsgAddItem{}, "market/MsgAddItem", nil)
	c.RegisterConcrete(MsgCancelItem{}, "market/MsgCancelItem", nil)
	c.RegisterConcrete(MsgAddAdmin{}, "market/MsgAddAdmin", nil)
	c.RegisterConcrete(MsgRemoveAdmin{}, "market/MsgRemoveAdmin", nil)
	c.RegisterConcrete(MsgUpdateParams{}, "market/MsgUpdateParams", nil)
}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)

	
	ModuleCodec.Seal()
}