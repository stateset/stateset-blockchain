package contact

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(MsgCreateContact{}, "stateset/MsgCreateContact", nil)
	c.RegisterConcrete(MsgEditContact{}, "stateset/MsgEditContact", nil)
	c.RegisterConcrete(MsgDeleteContact{}, "stateset/MsgDeleteContact", nil)
	c.RegisterConcrete(MsgAddAdmin{}, "contact/MsgAddAdmin", nil)
	c.RegisterConcrete(MsgRemoveAdmin{}, "contact/MsgRemoveAdmin", nil)
	c.RegisterConcrete(MsgUpdateParams{}, "contact/MsgUpdateParams", nil)

	c.RegisterConcrete(Contact{}, "stateset/Contact", nil)
}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}