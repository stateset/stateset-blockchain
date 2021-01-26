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

	"github.com/stateset/stateset-blockchain/x/invoice"

)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	invoiceQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the invoice module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	invoiceQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryInvoice(),
		GetCmdQueryInvoices(),
		GetCmdQueryInvoiceLineItem(),
		GetCmdQueryInvoicesLineItems(),
		GetCmdQueryPoolBatchSwap(),
	)

	return invoiceQueryCmd
}

//GetCmdQueryParams implements the params query command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current invoice parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as invoice parameters.
Example:
$ %s query invoice params
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

func GetCmdQueryInvoicePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoice [invoice-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a invoice",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of an invoice
Example:
$ %s query invoice 1
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
				return fmt.Errorf("invoice-id %s not a valid uint, please input a valid invoice-id", args[0])
			}

			// Query the pool
			res, err := queryClient.Invoice(
				context.Background(),
				&types.QueryInvoiceRequest{InvoiceId: invoiceId},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch invoiceId %d: %s", invoiceId, err)
			}

			params := &types.QueryInvoiceRequest{InvoiceId: invoiceId}
			res, err = queryClient.Invoice(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryInvoices() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoices",
		Args:  cobra.NoArgs,
		Short: "Query for all invoices",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all invoices on a network.
Example:
$ %s query invoices
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
			result, err := queryClient.Invoices(context.Background(), &types.QueryAgreementsRequest{Pagination: pageReq})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(result)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}