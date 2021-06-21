package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers all the necessary agreements module concrete types and interfaces with
// the provided codec reference. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*Agreement)(nil), nil)
	cdc.RegisterConcrete(MsgCreateAgreement{}, "stateset/MsgCreateAgreement", nil)
	cdc.RegisterConcrete(MsgEditAgreement{}, "stateset/MsgEditAgreement", nil)
	cdc.RegisterConcrete(MsgEditAgreement{}, "stateset/MsgDeleteAgreement", nil)
	cdc.RegisterConcrete(MsgAmendAgreement{}, "stateset/MsgAmendAgreement", nil)
	cdc.RegisterConcrete(MsgRenewAgreement{}, "stateset/MsgRenewAgreement", nil)
	cdc.RegisterConcrete(MsgTerminateAgreement{}, "stateset/MsgTerminateAgreement", nil)
	cdc.RegisterConcrete(MsgExpireAgreement{}, "stateset/MsgExpireAgreement", nil)
	cdc.RegisterConcrete(&MsgSendIbcAgreement{}, "stateset/SendIbcAgreement", nil)

	c.RegisterConcrete(Agreement{}, "stateset/Agreement", nil)
}

// RegisterInterfaces registers the agreement interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateAgreement{},
		&MsgEditAgreement{},
		&MsgDeleteAgreement{},
		&MsgAmendAgreement{},
		&MsgRenewAgreement{},
		&MsgTerminateAgreement{},
		&MsgExpireAgreement{},
		&MsgSendIbcAgreement{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	registry.RegisterInterface(
		"stateset.ageeement.v1alpha1.Agreement",
		(*Agreement)(nil),
	)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
}
