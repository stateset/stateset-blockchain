package market

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the stateset Querier
const (
	QueryMarket   = "market"
	QueryMarkets = "markets"
	QueryParams      = "params"
)

// QueryMarketParams are params for querying markets by id queries
type QueryMarketParams struct {
	ID string
}

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, request abci.RequestQuery) (result []byte, err sdk.Error) {
		switch path[0] {
		case QueryMarket:
			return queryMarket(ctx, request, k)
		case QueryMarkets:
			return queryMarkets(ctx, k)
		case QueryParams:
			return queryParams(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown stateset query endpoint: commmunity/%s", path[0]))
		}
	}
}

func queryMarket(ctx sdk.Context, req abci.RequestQuery, k Keeper) (result []byte, err sdk.Error) {
	var params QueryMarketParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}

	market, err := k.Market(ctx, params.ID)
	if err != nil {
		return
	}

	return mustMarshal(market)
}

func queryMarkets(ctx sdk.Context, k Keeper) (result []byte, err sdk.Error) {
	markets := k.Markets(ctx)
	return mustMarshal(markets)
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