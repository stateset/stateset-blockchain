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

// TxCmd returns the transaction commands for this module
func TxCmd(cdc *codec.LegacyAmino) *cobra.Command {
	agreementTxCmd := &cobra.Command{
		Use:                        "agreement",
		Short:                      "agreement transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	agreementTxCmd.AddCommand(flags.PostCommands(
		CmdCreateAgreement(cdc),
		CmdUpdateAgreement(cdc),
		CmdActivateAgreement(cdc),
		CmdAmendAgreement(cdc),
		CmdRenewAgreement(cdc),
		CmdTerminateAgreement(cdc),
		CmdExpireAgreement(cdc)
	)...)

	return agreementTxCmd
}

// CmdCreateAgreeement cli command for creating atomic swaps
func CmdCreateAgreement() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-agreement [agreementNumber] [agreementName] [agreementType] [agreementStatus] [totalAgreementValue] [party] [counterparty] [AgreementStartBlock] [AgreementEndBlock]",
		Short: "Creates a new agreement",
		Args: cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsAgreementNumber := string(args[0])
			argsAgreementName := string(args[1])
			argsAgreementType := string(args[2])
			argsAgreementStatus := string(args[3])
			argsTotalAgreementValue := string(args[4])
			argsParty := string(args[5])
			argsCounterparty := string(args[6])
			argsAgreementStartBlock := string(args[7])
			argsAgreementEndBlock := string(args[8])
			
				  clientCtx, err := client.GetClientTxContext(cmd)
				  if err != nil {
					  return err
				  }
	  
				  msg := types.NewMsgCreateAgreement(clientCtx.GetFromAddress().String(), string(argsAgreementNumber), string(argsAgreementName), string(argsAgreementType), string(argsAgreementStatus), string(argsTotalAgreementValue), string(argsParty), string(argsCounterparty), string(argsAgreementStartBlock), string(argsAgreementEndBlock))
				  if err := msg.ValidateBasic(); err != nil {
					  return err
				  }
				  return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
			  },
		  }
	  
		  flags.AddTxFlagsToCmd(cmd)
	  
		  return cmd
	  }

	  func CmdUpdateAgreement() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "update-agreement [id] [agreementNumber] [agreementName] [agreementType] [agreementStatus] [totalAgreementValue] [party] [counterparty] [AgreementStartBlock] [AgreementEndBlock]",
			Short: "Update a agreement",
			Args:  cobra.ExactArgs(10),
			RunE: func(cmd *cobra.Command, args []string) error {
				id := args[0]
		  argsAgreementNumber := string(args[1])
		  argsAgreementName := string(args[2])
		  argsAgreementType := string(args[3])
		  argsAgreementStatus := string(args[4])
		  argsTotalAgreementValue := string(args[5])
		  argsParty := string(args[6])
		  argsCounterparty := string(args[7])
		  argsAgreementStartBlock := string(args[8])
		  argsAgreementEndBlock := string(args[9])
		  
				clientCtx, err := client.GetClientTxContext(cmd)
				if err != nil {
					return err
				}
	
				msg := types.NewMsgUpdateAgreement(clientCtx.GetFromAddress().String(), id, string(argsAgreementNumber), string(argsAgreementName), string(argsAgreementType), string(argsAgreementStatus), string(argsTotalAgreementValue), string(argsParty), string(argsCounterparty), string(argsAgreementStartBlock), string(argsAgreementEndBlock))
				if err := msg.ValidateBasic(); err != nil {
					return err
				}
				return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
			},
		}
	
		flags.AddTxFlagsToCmd(cmd)
	
		return cmd
	}

// GetCmdActivateAgreement cli command for activating an agreement 
func GetCmdActivateAgreement(cdc *codec.LegacyAmino) *cobra.Command {
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
func GetCmdAmendAgreement(cdc *codec.LegacyAmino) *cobra.Command {
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
func GetCmdRenewAgreement(cdc *codec.LegacyAmino) *cobra.Command {
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
func GetCmdTerminateAgreement(cdc *codec.LegacyAmino) *cobra.Command {
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
func GetCmdExpireAgreement(cdc *codec.LegacyAmino) *cobra.Command {
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