package agreement

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(MsgCreateAgreement{}, "stateset/MsgCreateAgreement", nil)
	c.RegisterConcrete(MsgEditAgreement{}, "stateset/MsgEditAgreement", nil)
	c.RegisterConcrete(MsgAmendAgreement{}, "stateset/MsgAmendAgreement", nil)
	c.RegisterConcrete(MsgRenewAgreement{}, "stateset/MsgRenewAgreement", nil)
	c.RegisterConcrete(MsgTerminateAgreement{}, "stateset/MsgTerminateAgreement", nil)
	c.RegisterConcrete(MsgExpireAgreement{}, "stateset/MsgExpireAgreement", nil)

	c.RegisterConcrete(Agreement{}, "stateset/Agreement", nil)
}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}