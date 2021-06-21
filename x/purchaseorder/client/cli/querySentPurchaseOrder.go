package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/stateset/stateset-blockchain/x/purchaseorder/types"
)

func CmdListSentPurchaseOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-sentPurchaseOrder",
		Short: "list all sentPurchaseOrder",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllSentPurchaseOrderRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.SentPurchaseOrderAll(context.Background(), params)
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

func CmdShowSentPurchaseOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-sentPurchaseOrder [id]",
		Short: "shows a sentPurchaseOrder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			params := &types.QueryGetSentPurchaseOrderRequest{
				Id: id,
			}

			res, err := queryClient.SentPurchaseOrder(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
