package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

var (
	amino = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgCreateAgreement{}, "stateset/MsgCreateAgreement", nil)
	cdc.RegisterConcrete(MsgEditAgreement{}, "stateset/MsgEditAgreement", nil)
	cdc.RegisterConcrete(MsgAmendAgreement{}, "stateset/MsgAmendAgreement", nil)
	cdc.RegisterConcrete(MsgRenewAgreement{}, "stateset/MsgRenewAgreement", nil)
	cdc.RegisterConcrete(MsgTerminateAgreement{}, "stateset/MsgTerminateAgreement", nil)
	cdc.RegisterConcrete(MsgExpireAgreement{}, "stateset/MsgExpireAgreement", nil)

	c.RegisterConcrete(Agreement{}, "stateset/Agreement", nil)
}

// RegisterInterfaces registers the x/market interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateAgreement{},
		&MsgEditAgreement{},
		&MsgAmendAgreement{},
		&MsgRenewAgreement{},
		&MsgTerminateAgreement{},
		&MsgExpireAgreement{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}