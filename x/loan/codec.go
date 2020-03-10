package loan

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(MsgCreateLoan{}, "stateset/MsgCreateLoan", nil)
	c.RegisterConcrete(MsgEditLoan{}, "stateset/MsgEditLoan", nil)
	c.RegisterConcrete(MsgDeleteLoan{}, "stateset/MsgDeleteLoan", nil)
	c.RegisterConcrete(MsgPayBackLoan{}, "stateset/MsgPayBackLoan", nil)
	c.RegisterConcrete(MsgAddAdmin{}, "loan/MsgAddAdmin", nil)
	c.RegisterConcrete(MsgRemoveAdmin{}, "loan/MsgRemoveAdmin", nil)
	c.RegisterConcrete(MsgUpdateParams{}, "loan/MsgUpdateParams", nil)

	c.RegisterConcrete(Loan{}, "stateset/Loan", nil)
}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}