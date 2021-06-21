package cli

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	

	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/stateset/stateset-blockchain/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.LegacyAmino) *cobra.Command {
	invoiceTxCmd := &cobra.Command{
		Use:                        "invoice",
		Short:                      "invoice transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	invoiceTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreateInvoice(cdc),
		GetCmdCancelInvoice(cdc),
		GetCmdPayInvoice(cdc)
	)...)

	return agreementTxCmd
}


func GetCmdCreateInvoice(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "create [name] [number] [start] [end] [amount] [timestamp] [from] [to]",
		Short: "create a new invoice",
		Example: fmt.Sprintf("%s tx %s create state1xy7hrjy9r0algz9w3gzm8u6mrpq97kwta747gj bnb1urfermcg92dwq36572cx4xg84wpk3lfpksr5g7 bnb1uky3me9ggqypmrsvxk7ur6hqkzq7zmv4ed4ng7 now 100bnb 270 --from validator",
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

			msg := types.NewMsgCreateInvoice(
				
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}


func GetCmdCancelInvoice(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:     "cancel [invoice-id]",
		Short:   "cancel invoice by id",
		Example: fmt.Sprintf("%s tx %s cancel 6682c03cc3856879c8fb98c9733c6b0c30758299138166b6523fe94628b1d3af --from accA", version.ClientName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			invoiceID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelInvoice(from, invoiceID)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}



func GetCmdPayInvoice(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:     "pay [invoice-id]",
		Short:   "pay invoice by id",
		Example: fmt.Sprintf("%s tx %s pay 6682c03cc3856879c8fb98c9733c6b0c30758299138166b6523fe94628b1d3af --from accA", version.ClientName, types.ModuleName),
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

			msg := types.NewMsgPayInvoice(from, invoiceID)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

