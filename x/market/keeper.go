package market

import (
	"fmt"

	app "github.com/stateset/stateset-blockchain/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/gaskv"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	log "github.com/tendermint/tendermint/libs/log"
)

// Keeper data type storing keys to the KVStore
type Keeper struct {
	storeKey   sdk.StoreKey
	codec      *codec.Codec
	paramStore params.Subspace
}

// NewKeeper creates a new keeper of the market Keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec) Keeper {
	return Keeper{
		storeKey,
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
	}
}

// NewMarket creates a new Market
func (k Keeper) NewMarket(ctx sdk.Context, id string, name string, description string, creator sdk.AccAddress) (market Market, err sdk.Error) {
	err = k.validateParams(ctx, id, name, description, creator)
	if err != nil {
		return
	}

	market = Market{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedTime: ctx.BlockHeader().Time,
	}

	k.setMarket(ctx, market)
	logger(ctx).Info(fmt.Sprintf("Created %s", market))

	return
}

// Market returns a market by its ID
func (k Keeper) Market(ctx sdk.Context, id string) (market Market, err sdk.Error) {
	store := k.store(ctx)
	marketBytes := store.Get(key(id))
	if marketBytes == nil {
		return market, ErrMarketNotFound(market.ID)
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(marketBytes, &market)

	return market, nil
}

// Markets gets all markets from the KVStore
func (k Keeper) Markets(ctx sdk.Context) (markets []Market) {
	store := k.store(ctx)
	iterator := sdk.KVStorePrefixIterator(store, MarketKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var market Market
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &market)
		markets = append(markets, market)
	}

	return
}

// AddAdmin adds a new admin
func (k Keeper) AddAdmin(ctx sdk.Context, admin, creator sdk.AccAddress) (err sdk.Error) {
	params := k.GetParams(ctx)

	// first admin can be added without any authorisation
	if len(params.MarketAdmins) > 0 && !k.isAdmin(ctx, creator) {
		err = ErrAddressNotAuthorised()
	}

	// if already present, don't add again
	for _, currentAdmin := range params.MarketAdmins {
		if currentAdmin.Equals(admin) {
			return
		}
	}

	params.MarketAdmins = append(params.MarketAdmins, admin)

	k.SetParams(ctx, params)

	return
}

// RemoveAdmin removes an admin
func (k Keeper) RemoveAdmin(ctx sdk.Context, admin, remover sdk.AccAddress) (err sdk.Error) {
	if !k.isAdmin(ctx, remover) {
		err = ErrAddressNotAuthorised()
	}

	params := k.GetParams(ctx)
	for i, currentAdmin := range params.MarketAdmins {
		if currentAdmin.Equals(admin) {
			params.MarketAdmins = append(params.MarketAdmins[:i], params.MarketAdmins[i+1:]...)
		}
	}

	k.SetParams(ctx, params)

	return
}

func (k Keeper) validateParams(ctx sdk.Context, id, name, description string, creator sdk.AccAddress) (err sdk.Error) {
	params := k.GetParams(ctx)
	if len(id) < params.MinIDLength || len(id) > params.MaxIDLength {
		err = ErrInvalidMarketMsg(
			fmt.Sprintf("ID must be between %d-%d chars in length", params.MinIDLength, params.MaxIDLength),
		)
	}
	if len(name) < params.MinNameLength || len(name) > params.MaxNameLength {
		err = ErrInvalidMarketMsg(
			fmt.Sprintf("Name must be between %d-%d chars in length", params.MinNameLength, params.MaxNameLength),
		)
	}
	if len(description) > params.MaxDescriptionLength {
		err = ErrInvalidMarketyMsg(
			fmt.Sprintf("Description must be less than %d chars in length", params.MaxDescriptionLength),
		)
	}

	if !k.isAdmin(ctx, creator) {
		err = ErrAddressNotAuthorised()
	}

	return
}

func (k Keeper) setMarket(ctx sdk.Context, market Market) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(market)
	store.Set(key(market.ID), bz)
}

func (k Keeper) isAdmin(ctx sdk.Context, address sdk.AccAddress) bool {
	for _, admin := range k.GetParams(ctx).MarketAdmins {
		if address.Equals(admin) {
			return true
		}
	}
	return false
}

func (k Keeper) store(ctx sdk.Context) sdk.KVStore {
	return gaskv.NewStore(ctx.MultiStore().GetKVStore(k.storeKey), ctx.GasMeter(), app.KVGasConfig())
}

func logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}