package bank

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(c *codec.LegacyAmino) {
	c.RegisterConcrete(MsgSend{}, "stateset/MsgSend", nil)
	c.RegisterConcrete(MsgIssue{}, "stateset/MsgIssue", nil)
	c.RegisterConcrete(MsgBurn{}, "stateset/MsgBurn", nil)
	c.RegisterConcrete(MsgUpdateParams{}, "bank/MsgUpdateParams", nil)

	c.RegisterConcrete(Transaction{}, "stateset/Transaction", nil)
}

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}