package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	appParams "github.com/stateset/stateset-blockchain/app/params"

)

func PrepareGenesisCmd(defaultNodeHome string, mbm module.BasicManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prepare-genesis",
		Short: "Prepare a genesis file with initial setup",
		Long: `Prepare a genesis file with initial setup.
Examples include:
	- Setting module initial params
	- Setting denom metadata
Example:
	statesetd prepare-genesis mainnet stateset-1
	- Check input genesis:
		file is at ~/.statesetd/config/genesis.json
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			depCdc := clientCtx.JSONMarshaler
			cdc := depCdc.(codec.Marshaler)
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			// read genesis file
			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			// get genesis params
			var genesisParams GenesisParams
			network := args[0]
			if network == "testnet" {
				genesisParams = TestnetGenesisParams()
			} else if network == "mainnet" {
				genesisParams = MainnetGenesisParams()
			} else {
				return fmt.Errorf("please choose 'mainnet' or 'testnet'")
			}

			// get genesis params
			chainID := args[1]

			// run Prepare Genesis
			appState, genDoc, err = PrepareGenesis(clientCtx, appState, genDoc, genesisParams, chainID)

			// validate genesis state
			if err = mbm.ValidateGenesis(cdc, clientCtx.TxConfig, appState); err != nil {
				return fmt.Errorf("error validating genesis file: %s", err.Error())
			}

			// save genesis
			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			err = genutil.ExportGenesisFile(genDoc, genFile)
			return err
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func PrepareGenesis(clientCtx client.Context, appState map[string]json.RawMessage, genDoc *tmtypes.GenesisDoc, genesisParams GenesisParams, chainID string) (map[string]json.RawMessage, *tmtypes.GenesisDoc, error) {
	depCdc := clientCtx.JSONMarshaler
	cdc := depCdc.(codec.Marshaler)

	// chain params genesis
	genDoc.GenesisTime = genesisParams.GenesisTime

	genDoc.ConsensusParams = genesisParams.ConsensusParams

	// ---
	// staking module genesis
	stakingGenState := stakingtypes.GetGenesisStateFromAppState(depCdc, appState)
	stakingGenState.Params = genesisParams.StakingParams
	stakingGenStateBz, err := cdc.MarshalJSON(stakingGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal staking genesis state: %w", err)
	}
	appState[stakingtypes.ModuleName] = stakingGenStateBz

	// mint module genesis
	mintGenState := minttypes.DefaultGenesisState()
	mintGenState.Params = genesisParams.MintParams
	mintGenStateBz, err := cdc.MarshalJSON(mintGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal mint genesis state: %w", err)
	}
	appState[minttypes.ModuleName] = mintGenStateBz

	// distribution module genesis
	distributionGenState := distributiontypes.DefaultGenesisState()
	distributionGenState.Params = genesisParams.DistributionParams
	// TODO Set initial community pool
	// distributionGenState.FeePool.CommunityPool = sdk.NewDecCoins()
	distributionGenStateBz, err := cdc.MarshalJSON(distributionGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal distribution genesis state: %w", err)
	}
	appState[distributiontypes.ModuleName] = distributionGenStateBz

	// gov module genesis
	govGenState := govtypes.DefaultGenesisState()
	govGenState.DepositParams = genesisParams.GovParams.DepositParams
	govGenState.TallyParams = genesisParams.GovParams.TallyParams
	govGenState.VotingParams = genesisParams.GovParams.VotingParams
	govGenStateBz, err := cdc.MarshalJSON(govGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal gov genesis state: %w", err)
	}
	appState[govtypes.ModuleName] = govGenStateBz

	// crisis module genesis
	crisisGenState := crisistypes.DefaultGenesisState()
	crisisGenState.ConstantFee = genesisParams.CrisisConstantFee
	// TODO Set initial community pool
	// distributionGenState.FeePool.CommunityPool = sdk.NewDecCoins()
	crisisGenStateBz, err := cdc.MarshalJSON(crisisGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal crisis genesis state: %w", err)
	}
	appState[crisistypes.ModuleName] = crisisGenStateBz

	// slashing module genesis
	slashingGenState := slashingtypes.DefaultGenesisState()
	slashingGenState.Params = genesisParams.SlashingParams
	slashingGenStateBz, err := cdc.MarshalJSON(slashingGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal slashing genesis state: %w", err)
	}
	appState[slashingtypes.ModuleName] = slashingGenStateBz

	// return appState and genDoc
	return appState, genDoc, nil
}

type GenesisParams struct {
	AirdropSupply sdk.Int

	StrategicReserveAccounts []banktypes.Balance

	ConsensusParams *tmproto.ConsensusParams

	GenesisTime         time.Time
	NativeCoinMetadatas []banktypes.Metadata

	StakingParams      stakingtypes.Params
	MintParams         minttypes.Params
	DistributionParams distributiontypes.Params
	GovParams          govtypes.Params

	CrisisConstantFee sdk.Coin

	SlashingParams    slashingtypes.Params
}

func MainnetGenesisParams() GenesisParams {
	genParams := GenesisParams{}

	genParams.AirdropSupply = sdk.NewIntWithDecimal(5, 13)                // 5*10^13 uparo, 5*10^7 (50 million) paro
	genParams.GenesisTime = time.Date(2021, 6, 18, 17, 0, 0, 0, time.UTC) // Jun 18, 2021 - 17:00 UTC

	genParams.NativeCoinMetadatas = []banktypes.Metadata{
		{
			Description: fmt.Sprintf("The native token of Stateset"),
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    appParams.BaseCoinUnit,
					Exponent: 0,
					Aliases:  nil,
				},
				{
					Denom:    appParams.HumanCoinUnit,
					Exponent: appParams.SsetExponent,
					Aliases:  nil,
				},
			},
			Base:    appParams.BaseCoinUnit,
			Display: appParams.HumanCoinUnit,
		},
		{
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    "uion",
					Exponent: 0,
					Aliases:  nil,
				},
				{
					Denom:    "ion",
					Exponent: 6,
					Aliases:  nil,
				},
			},
			Base:    "uion",
			Display: "ion",
		},
	}

	genParams.StrategicReserveAccounts = []banktypes.Balance{
		{
			Address: "paro1ekmcnlkmlksck218oiuoijioj1232121",
			Coins:   sdk.NewCoins(sdk.NewCoin(genParams.NativeCoinMetadatas[0].Base, sdk.NewInt(47_874_500_000_000))), // 47.8745 million OSMO
		},
		{
			Address: "paro1ekmcnlkmlksck218oiuoijioj1232121",
			Coins:   sdk.NewCoins(sdk.NewCoin(genParams.NativeCoinMetadatas[0].Base, sdk.NewInt(500_000_000_000))), // 500 thousand OSMO
		},
		{
			Address: "paro1ekmcnlkmlksck218oiuoijioj1232121",
			Coins:   sdk.NewCoins(sdk.NewCoin(genParams.NativeCoinMetadatas[0].Base, sdk.NewInt(1_000_000_000_000))), // 1 million OSMO
		},
		{
			Address: "paro16n7070n4whce0wlu76j42dyrxh9f7nlapg6c4a",
			Coins:   sdk.NewCoins(sdk.NewCoin(genParams.NativeCoinMetadatas[0].Base, sdk.NewInt(50_000_000_000))),
		},
		{
			Address: "paro1vnyc6q49sr0hs9ddjepcmtlaq3l6wwj0rrw6hd",
			Coins:   sdk.NewCoins(sdk.NewCoin(genParams.NativeCoinMetadatas[0].Base, sdk.NewInt(1_000_000_000))),
		},
		{
			Address: "paro1ujx3rqerqdksnxnjka2n9tde874mzut75hzx92",
			Coins:   sdk.NewCoins(sdk.NewCoin(genParams.NativeCoinMetadatas[0].Base, sdk.NewInt(1_000_000_000))),
		},
		{
			Address: "paro1a2r7nqnc9e032wj37ptskh207e232462vcrrjf",
			Coins:   sdk.NewCoins(sdk.NewCoin(genParams.NativeCoinMetadatas[0].Base, sdk.NewInt(1_000_000_000))),
		},
	}

	genParams.StakingParams = stakingtypes.DefaultParams()
	genParams.StakingParams.UnbondingTime = time.Hour * 24 * 7 * 2 // 2 weeks
	genParams.StakingParams.MaxValidators = 100
	genParams.StakingParams.BondDenom = genParams.NativeCoinMetadatas[0].Base
	genParams.StakingParams.MinCommissionRate = sdk.MustNewDecFromStr("0.05")

	genParams.MintParams = minttypes.DefaultParams()
	genParams.MintParams.EpochIdentifier = "day"                                                // 1 day
	genParams.MintParams.GenesisEpochProvisions = sdk.NewDec(300_000_000_000_000).QuoInt64(365) // 300M * 10^6 / 365 = ~821917.8082191781 * 10^6
	genParams.MintParams.MintDenom = genParams.NativeCoinMetadatas[0].Base
	genParams.MintParams.ReductionFactor = sdk.NewDec(2).QuoInt64(3) // 2/3
	genParams.MintParams.ReductionPeriodInEpochs = 365               // 1 year (screw leap years)
	genParams.MintParams.DistributionProportions = minttypes.DistributionProportions{
		Staking:          sdk.MustNewDecFromStr("0.25"), // 25%
		DeveloperRewards: sdk.MustNewDecFromStr("0.25"), // 25%
		PoolIncentives:   sdk.MustNewDecFromStr("0.45"), // 45%
		CommunityPool:    sdk.MustNewDecFromStr("0.05"), // 5%
	}
	genParams.MintParams.MintingRewardsDistributionStartEpoch = 1
	genParams.MintParams.WeightedDeveloperRewardsReceivers = []minttypes.WeightedAddress{
		{
			Address: "paro14kjcwdwcqsujkdt8n5qwpd8x8ty2rys5rjrdjj",
			Weight:  sdk.MustNewDecFromStr("0.2887"),
		},
		{
			Address: "paro1gw445ta0aqn26suz2rg3tkqfpxnq2hs224d7gq",
			Weight:  sdk.MustNewDecFromStr("0.2290"),
		},
		{
			Address: "paro13lt0hzc6u3htsk7z5rs6vuurmgg4hh2ecgxqkf",
			Weight:  sdk.MustNewDecFromStr("0.1625"),
		},
		{
			Address: "paro1kvc3he93ygc0us3ycslwlv2gdqry4ta73vk9hu",
			Weight:  sdk.MustNewDecFromStr("0.109"),
		},
		{
			Address: "paro19qgldlsk7hdv3ddtwwpvzff30pxqe9phq9evxf",
			Weight:  sdk.MustNewDecFromStr("0.0995"),
		},
		{
			Address: "paro19fs55cx4594een7qr8tglrjtt5h9jrxg458htd",
			Weight:  sdk.MustNewDecFromStr("0.06"),
		},
		{
			Address: "paro1ssp6px3fs3kwreles3ft6c07mfvj89a544yj9k",
			Weight:  sdk.MustNewDecFromStr("0.015"),
		},
		{
			Address: "paro1c5yu8498yzqte9cmfv5zcgtl07lhpjrj0skqdx",
			Weight:  sdk.MustNewDecFromStr("0.01"),
		},
		{
			Address: "paro1yhj3r9t9vw7qgeg22cehfzj7enwgklw5k5v7lj",
			Weight:  sdk.MustNewDecFromStr("0.0075"),
		},
		{
			Address: "paro18nzmtyn5vy5y45dmcdnta8askldyvehx66lqgm",
			Weight:  sdk.MustNewDecFromStr("0.007"),
		},
		{
			Address: "paro1z2x9z58cg96ujvhvu6ga07yv9edq2mvkxpgwmc",
			Weight:  sdk.MustNewDecFromStr("0.005"),
		},
		{
			Address: "paro1tvf3373skua8e6480eyy38avv8mw3hnt8jcxg9",
			Weight:  sdk.MustNewDecFromStr("0.0025"),
		},
		{
			Address: "paro1zs0txy03pv5crj2rvty8wemd3zhrka2ne8u05n",
			Weight:  sdk.MustNewDecFromStr("0.0025"),
		},
		{
			Address: "paro1djgf9p53n7m5a55hcn6gg0cm5mue4r5g3fadee",
			Weight:  sdk.MustNewDecFromStr("0.001"),
		},
		{
			Address: "paro1488zldkrn8xcjh3z40v2mexq7d088qkna8ceze",
			Weight:  sdk.MustNewDecFromStr("0.0008"),
		},
	}

	genParams.DistributionParams = distributiontypes.DefaultParams()
	genParams.DistributionParams.BaseProposerReward = sdk.MustNewDecFromStr("0.01")
	genParams.DistributionParams.BonusProposerReward = sdk.MustNewDecFromStr("0.04")
	genParams.DistributionParams.CommunityTax = sdk.MustNewDecFromStr("0")
	genParams.DistributionParams.WithdrawAddrEnabled = true

	genParams.GovParams = govtypes.DefaultParams()
	genParams.GovParams.DepositParams.MaxDepositPeriod = time.Hour * 24 * 14 // 2 weeks
	genParams.GovParams.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		sdk.NewInt(2_500_000_000),
	))
	genParams.GovParams.TallyParams.Quorum = sdk.MustNewDecFromStr("0.2") // 20%
	genParams.GovParams.VotingParams.VotingPeriod = time.Hour * 24 * 3    // 3 days

	genParams.CrisisConstantFee = sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		sdk.NewInt(500_000_000_000),
	)

	genParams.SlashingParams = slashingtypes.DefaultParams()
	genParams.SlashingParams.SignedBlocksWindow = int64(30000)                       // 30000 blocks (~41 hr at 5 second blocks)
	genParams.SlashingParams.MinSignedPerWindow = sdk.MustNewDecFromStr("0.05")      // 5% minimum liveness
	genParams.SlashingParams.DowntimeJailDuration = time.Minute                      // 1 minute jail period
	genParams.SlashingParams.SlashFractionDoubleSign = sdk.MustNewDecFromStr("0.05") // 5% double sign slashing
	genParams.SlashingParams.SlashFractionDowntime = sdk.ZeroDec()                   // 0% liveness slashing

	genParams.Epochs = epochstypes.DefaultGenesis().Epochs
	for _, epoch := range genParams.Epochs {
		epoch.StartTime = genParams.GenesisTime
	}

	genParams.IncentivesGenesis = *incentivestypes.DefaultGenesis()
	genParams.IncentivesGenesis.Params.DistrEpochIdentifier = "day"
	genParams.IncentivesGenesis.LockableDurations = []time.Duration{
		time.Hour * 24,      // 1 day
		time.Hour * 24 * 7,  // 7 day
		time.Hour * 24 * 14, // 14 days
	}

	genParams.ClaimParams = claimtypes.Params{
		AirdropStartTime:   genParams.GenesisTime,
		DurationUntilDecay: time.Hour * 24 * 60,  // 60 days = ~2 months
		DurationOfDecay:    time.Hour * 24 * 120, // 120 days = ~4 months
		ClaimDenom:         genParams.NativeCoinMetadatas[0].Base,
	}

	genParams.ConsensusParams = tmtypes.DefaultConsensusParams()
	genParams.ConsensusParams.Block.MaxBytes = 5 * 1024 * 1024
	genParams.ConsensusParams.Block.MaxGas = 6_000_000
	genParams.ConsensusParams.Evidence.MaxAgeDuration = genParams.StakingParams.UnbondingTime
	genParams.ConsensusParams.Evidence.MaxAgeNumBlocks = int64(genParams.StakingParams.UnbondingTime.Seconds()) / 3
	genParams.ConsensusParams.Version.AppVersion = 1

	return genParams
}

func TestnetGenesisParams() GenesisParams {

	genParams := MainnetGenesisParams()

	genParams.GenesisTime = time.Now()

	genParams.Epochs = append(genParams.Epochs, epochstypes.EpochInfo{
		Identifier:            "15min",
		StartTime:             time.Time{},
		Duration:              15 * time.Minute,
		CurrentEpoch:          0,
		CurrentEpochStartTime: time.Time{},
		EpochCountingStarted:  false,
	})

	for _, epoch := range genParams.Epochs {
		epoch.StartTime = genParams.GenesisTime
	}

	genParams.StakingParams.UnbondingTime = time.Hour * 24 * 7 * 2 // 2 weeks

	genParams.MintParams.EpochIdentifier = "15min"     // 15min
	genParams.MintParams.ReductionPeriodInEpochs = 192 // 2 days

	genParams.GovParams.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		sdk.NewInt(1000000), // 1 OSMO
	))
	genParams.GovParams.TallyParams.Quorum = sdk.MustNewDecFromStr("0.0000000001") // 0.00000001%
	genParams.GovParams.VotingParams.VotingPeriod = time.Second * 300              // 300 seconds

	genParams.IncentivesGenesis = *incentivestypes.DefaultGenesis()
	genParams.IncentivesGenesis.Params.DistrEpochIdentifier = "15min"
	genParams.IncentivesGenesis.LockableDurations = []time.Duration{
		time.Minute * 30, // 30 min
		time.Hour * 1,    // 1 hour
		time.Hour * 2,    // 2 hours
	}

	genParams.ClaimParams.AirdropStartTime = genParams.GenesisTime
	genParams.ClaimParams.DurationUntilDecay = time.Hour * 48 // 2 days
	genParams.ClaimParams.DurationOfDecay = time.Hour * 48    // 2 days

	genParams.PoolIncentivesGenesis.LockableDurations = genParams.IncentivesGenesis.LockableDurations

	return genParams
}