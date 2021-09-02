package cli


import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/stateset/stateset-blockchain/x/purchaseorder"

)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	agreementQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the purchase order module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	agreementQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryPurchaseOrder(),
		GetCmdQueryPurchaseOrders(),
		GetCmdQueryPurchaseOrderLineItem(),
		GetCmdQueryPurchaseOrderLineItems(),
		GetCmdQueryPoolBatchSwap(),
	)

	return agreementQueryCmd
}

//GetCmdQueryParams implements the params query command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current purchase order parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as purchase order parameters.
Example:
$ stateset query purchaseorder params
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryPurchaseOrderPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purchase order [purchaseorder-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a purchase order",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of an purchase order
Example:
$ stateset query purchaseorder 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("purchaseorder-id %s not a valid uint, please input a valid purchaseorder-id", args[0])
			}

			// Query the pool
			res, err := queryClient.PurchaseOrder(
				context.Background(),
				&types.QueryPurchaseOrderRequest{PurchaseOrderId: purchaseorderId},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch purchaseorderId %d: %s", purchaseorderId, err)
			}

			params := &types.QueryPurchaseOrderRequest{PurchaseOrderId: purchaseorderId}
			res, err = queryClient.PurchaseOrder(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryPurchaseOrders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purchasorders",
		Args:  cobra.NoArgs,
		Short: "Query for all purchase orders",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all purchase orders on a network.
Example:
$ stateset query purchaseorders
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			result, err := queryClient.PurchaseOrders(context.Background(), &types.QueryPurchaseOrdersRequest{Pagination: pageReq})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(result)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}