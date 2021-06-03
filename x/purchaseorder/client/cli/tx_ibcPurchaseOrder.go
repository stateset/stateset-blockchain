package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	channelutils "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/client/utils"
	"github.com/domsteil/stateset/x/purchaseorder/types"
)

var _ = strconv.Itoa(0)

func CmdSendIbcPurchaseOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-ibcPurchaseOrder [src-port] [src-channel] [purchaseordername] [purchaseordernumber] [status] [total]",
		Short: "Send a ibcPurchaseOrder over IBC",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsPurchaseordername := string(args[2])
			argsPurchaseordernumber := string(args[3])
			argsStatus := string(args[4])
			argsTotal := string(args[5])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress().String()
			srcPort := args[0]
			srcChannel := args[1]

			// Get the relative timeout timestamp
			timeoutTimestamp, err := cmd.Flags().GetUint64(flagPacketTimeoutTimestamp)
			if err != nil {
				return err
			}
			consensusState, _, _, err := channelutils.QueryLatestConsensusState(clientCtx, srcPort, srcChannel)
			if err != nil {
				return err
			}
			if timeoutTimestamp != 0 {
				timeoutTimestamp = consensusState.GetTimestamp() + timeoutTimestamp
			}

			msg := types.NewMsgSendIbcPurchaseOrder(sender, srcPort, srcChannel, timeoutTimestamp, string(argsPurchaseordername), string(argsPurchaseordernumber), string(argsStatus), string(argsTotal))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Uint64(flagPacketTimeoutTimestamp, DefaultRelativePacketTimeoutTimestamp, "Packet timeout timestamp in nanoseconds. Default is 10 minutes.")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}