package cli

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/stateset/stateset-blockchain/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	agreementTxCmd := &cobra.Command{
		Use:                        "agreement",
		Short:                      "agreement transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	agreementTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreateAgreement(cdc),
		GetCmdActivateAgreement(cdc),
		GetCmdAmendAgreement(cdc),
		GetCmdRenewAgreement(cdc),
		GetCmdTerminateAgreement(cdc),
		GetCmdExpireAgreement(cdc)
	)...)

	return agreementTxCmd
}

// GetCmdCreateAgreeement cli command for creating atomic swaps
func GetCmdCreateAgreement(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create [name] [number] [start] [end] [amount] [timestamp] [from] [to]",
		Short: "create a new agreement",
		Example: fmt.Sprintf("%s tx %s create kava1xy7hrjy9r0algz9w3gzm8u6mrpq97kwta747gj bnb1urfermcg92dwq36572cx4xg84wpk3lfpksr5g7 bnb1uky3me9ggqypmrsvxk7ur6hqkzq7zmv4ed4ng7 now 100bnb 270 --from validator",
			version.ClientName, types.ModuleName),
		Args: cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()
			to, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// Timestamp defaults to time.Now() unless it's explicitly set
			var timestamp int64
			if strings.Compare(args[3], "now") == 0 {
				timestamp = tmtime.Now().Unix()
			} else {
				timestamp, err = strconv.ParseInt(args[3], 10, 64)
				if err != nil {
					return err
				}
			}

			// Generate cryptographically strong pseudo-random number
			randomNumber, err := types.GenerateSecureRandomNumber()
			if err != nil {
				return err
			}

			randomNumberHash := types.CalculateRandomHash(randomNumber, timestamp)

			// Print random number, timestamp, and hash to user's console
			fmt.Printf("\nRandom number: %s\n", hex.EncodeToString(randomNumber))
			fmt.Printf("Timestamp: %d\n", timestamp)
			fmt.Printf("Random number hash: %s\n\n", hex.EncodeToString(randomNumberHash))

			coins, err := sdk.ParseCoins(args[4])
			if err != nil {
				return err
			}

			heightSpan, err := strconv.ParseUint(args[5], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateAgreement(
				agreementID, agreementNumber, agreementName, agreementType, agreementStatus, totalAgreementValue, from, to, startDate, endDate, paid, active, timestamp
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdActivateAgreement cli command for activating an agreement 
func GetCmdActivateAgreement(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "activate [agreement-id]",
		Short:   "activate agreement by id",
		Example: fmt.Sprintf("%s tx %s activate 6682c03cc3856879c8fb98c9733c6b0c30758299138166b6523fe94628b1d3af --from accA", version.ClientName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			agreementID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgActivateAgreement(from, agreementID)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdAmendAgreement cli command for activating an agreement 
func GetCmdAmendAgreement(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "amend [agreement-id]",
		Short:   "amend agreement by id",
		Example: fmt.Sprintf("%s tx %s amend 6682c03cc3856879c8fb98c9733c6b0c30758299138166b6523fe94628b1d3af --from accA", version.ClientName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			agreementID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgAmendAgreement(from, agreementID)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdRenewAgreement cli command for activating an agreement 
func GetCmdRenewAgreement(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "renew [agreement-id]",
		Short:   "renew agreement by id",
		Example: fmt.Sprintf("%s tx %s amend 6682c03cc3856879c8fb98c9733c6b0c30758299138166b6523fe94628b1d3af --from accA", version.ClientName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			agreementID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgRenewAgreement(from, agreementID)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdTeriminateAgreement cli command for activating an agreement 
func GetCmdTerminateAgreement(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "terminate [agreement-id]",
		Short:   "terminate agreement by id",
		Example: fmt.Sprintf("%s tx %s expire 6682c03cc3856879c8fb98c9733c6b0c30758299138166b6523fe94628b1d3af --from accA", version.ClientName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			agreementID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgTeriminateAgreement(from, agreementID)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdExpireAgreement cli command for activating an agreement 
func GetCmdExpireAgreement(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "expire [agreement-id]",
		Short:   "expire agreement by id",
		Example: fmt.Sprintf("%s tx %s expire 6682c03cc3856879c8fb98c9733c6b0c30758299138166b6523fe94628b1d3af --from accA", version.ClientName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			agreementID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgExpireAgreement(from, agreementID)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}