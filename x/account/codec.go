package account

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(MsgCreateAccount{}, "stateset/MsgCreateAccount", nil)
	c.RegisterConcrete(MsgEditAccount{}, "stateset/MsgEditAccount", nil)
	c.RegisterConcrete(MsgDeleteAccount{}, "stateset/MsgDeleteAccount", nil)
	c.RegisterConcrete(MsgAddAdmin{}, "account/MsgAddAdmin", nil)
	c.RegisterConcrete(MsgRemoveAdmin{}, "account/MsgRemoveAdmin", nil)
	c.RegisterConcrete(MsgUpdateParams{}, "account/MsgUpdateParams", nil)

	c.RegisterConcrete(Account{}, "stateset/Account", nil)
}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}