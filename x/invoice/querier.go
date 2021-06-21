package invoice

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints
const (
	QueryInvoice              = "invoice"
	QueryInvoices             = "invoices"
	QueryInvoiceByIDs         = "invoices_ids"
	QueryInvoicesIDRange      = "invoices_id_range"
	QueryInvoiceBeforeTime    = "invoices_before_time"
	QueryInvoicesAfterTime    = "invoices_after_time"
	QueryParams               = "params"
)

// QueryInvoiceParams for a single invoice
type QueryInvoiceParams struct {
	ID uint64 `json:"id"`
}

// QueryInvoicesParams for many invoices
type QueryInvoicesParams struct {
	IDs []uint64 `json:"ids"`
}

// QueryInvoiceIDRangeParams for invoices by an id range
type QueryInvoiceIDRangeParams struct {
	StartID uint64 `json:"start_id"`
	EndID   uint64 `json:"end_id"`
}

// QueryInvoicesTimeParams for invoices by time
type QueryInvoicesTimeParams struct {
	CreatedTime time.Time `json:"created_time"`
}

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryInvoice:
			return queryInvoice(ctx, req, keeper)
		case QueryInvoices:
			return queryInvoices(ctx, req, keeper)
		case QueryInvoicesByIDs:
			return queryInvoicesByIDs(ctx, req, keeper)
		case QueryInvoicesIDRange:
			return queryInvoicesByIDRange(ctx, req, keeper)
		case QueryInvoicesBeforeTime:
			return queryInvoicesBeforeTime(ctx, req, keeper)
		case QueryInvoicesAfterTime:
			return queryInvoicesAfterTime(ctx, req, keeper)
		case QueryParams:
			return queryParams(ctx, keeper)
		}

		return nil, sdk.ErrUnknownRequest("Unknown invoice query endpoint")
	}
}

func queryInvoice(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryInvoiceParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}

	invoice, ok := keeper.Invoice(ctx, params.ID)
	if !ok {
		return nil, ErrUnknownInvoice(params.ID)
	}

	return mustMarshal(invoice)
}

func queryInvoices(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	invoices := keeper.Invoices(ctx)

	return mustMarshal(invoices)
}

func queryInvoicesByIDs(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryInvoicesParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}

	var invoices Invoices
	for _, id := range params.IDs {
		invoice, ok := keeper.Invoice(ctx, id)
		if !ok {
			return nil, ErrUnknownInvoice(id)
		}
		invoices = append(invoices, invoice)
	}

	return mustMarshal(invoices)
}


func queryInvoicesByIDRange(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryInvoicesIDRangeParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}
	invoices := keeper.InvoicesBetweenIDs(ctx, params.StartID, params.EndID)

	return mustMarshal(invoices)
}

func queryInvoicesBeforeTime(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryInvoicesTimeParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}
	invoices := keeper.InvoicesBeforeTime(ctx, params.CreatedTime)

	return mustMarshal(invoices)
}

func queryInvoicesAfterTime(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryInvoicesTimeParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}
	invoices := keeper.InvoicesAfterTime(ctx, params.CreatedTime)

	return mustMarshal(invoices)
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