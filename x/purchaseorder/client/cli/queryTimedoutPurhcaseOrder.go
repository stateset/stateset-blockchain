package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/stateset/stateset-blockchain/x/purchaseorder/types"
)

func CmdListTimedoutPurchaseOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-timedoutPurchaseOrder",
		Short: "list all timedoutPurchaseOrder",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllTimedoutPurchaseOrderRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.TimedoutPurchaseOrderAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowTimedoutPurchaseOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-timedoutPurchaseOrder [id]",
		Short: "shows a timedoutPurchaseOrder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			params := &types.QueryGetTimedoutPurchaseOrderRequest{
				Id: id,
			}

			res, err := queryClient.TimedoutPurchaseOrder(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
