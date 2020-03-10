package marketplace

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers messages into the codec
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(MsgNewMarkerplace{}, "community/MsgNewMarkeplace", nil)
	c.RegisterConcrete(MsgAddAdmin{}, "community/MsgAddAdmin", nil)
	c.RegisterConcrete(MsgRemoveAdmin{}, "community/MsgRemoveAdmin", nil)
	c.RegisterConcrete(MsgUpdateParams{}, "community/MsgUpdateParams", nil)
}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}