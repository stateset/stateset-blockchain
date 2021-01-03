package purchaseorder

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keys for params
var (
	KeyMinPurchaseOrderLength = []byte("minPurchaseOrderLength")
	KeyMaxPurchaseOrderLength = []byte("maxPurchaseOrderLength")
	KeyPurchaseOrderAdmins    = []byte("purchaseorderAdmins")
)

// Params holds parameters for a PurchaseOrder
type Params struct {
	MinPurchaseOrderLength int              `json:"min_purchaseorder_length"`
	MaxPurchaseOrderLength int              `json:"max_purchaseorder_length"`
	PurchaseOrderAdmins    []sdk.AccAddress `json:"purchaseorder_admins"`
}

// DefaultParams is the PurchaseOrder params for testing
func DefaultParams() Params {
	return Params{
		MinPurchaseOrderLength: 25,
		MaxPurchaseOrderLength: 140,
		PurchaseOrderAdmins:    []sdk.AccAddress{},
	}
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMinPurchaseOrderLength, Value: &p.MinPurchaseOrderLength},
		{Key: KeyMaxPurchaseOrderLength, Value: &p.MaxPurchaseOrderLength},
		{Key: KeyPurchaseOrderAdmins, Value: &p.PurchaseOrderAdmins},
	}
}

// ParamKeyTable for purchaseorder module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// GetParams gets the genesis params for the purchaseorder
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var paramSet Params
	k.paramStore.GetParamSet(ctx, &paramSet)
	return paramSet
}

// SetParams sets the params for the purchaseorder
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramStore.SetParamSet(ctx, &params)
	logger(ctx).Info(fmt.Sprintf("Loaded purchaseorder params: %+v", params))
}

// UpdateParams updates the required params
func (k Keeper) UpdateParams(ctx sdk.Context, updater sdk.AccAddress, updates Params, updatedFields []string) sdk.Error {
	if !k.isAdmin(ctx, updater) {
		err := ErrAddressNotAuthorised()
		return err
	}

	current := k.GetParams(ctx)
	updated := k.getUpdatedParams(current, updates, updatedFields)
	k.SetParams(ctx, updated)

	return nil
}

func (k Keeper) getUpdatedParams(current Params, updates Params, updatedFields []string) Params {
	updated := current
	mapParams(updates, func(param string, index int, field reflect.StructField) {
		if isIn(param, updatedFields) {
			reflect.ValueOf(&updated).Elem().FieldByName(field.Name).Set(
				reflect.ValueOf(
					reflect.ValueOf(updates).FieldByName(field.Name).Interface(),
				),
			)
		}
	})

	return updated
}

func isIn(needle string, haystack []string) bool {
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}

	return false
}

// mapParams walks over each param, and ignores the *_admins param because they are out of scope for this CLI command
func mapParams(params interface{}, fn func(param string, index int, field reflect.StructField)) {
	rParams := reflect.TypeOf(params)
	for i := 0; i < rParams.NumField(); i++ {
		field := rParams.Field(i)
		param := field.Tag.Get("json")
		fn(param, i, field)
	}
}