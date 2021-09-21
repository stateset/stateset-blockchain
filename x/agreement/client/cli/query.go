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

)

// QueryCmd returns the cli query commands for this module
func QueryCmd() *cobra.Command {
	agreementQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the agreement module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	agreementQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryAgreement(),
		GetCmdQueryAgreements(),
		GetCmdQueryAgreementLineItem(),
		GetCmdQueryAgreementLineItems(),
	)

	return agreementQueryCmd
}

//GetCmdQueryParams implements the params query command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current agreement parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as agreement parameters.
Example:
$ %s query agreement params
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

func GetCmdQueryAgreementPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agreement [agreement-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a agreement",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of an agreement
Example:
$ %s query agreement 1
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
				return fmt.Errorf("agreement-id %s not a valid uint, please input a valid agreement-id", args[0])
			}

			// Query the pool
			res, err := queryClient.Agreement(
				context.Background(),
				&types.QueryAgreementRequest{AgreementId: agreementId},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch agreementId %d: %s", agreementId, err)
			}

			params := &types.QueryAgreementRequest{AgreementId: agreementId}
			res, err = queryClient.Agreement(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryAgreements() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agreements",
		Args:  cobra.NoArgs,
		Short: "Query for all agreements",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all agreements on a network.
Example:
$ %s query agreements
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
			result, err := queryClient.Agreements(context.Background(), &types.QueryAgreementsRequest{Pagination: pageReq})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(result)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}