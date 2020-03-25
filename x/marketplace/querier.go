package markerplace

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the stateset Querier
const (
	QueryMarketplace   = "marketplace"
	QueryMarketplaces = "marketplaces"
	QueryParams      = "params"
)

// QueryMarketplaceParams are params for querying marketplaces by id queries
type QueryMarketplaceParams struct {
	ID string
}

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, request abci.RequestQuery) (result []byte, err sdk.Error) {
		switch path[0] {
		case QueryMarketplace:
			return queryMarketplace(ctx, request, k)
		case QueryMarketplaces:
			return queryMarketplaces(ctx, k)
		case QueryParams:
			return queryParams(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown stateset query endpoint: commmunity/%s", path[0]))
		}
	}
}

func queryMarketplace(ctx sdk.Context, req abci.RequestQuery, k Keeper) (result []byte, err sdk.Error) {
	var params QueryMarketplaceParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}

	marketplace, err := k.Marketplace(ctx, params.ID)
	if err != nil {
		return
	}

	return mustMarshal(marketplace)
}

func queryMarketplaces(ctx sdk.Context, k Keeper) (result []byte, err sdk.Error) {
	marketplaces := k.Marketplaces(ctx)
	return mustMarshal(marketplaces)
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