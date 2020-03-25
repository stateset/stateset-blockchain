package marketplace

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers messages into the codec
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(MsgNewMarketplace{}, "marketplace/MsgNewMarketplace", nil)
	c.RegisterConcrete(MsgAddItem{}, "marketplace/MsgAddItem", nil)
	c.RegisterConcrete(MsgCancelItem{], "marketplace/MsgCancelItem", nil)
	c.RegisterConcrete(MsgAddAdmin{}, "marketplace/MsgAddAdmin", nil)
	c.RegisterConcrete(MsgRemoveAdmin{}, "marketplace/MsgRemoveAdmin", nil)
	c.RegisterConcrete(MsgUpdateParams{}, "marketplace/MsgUpdateParams", nil)
}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)

	
	ModuleCodec.Seal()
}