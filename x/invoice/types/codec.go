package invoice

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgCreateInvoice{}, "stateset/MsgCreateInvoice", nil)
	cdc.RegisterConcrete(MsgCancelInvoice{}, "stateset/MsgCancelInvoice", nil)
	cdc.RegisterConcrete(MsgEditInvoice{}, "stateset/MsgEditInvoice", nil)
	cdc.RegisterConcrete(MsgDeleteInvoice{}, "stateset/MsgDeleteInvoice", nil)
	cdc.RegisterConcrete(MsgFactorInvoice{}, "stateset/MsgFactorInvoice", nil)
	cdc.RegisterConcrete(MsgPayInvoice{}, "stateset/MsgPayInvoice", nil)
	cdc.RegisterConcrete(MsgAddAdmin{}, "stateset/MsgAddAdmin", nil)
	cdc.RegisterConcrete(MsgRemoveAdmin{}, "stateset/MsgRemoveAdmin", nil)
	cdc.RegisterConcrete(MsgUpdateParams{}, "stateset/MsgUpdateParams", nil)

	cdc.RegisterConcrete(Invoice{}, "stateset/Invoice", nil)
}

// RegisterInterfaces registers the x/purchaseorder interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {

	registry.RegisterInterface(
		"stateset.invoice.v1alpha1.Invoice",
		(*InvoiceI)(nil),
		&Invoice{},
	)
	
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateInvoice{},
		&MsgEditInvoice{},
		&MsgDeleteInvoice{},
		&MsgCompleteInvoice{},
		&MsgCancelInvoice{},
		&MsgPayPurchaseOrder{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
