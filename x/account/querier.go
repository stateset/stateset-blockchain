package account

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints
const (
	QueryAccount             = "account"
	QueryAccounts            = "accounts"
	QueryAccountsByIDs       = "accounts_ids"
	QueryAccountsIDRange     = "accounts_id_range"
	QueryAccountsBeforeTime  = "accounts_before_time"
	QueryAccountsAfterTime   = "accounts_after_time"
	QueryParams              = "params"
)

// QueryAccountParams for a single account
type QueryAccountParams struct {
	ID uint64 `json:"id"`
}

// QueryAccountsParams for many account
type QueryAccountsParams struct {
	IDs []uint64 `json:"ids"`
}

// QueryAccountsIDRangeParams for accounts by an id range
type QueryAccountsIDRangeParams struct {
	StartID uint64 `json:"start_id"`
	EndID   uint64 `json:"end_id"`
}

// QueryAccountsTimeParams for accounts by time
type QueryAccountsTimeParams struct {
	CreatedTime time.Time `json:"created_time"`
}

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryAccount:
			return queryAccount(ctx, req, keeper)
		case QueryAccounts:
			return queryAccounts(ctx, req, keeper)
		case QueryAccountsByIDs:
			return queryAccountsByIDs(ctx, req, keeper)
		case QueryAccountsIDRange:
			return queryAccountsByIDRange(ctx, req, keeper)
		case QueryAccountsBeforeTime:
			return queryAccountsBeforeTime(ctx, req, keeper)
		case QueryAccountsAfterTime:
			return queryAccountsAfterTime(ctx, req, keeper)
		case QueryParams:
			return queryParams(ctx, keeper)
		}

		return nil, sdk.ErrUnknownRequest("Unknown account query endpoint")
	}
}

func queryAccount(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryAccountParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}

	account, ok := keeper.Account(ctx, params.ID)
	if !ok {
		return nil, ErrUnknownClaim(params.ID)
	}

	return mustMarshal(claim)
}

func queryAccounts(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	accounts := keeper.Accounts(ctx)

	return mustMarshal(accounts)
}

func queryAccountsByIDs(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryAccountsParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}

	var accounts Accounts
	for _, id := range params.IDs {
		account, ok := keeper.Account(ctx, id)
		if !ok {
			return nil, ErrUnknownAccount(id)
		}
		accounts = append(accounts, claim)
	}

	return mustMarshal(accounts)
}


func queryAccountssByIDRange(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryAccountsIDRangeParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}
	accounts := keeper.AccountsBetweenIDs(ctx, params.StartID, params.EndID)

	return mustMarshal(accounts)
}

func queryAccountsBeforeTime(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryAccountsTimeParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}
	accounts := keeper.AccountsBeforeTime(ctx, params.CreatedTime)

	return mustMarshal(accounts)
}

func queryAccountsAfterTime(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryAccountsTimeParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}
	accounts := keeper.AccountsAfterTime(ctx, params.CreatedTime)

	return mustMarshal(accounts)
}

func queryParams(ctx sdk.Context, keeper Keeper) (result []byte, err sdk.Error) {
	params := keeper.GetParams(ctx)

	result, jsonErr := ModuleCodec.MarshalJSON(params)
	if jsonErr != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marsal result to JSON", jsonErr.Error()))
	}

	return result, nil
}

func mustMarshal(v interface{}) (result []byte, err sdk.Error) {
	result, jsonErr := codec.MarshalJSONIndent(ModuleCodec, v)
	if jsonErr != nil {
		return nil, ErrJSONParse(jsonErr)
	}

	return
}