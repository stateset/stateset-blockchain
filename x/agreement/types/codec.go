package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/bank interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
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
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateAgreement{},
		&MsgEditAgreement{},
		&MsgAmendAgreement{},
		&MsgRenewAgreement{},
		&MsgTerminateAgreement{},
		&MsgExpireAgreement{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}


var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/bank module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}