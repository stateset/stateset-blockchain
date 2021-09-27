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
	cdc.RegisterConcrete(&MsgUpdatePurchaseOrder{}, "stateset/MsgUpdatePurchaseOrder", nil)
	cdc.RegisterConcrete(&MsgDeletePurchaseOrder{}, "stateset/MsgDeletePurchaseOrder", nil)
	cdc.RegisterConcrete(&MsgCompletePurchaseOrder{}, "stateset/MsgCompletePurchaseOrder", nil)
	cdc.RegisterConcrete(&MsgCancelPurchaseOrder{}, "stateset/MsgCancelPurchaseOrder", nil)
	cdc.RegisterConcrete(&MsgLockPurchaseOrder{}, "stateset/MsgLockPurchaseOrder", nil)
	cdc.RegisterConcrete(&MsgFinancePurchaseOrder{}, "stateset/MsgFinancePurchaseOrder", nil)
	cdc.RegisterConcrete(&MsgSendIbcPurchaseOrder{}, "stateset/MsgSendIbcPurchaseOrder", nil)

}

// RegisterInterfaces registers the x/purchaseorder interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {

	registry.RegisterInterface(
		"stateset.purchaseorder.v1alpha1.PurchaseOrder",
		(*PurchaseOrderI)(nil),
		&PurchaseOrder{},
	)
	
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreatePurchaseOrder{},
		&MsgEditPurchaseOrder{},
		&MsgDeletePurchaseOrder{},
		&MsgCompletePurchaseOrder{},
		&MsgCancelPurchaseOrder{},
		&MsgLockPurchaseOrder{},
		&MsgFinancePurchaseOrder{},
		&MsgSendIbcPurchaseOrder{},
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
