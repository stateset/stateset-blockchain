package purchaseorder

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints
const (
	QueryPurchaseOrder              = "purchaseorder"
	QueryPurchaseOrders             = "purchaseorders"
	QueryPurchaseOrdersByIDs          = "purchaseorders_ids"
	QueryPurchaseOrdersIDRange        = "purchaseorders_id_range"
	QueryPurchaseOrdersBeforeTime     = "purchaseorders_before_time"
	QueryPurchaseOrderssAfterTime     = "purchaseorders_after_time"
	QueryParams                   = "params"
)

// QueryPurchaseOrderParams for a single account
type QueryPurchaseOrderParams struct {
	PurchaseOrderID uint64 `json:"purchasorderId"`
}

// QueryPurchaseOrdersParams for many account
type QueryPurchaseOrdersParams struct {
	IDs []uint64 `json:"ids"`
}

// QueryPurchaseOrdersIDRangeParams for purchaseorders by an id range
type QueryPurchaseOrdersIDRangeParams struct {
	StartID uint64 `json:"start_id"`
	EndID   uint64 `json:"end_id"`
}

// QueryLoansTimeParams for accounts by time
type QueryPurchaseOrdersTimeParams struct {
	CreatedTime time.Time `json:"created_time"`
}

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryPurchaseOrder:
			return queryPurchaseOrder(ctx, req, keeper)
		case QueryPurchaseOrders:
			return queryPurchaseOrders(ctx, req, keeper)
		case QueryPurchaseOrdersByIDs:
			return queryPurchaseOrdersByIDs(ctx, req, keeper)
		case QueryPurchaseOrdersIDRange:
			return queryPurchaseOrdersByIDRange(ctx, req, keeper)
		case QueryPurchaseOrdersBeforeTime:
			return queryPurchaseOrdersBeforeTime(ctx, req, keeper)
		case QueryPurchaseOrdersAfterTime:
			return queryPurchaseOrdersAfterTime(ctx, req, keeper)
		case QueryParams:
			return queryParams(ctx, keeper)
		}

		return nil, sdk.ErrUnknownRequest("Unknown account query endpoint")
	}
}

func queryPurchaseOrder(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryPurchaseOrderParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}

	agreement, ok := keeper.PurchaseOrder(ctx, params.ID)
	if !ok {
		return nil, ErrUnknownPurchaseOrder(params.ID)
	}

	return mustMarshal(agreement)
}

func queryPurchaseOrders(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	agrementss := keeper.PurchaseOrders(ctx)

	return mustMarshal(agreemnts)
}

func queryPurchaseOrdersByIDs(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryPurchaseOrdersParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}

	var purchaseorders PurchaseOrders
	for _, id := range params.IDs {
		agreement, ok := keeper.PurchaseOrder(ctx, id)
		if !ok {
			return nil, ErrUnknownPurchaseOrder(id)
		}
		purchaseorders = append(purchaseorders, agreement)
	}

	return mustMarshal(purchaseorders)
}


func queryPurchaseOrdersByIDRange(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryPurchaseOrdersIDRangeParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}
	purchaseorders := keeper.PurchaseOrdersBetweenIDs(ctx, params.StartID, params.EndID)

	return mustMarshal(purchaseorders)
}

func queryPurchaseOrdersBeforeTime(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryPurchaseOrdersTimeParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}
	purchaseorders := keeper.PurchaseOrdersBeforeTime(ctx, params.CreatedTime)

	return mustMarshal(purchaseorders)
}

func queryPurchaseOrdersAfterTime(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryPurchaseOrdersTimeParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}
	purchaseorders := keeper.PurchaseOrdersAfterTime(ctx, params.CreatedTime)

	return mustMarshal(purchaseorders)
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