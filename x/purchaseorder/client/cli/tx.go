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
	purchaseorderTxCmd := &cobra.Command{
		Use:                        "purchaseorder",
		Short:                      "Purchase Order transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	purchaseorderTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreatePurchaseOrder(cdc),
		getCmdCompletePurchaseOrder(cdc),
		GetCmdCancelPurchaseOrder(cdc)
		GetCmdLockPurchaseOrder(cdc)
		GetCmdFinancePurchaseOrder(cdc)
	)...)

	return purchaseorderTxCmd
}

// GetCmdCreatePurchaseOrder cli command for creating atomic swaps
func GetCmdCreatePurchaseOrder(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "create [name] [number] [start] [end] [amount] [timestamp] [from] [to]",
		Short: "create a new purchase order",
		Example: fmt.Sprintf("%s tx %s --from validator",
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

			msg := types.NewMsgCreatePurchaseOrder(
				AgreementID, PurchaseOrderID, Body, Lender, Source
			
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}


func GetCmdCancelPurchaseOrder(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:     "cancel [purchaseorder-id]",
		Short:   "cancel purchaseorder by id",
		Example: fmt.Sprintf("%s tx %s cancel 6682c03cc3856879c8fb98c9733c6b0c30758299138166b6523fe94628b1d3af --from accA", version.ClientName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			purchaseorderID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelPurchaseOrder(from, purchaseorderID)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdCompletePurchaseOrder cli command for completeing a purchase order
func GetCmdCompletePurchaseOrder(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:     "complete [purchaseorder-id]",
		Short:   "complete purchaseorder by id",
		Example: fmt.Sprintf("%s tx %s complete 6682c03cc3856879c8fb98c9733c6b0c30758299138166b6523fe94628b1d3af --from accA", version.ClientName, types.ModuleName),
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

			msg := types.NewMsgCompletePurchaseOrder(from, purchaseOrderID)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdLockPurchaseOrder cli command for locking a purchase order
func GetCmdLockPurchaseOrder(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:     "lock [purchaseorder-id]",
		Short:   "lock purchaseorder by id",
		Example: fmt.Sprintf("%s tx %s complete 6682c03cc3856879c8fb98c9733c6b0c30758299138166b6523fe94628b1d3af --from accA", version.ClientName, types.ModuleName),
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

			msg := types.NewMsgLockPurchaseOrder(from, purchaseOrderID)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdFinancePurchaseOrder cli command for financing a purchase order
func GetCmdFinancePurchaseOrder(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:     "finance [purchaseorder-id]",
		Short:   "finance purchaseorder by id",
		Example: fmt.Sprintf("%s tx %s complete 6682c03cc3856879c8fb98c9733c6b0c30758299138166b6523fe94628b1d3af --from accA", version.ClientName, types.ModuleName),
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

			msg := types.NewMsgFinancePurchaseOrder(from, purchaseOrderID)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

flags.AddTxFlagsToCmd(cmd)

return cmd
}