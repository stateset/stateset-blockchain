package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types for codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "stateset-blockchain/staking/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "stateset-blockchain/staking/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgDestroyValidator{}, "stateset-blockchain/staking/MsgDestroyValidator", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "stateset-blockchain/staking/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgWithdraw{}, "stateset-blockchain/staking/MsgWithdraw", nil)
	cdc.RegisterConcrete(MsgAddShares{}, "stateset-blockchain/staking/MsgAddShares", nil)
	cdc.RegisterConcrete(MsgRegProxy{}, "stateset-blockchain/staking/MsgRegProxy", nil)
	cdc.RegisterConcrete(MsgBindProxy{}, "stateset-blockchain/staking/MsgBindProxy", nil)
	cdc.RegisterConcrete(MsgUnbindProxy{}, "stateset-blockchain/staking/MsgUnbindProxy", nil)
}

// ModuleCdc is generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
