package factoring

import (
	app "github.com/stateset/stateset-blockchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryInvoiceLoan         = "invoice_loan"
	QueryInvoiceArguments      = "invoice_loans"
	QueryUserArguments       = "user_loans"
	QueryLoanStakes      = "loan_stakes"
	QueryMarketStakes     = "market_stakes"
	QueryStake               = "stake"
	QueryLoansByIDs      = "loans_ids"
	QueryUserStakes          = "user_stakes"
	QueryUserMarketStakes = "user_markeplaace_stakes"
	QueryInvoiceTopArgument    = "invoice_top_loan"
	QueryEarnedCoins         = "earned_coins"
	QueryTotalEarnedCoins    = "total_earned_coins"
	QueryParams              = "params"
)

type QueryInvoiceLoanParams struct {
	LoanID uint64 `json:"loan_id"`
}

type QueryInvoiceLoansParams struct {
	InvoiceID uint64 `json:"invoice_id"`
}

type QueryUserLoansParams struct {
	Address sdk.AccAddress `json:"address"`
}

type QueryLoanStakesParams struct {
	LoanID uint64 `json:"loan_id"`
}

type QueryMarketStakesParams struct {
	MarketID string `json:"market_id"`
}

type QueryStakeParams struct {
	StakeID uint64 `json:"stake_id"`
}

type QueryLoansByIDsParams struct {
	LoanIDs []uint64 `json:"loan_ids"`
}

type QueryUserStakesParams struct {
	Address sdk.AccAddress `json:"address"`
}

type QueryUserMarketStakesParams struct {
	Address     sdk.AccAddress `json:"address"`
	MarketID string         `json:"market_id"`
}

type QueryInvoiceTopLoanParams struct {
	InvoiceID uint64 `json:"invoice_id"`
}

type QueryEarnedCoinsParams struct {
	Address sdk.AccAddress `json:"address"`
}

type QueryTotalEarnedCoinsParams struct {
	Address sdk.AccAddress `json:"address"`
}

// NewQuerier creates a new querier
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryInvoiceLoan:
			return queryInvoiceLoan(ctx, req, keeper)
		case QueryInvoiceLoans:
			return queryInvoiceLoans(ctx, req, keeper)
		case QueryUserLoans:
			return queryUserLoans(ctx, req, keeper)
		case QueryLoanStakes:
			return queryLoanStakes(ctx, req, keeper)
		case QueryMarketStakes:
			return queryMarketStakes(ctx, req, keeper)
		case QueryStake:
			return queryStake(ctx, req, keeper)
		case QueryLoansByIDs:
			return queryLoanssByIDs(ctx, req, keeper)
		case QueryUserStakes:
			return queryUserStakes(ctx, req, keeper)
		case QueryUserMarkeptlaceStakes:
			return queryUserMarketStakes(ctx, req, keeper)
		case QueryInvoiceTopLoan:
			return queryInvoiceTopLoan(ctx, req, keeper)
		case QueryEarnedCoins:
			return queryEarnedCoins(ctx, req, keeper)
		case QueryTotalEarnedCoins:
			return queryTotalEarnedCoins(ctx, req, keeper)
		case QueryParams:
			return queryParams(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown factoring query endpoint")
		}
	}
}

func queryInvoiceLoan(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryInvoiceLoanParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	loan, ok := keeper.Loan(ctx, params.LoanID)
	if !ok {
		return nil, ErrCodeUnknownLoan(params.LoanID)
	}
	bz, err := keeper.codec.MarshalJSON(loan)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryUserLoans(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryUserLoansParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	loans := keeper.UserLoans(ctx, params.Address)
	bz, err := keeper.codec.MarshalJSON(loans)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryInvoiceLoans(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryInvoiceLoansParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	loans := keeper.InvoiceLoans(ctx, params.InvoiceID)
	bz, err := keeper.codec.MarshalJSON(loans)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryLoanStakes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryLoanStakesParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	stakes := keeper.LoanStakes(ctx, params.LoanID)
	bz, err := keeper.codec.MarshalJSON(stakes)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryMarketStakes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryMarketStakesParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	stakes := keeper.MarketStakes(ctx, params.MarketID)
	bz, err := keeper.codec.MarshalJSON(stakes)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryStake(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryStakeParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}

	stakes, ok := keeper.Stake(ctx, params.StakeID)
	if !ok {
		return nil, ErrCodeUnknownStake(params.StakeID)
	}

	bz, err := keeper.codec.MarshalJSON(stakes)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryLoansByIDs(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryLoansByIDsParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	var loans []Loan
	for _, id := range params.LoanIDs {
		a, ok := keeper.Loan(ctx, id)
		if !ok {
			return nil, ErrCodeUnknownLoan(id)
		}
		loans = append(loans, a)
	}

	bz, err := keeper.codec.MarshalJSON(loans)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryUserStakes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryUserStakesParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	stakes := keeper.UserStakes(ctx, params.Address)
	bz, err := keeper.codec.MarshalJSON(stakes)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryUserMarketStakes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryUserMarketStakesParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	stakes := keeper.UserMarketStakes(ctx, params.Address, params.MarktplaceID)
	bz, err := keeper.codec.MarshalJSON(stakes)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryInvoiceTopLoan(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryInvoiceTopLoanParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	loans := keeper.InvoiceLoans(ctx, params.InvoiceID)
	topLoan := Loan{}
	if len(loans) == 0 {
		bz, err := keeper.codec.MarshalJSON(topLoan)
		if err != nil {
			return nil, ErrJSONParse(err)
		}
		return bz, nil
	}
	for _, l := range loans {
		if l.IsUnhelpful {
			continue
		}
		if topLoan.ID == 0 {
			topLoan = l
		}
		if topLoan.UpvotedStake.IsLT(a.UpvotedStake) {
			topLoan = l
		}
	}
	bz, err := keeper.codec.MarshalJSON(topLoan)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryEarnedCoins(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryEarnedCoinsParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	earnedCoins := keeper.getEarnedCoins(ctx, params.Address)
	bz, err := keeper.codec.MarshalJSON(earnedCoins)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryTotalEarnedCoins(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryTotalEarnedCoinsParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	totalStakeEarned := sdk.NewCoin(app.StakeDenom, keeper.TotalEarnedCoins(ctx, params.Address))
	bz, err := keeper.codec.MarshalJSON(totalStakeEarned)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryParams(ctx sdk.Context, keeper Keeper) (result []byte, err sdk.Error) {
	params := keeper.GetParams(ctx)

	result, jsonErr := ModuleCodec.MarshalJSON(params)
	if jsonErr != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marsal result to JSON", jsonErr.Error()))
	}

	return result, nil
}