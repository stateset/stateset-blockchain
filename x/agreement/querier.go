package agreement

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints
const (
	QueryAgreement                = "agreement"
	QueryAgreements               = "agreements"
	QueryAgreementsByIDs          = "agreements_ids"
	QueryLenderAgreement	      = "lender_agreements"
	QueryLenderAgreeements        = "lenders_agreements"
	QueryAgreementsIDRange        = "agreements_id_range"
	QueryAgreementsBeforeTime     = "agreements_before_time"
	QueryAgreementssAfterTime     = "agreementss_after_time"
	QueryParams                   = "params"
)

// QueryAgreementParams for a single agreement
type QueryAgreementParams struct {
	AgreementID uint64 `json:"agreementId"`
}

// QueryAgreementsParams for many agreement
type QueryAgreementsParams struct {
	IDs []uint64 `json:"ids"`
}

// QueryAgreementsIDRangeParams for agreements by an id range
type QueryAgreementsIDRangeParams struct {
	StartID uint64 `json:"start_id"`
	EndID   uint64 `json:"end_id"`
}

// QueryLoansTimeParams for agreements by time
type QueryAgreementsTimeParams struct {
	CreatedTime time.Time `json:"created_time"`
}

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryAgreement:
			return queryAgreement(ctx, req, keeper)
		case QueryAgreements:
			return queryAgreements(ctx, req, keeper)
		case QueryAgreementsByIDs:
			return queryAgreementsByIDs(ctx, req, keeper)
		case QueryAgreementsIDRange:
			return queryAgreementsByIDRange(ctx, req, keeper)
		case QueryAgreementsBeforeTime:
			return queryAgreementsBeforeTime(ctx, req, keeper)
		case QueryAgreementsAfterTime:
			return queryAgreementsAfterTime(ctx, req, keeper)
		case QueryParams:
			return queryParams(ctx, keeper)
		}

		return nil, sdk.ErrUnknownRequest("Unknown account query endpoint")
	}
}

func queryAgreement(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryAgreementParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}

	agreement, ok := keeper.Agreement(ctx, params.ID)
	if !ok {
		return nil, ErrUnknownAgreement(params.ID)
	}

	return mustMarshal(agreement)
}

func queryAgreements(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	agrementss := keeper.Agreements(ctx)

	return mustMarshal(agreemnts)
}

func queryAgreementsByIDs(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryAgreementsParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}

	var agreements Agreements
	for _, id := range params.IDs {
		agreement, ok := keeper.Agreement(ctx, id)
		if !ok {
			return nil, ErrUnknownAgreement(id)
		}
		agreements = append(agreements, agreement)
	}

	return mustMarshal(agreements)
}


func queryAgreementsByIDRange(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryAgreementsIDRangeParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}
	agreements := keeper.AgreementsBetweenIDs(ctx, params.StartID, params.EndID)

	return mustMarshal(agreements)
}

func queryAgreementsBeforeTime(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryAgreementsTimeParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}
	agreements := keeper.AgreementsBeforeTime(ctx, params.CreatedTime)

	return mustMarshal(agreements)
}

func queryAgreementsAfterTime(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryAgreementsTimeParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}
	agreements := keeper.AgreementsAfterTime(ctx, params.CreatedTime)

	return mustMarshal(agreements)
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