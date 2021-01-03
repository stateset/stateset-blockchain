package purchaseorder

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(MsgCreatePurchaseOrder{}, "stateset/MsgCreatePurchaseOrder", nil)
	c.RegisterConcrete(MsgEditPurchaseOrder{}, "stateset/MsgEditPurchaseOrder", nil)
	c.RegisterConcrete(MsgCompletePurchaseOrder{}, "stateset/MsgCompletePurchaseOrder", nil)
	c.RegisterConcrete(MsgCancelPurchaseOrder{}, "stateset/MsgCancelPurchaseOrder", nil)
	c.RegisterConcrete(MsgFinancePurchaseOrder{}, "stateset/MsgFinancePurchaseOrder", nil)

	c.RegisterConcrete(PurchaseOrder{}, "stateset/PurchaseOrder", nil)
}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}