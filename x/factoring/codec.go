package factoring

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(MsgFactorInvoice{}, "stateset/MsgFactorInvoice", nil)
	c.RegisterConcrete(MsgAddAdmin{}, "factoring/MsgAddAdmin", nil)
	c.RegisterConcrete(MsgRemoveAdmin{}, "factoring/MsgRemoveAdmin", nil)
	c.RegisterConcrete(MsgUpdateParams{}, "factoring/MsgUpdateParams", nil)

	c.RegisterConcrete(Factor{}, "stateset/Factor", nil)
	c.RegisterConcrete(Loan{}, "stateset/Loan", nil)
	c.RegisterConcrete(Invoice{}, "stateset/Invoice", nil)

}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}