package types

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

// RegisterCodec register concrete types on codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreatePurchaseOrder{}, "stateset/MsgCreatePurchaseOrder", nil)
	cdc.RegisterConcrete(&MsgEditPurchaseOrder{}, "stateset/MsgEditPurchaseOrder", nil)
	cdc.RegisterConcrete(&MsgDeletePurchaseOrder{}, "stateset/MsgDeletePurchaseOrder", nil)
	cdc.RegisterConcrete(&MsgCompletePurchaseOrder{}, "stateset/MsgCompletePurchaseOrder", nil)
	cdc.RegisterConcrete(&MsgCancelPurchaseOrder{}, "stateset/MsgCancelPurchaseOrder", nil)
	cdc.RegisterConcrete(&MsgFinancePurchaseOrder{}, "stateset/MsgFinancePurchaseOrder", nil)
}

// RegisterInterfaces registers the x/purchaseorder interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreatePurchaseOrder{},
		&MsgEditPurchaseOrder{},
		&MsgDeletePurchaseOrder{},
		&MsgCompletePurchaseOrder{},
		&MsgCancelPurchaseOrder{},
		&MsgFinancePurchaseOrder{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}