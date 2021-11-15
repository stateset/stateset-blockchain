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
	cdc.RegisterConcrete(&MsgRequestLoan{}, "stateset/MsgRequestLoan", nil)
	cdc.RegisterConcrete(&MsgApproveLoan{}, "stateset/MsgApproveLoan", nil)
	cdc.RegisterConcrete(&MsgBorrowLoan{}, "stateset/MsgBorrowLoan", nil)
	cdc.RegisterConcrete(&MsgCancelLoan{}, "stateset/MsgCancelLoan", nil)
	cdc.RegisterConcrete(&MsgRepayLoan{}, "stateset/MsgRepayLoan", nil)
	cdc.RegisterConcrete(&MsgLiqudidateLoan{}, "stateset/MsgLiqudidateLoan", nil)
	cdc.RegisterConcrete(&MsgDepositCollateral{}, "stateset/MsgDepositCollateral", nil)
	cdc.RegisterConcrete(&MsgWithdrawCollateral{}, "stateset/MsgWithdrawCollateral", nil)

}

// RegisterInterfaces registers the x/loan interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {

	registry.RegisterInterface(
		"stateset.purchaseorder.v1alpha1.PurchaseOrder",
		(*PurchaseOrderI)(nil),
		&PurchaseOrder{},
	)
	
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRequestLoan{},
		&MsgApproveLoan{},
		&MsgBorrowLoan{},
		&MsgCancelLoan{},
		&MsgRepayLoan{},
		&MsgLiquidateLoan{},
		&MsgDepositCollateral{},
		&MsgWithdrawCollateral{},
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
