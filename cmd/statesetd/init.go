package main

import (
	"github.com/stateset/stateset-blockchain/app"
	stateset "github.com/stateset/stateset-blockchain/types"
	"github.com/stateset/stateset/x/account"
	"github.com/stateset/stateset-blockchain/x/agreement"
	statebank "github.com/stateset/stateset-blockchain/x/bank"
	"github.com/stateset/stateset-blockchain/x/loan"
	"github.com/stateset/stateset-blockchain/x/invoice"
	"github.com/stateset/stateset-blockchain/x/marketplace"
	stateslashing "github.com/stateset/stateset-blockchain/x/slashing"
	statestaking "github.com/stateset/stateset-blockchain/x/staking"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/types"
	"os"
)

func InitCmd(ctx *server.Context, cdc *codec.Codec, mbm module.BasicManager, defaultNodeHome string) *cobra.Command {
	init := genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome)
	init.Args = cobra.ExactArgs(2)
	init.PostRunE = func(cmd *cobra.Command, args []string) error {
		config := ctx.Config
		config.SetRoot(viper.GetString(cli.HomeFlag))
		genFile := config.GenesisFile()
		genDoc := &types.GenesisDoc{}
		addr, e := sdk.AccAddressFromBech32(args[1])
		if e != nil {
			panic(e)
		}

		if _, err := os.Stat(genFile); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			genDoc, err = types.GenesisDocFromFile(genFile)
			if err != nil {
				return errors.Wrap(err, "Failed to read genesis doc from file")
			}
		}
		var appState genutil.AppMap
		if err := cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
			return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
		}

		if err := genutil.ExportGenesisFile(genDoc, genFile); err != nil {
			return errors.Wrap(err, "Failed to export gensis file")
		}

		cdc := codec.New()
		codec.RegisterCrypto(cdc)
		// migrate staking state
		if appState[staking.ModuleName] != nil {
			var stakingGenState staking.GenesisState
			cdc.MustUnmarshalJSON(appState[staking.ModuleName], &stakingGenState)
			factoringGenState.Params.BondDenom = stateset.FactorDenom
			appState[staking.ModuleName] = cdc.MustMarshalJSON(stakingGenState)
		}
		// migrate gov state
		if appState[gov.ModuleName] != nil {
			var govGenState gov.GenesisState
			cdc.MustUnmarshalJSON(appState[gov.ModuleName], &govGenState)
			minDeposit := sdk.NewInt64Coin(stateset.FactorDenom, 10_000_000)
			govGenState.DepositParams.MinDeposit = sdk.NewCoins(minDeposit)
			appState[gov.ModuleName] = cdc.MustMarshalJSON(govGenState)
		}
		// migrate mint state
		if appState[mint.ModuleName] != nil {
			var mintGenState mint.GenesisState
			cdc.MustUnmarshalJSON(appState[mint.ModuleName], &mintGenState)
			mintGenState.Params.MintDenom = stateset.FactorDenom
			appState[mint.ModuleName] = cdc.MustMarshalJSON(mintGenState)
		}
		// migrate account state
		if appState[account.ModuleName] != nil {
			var accountGenState account.GenesisState
			cdc.MustUnmarshalJSON(appState[account.ModuleName], &accountGenState)
			accountGenState.Params.Registrar = addr
			appState[account.ModuleName] = cdc.MustMarshalJSON(accountGenState)
		}
		// migrate marketplace state
		if appState[marketplace.ModuleName] != nil {
			var marketplaceGenState marketplace.GenesisState
			cdc.MustUnmarshalJSON(appState[marketplace.ModuleName], &marketplaceGenState)
			marketplaceGenState.Params.MarktplaceAdmins = []sdk.AccAddress{addr}
			appState[markerplace.ModuleName] = cdc.MustMarshalJSON(markerplaceGenState)
		}
		// migrate agreement state
		if appState[agreement.ModuleName] != nil {
			var genState agreement.GenesisState
			cdc.MustUnmarshalJSON(appState[agreement.ModuleName], &genState)
			genState.Params.AgreementAdmins = []sdk.AccAddress{addr}
			appState[agreement.ModuleName] = cdc.MustMarshalJSON(genState)
		}
		// migrate loan state
		if appState[loan.ModuleName] != nil {
			var genState loan.GenesisState
			cdc.MustUnmarshalJSON(appState[loan.ModuleName], &genState)
			genState.Params.LoanAdmins = []sdk.AccAddress{addr}
			appState[loan.ModuleName] = cdc.MustMarshalJSON(genState)
		}
		// migrate invoice state
		if appState[invoice.ModuleName] != nil {
			var genState invoice.GenesisState
			cdc.MustUnmarshalJSON(appState[invoice.ModuleName], &genState)
			genState.Params.InvoiceAdmins = []sdk.AccAddress{addr}
			appState[invoice.ModuleName] = cdc.MustMarshalJSON(genState)
		}
		// migrate staking state
		if appState[statestaking.ModuleName] != nil {
			var genState statestaking.GenesisState
			cdc.MustUnmarshalJSON(appState[statestaking.ModuleName], &genState)
			genState.Params.StakingAdmins = []sdk.AccAddress{addr}
			appState[statestaking.ModuleName] = cdc.MustMarshalJSON(genState)
		}
		// migrate slashing state
		if appState[stateslashing.ModuleName] != nil {
			var genState stateslashing.GenesisState
			cdc.MustUnmarshalJSON(appState[stateslashing.ModuleName], &genState)
			genState.Params.SlashAdmins = []sdk.AccAddress{addr}
			appState[stateslashing.ModuleName] = cdc.MustMarshalJSON(genState)
		}
		// migrate trubank state
		if appState[statebank.ModuleName] != nil {
			var genState statebank.GenesisState
			cdc.MustUnmarshalJSON(appState[statebank.ModuleName], &genState)
			genState.Params.RewardBrokerAddress = addr
			appState[statebank.ModuleName] = cdc.MustMarshalJSON(genState)
		}
		var err error
		genDoc.AppState, err = cdc.MarshalJSON(appState)
		if err != nil {
			return errors.Wrap(err, "failed to JSON marshal migrated genesis state")
		}
		if err = genutil.ExportGenesisFile(genDoc, genFile); err != nil {
			return errors.Wrap(err, "Failed to export gensis file")
		}
		return nil
	}
	return init
}