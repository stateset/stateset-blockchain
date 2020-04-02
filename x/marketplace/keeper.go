package marketplace

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

// NewKeeper creates a new keeper of the marketplace Keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec) Keeper {
	return Keeper{
		storeKey,
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
	}
}

// NewMarketplace creates a new Marketplace
func (k Keeper) NewMarketplace(ctx sdk.Context, id string, name string, description string, creator sdk.AccAddress) (marketplace Marketplace, err sdk.Error) {
	err = k.validateParams(ctx, id, name, description, creator)
	if err != nil {
		return
	}

	marketplace = Marketplace{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedTime: ctx.BlockHeader().Time,
	}

	k.setMarketplace(ctx, marketplace)
	logger(ctx).Info(fmt.Sprintf("Created %s", marketplace))

	return
}

// Marketplace returns a marketplace by its ID
func (k Keeper) Marketplace(ctx sdk.Context, id string) (marketplace Marketplace, err sdk.Error) {
	store := k.store(ctx)
	marketplaceBytes := store.Get(key(id))
	if marketplaceBytes == nil {
		return marketplace, ErrMarketplaceNotFound(marketplace.ID)
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(marketplaceBytes, &marketplace)

	return marketplace, nil
}

// Marketplaces gets all marketplaces from the KVStore
func (k Keeper) Marketplaces(ctx sdk.Context) (marketplaces []Marketplace) {
	store := k.store(ctx)
	iterator := sdk.KVStorePrefixIterator(store, MarketplaceKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var marketplace Marktplace
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &marketplace)
		marketplaces = append(marketplaces, marketplace)
	}

	return
}

// AddAdmin adds a new admin
func (k Keeper) AddAdmin(ctx sdk.Context, admin, creator sdk.AccAddress) (err sdk.Error) {
	params := k.GetParams(ctx)

	// first admin can be added without any authorisation
	if len(params.MarketplaceAdmins) > 0 && !k.isAdmin(ctx, creator) {
		err = ErrAddressNotAuthorised()
	}

	// if already present, don't add again
	for _, currentAdmin := range params.MarketplaceAdmins {
		if currentAdmin.Equals(admin) {
			return
		}
	}

	params.MarketplaceAdmins = append(params.MarketplaceAdmins, admin)

	k.SetParams(ctx, params)

	return
}

// RemoveAdmin removes an admin
func (k Keeper) RemoveAdmin(ctx sdk.Context, admin, remover sdk.AccAddress) (err sdk.Error) {
	if !k.isAdmin(ctx, remover) {
		err = ErrAddressNotAuthorised()
	}

	params := k.GetParams(ctx)
	for i, currentAdmin := range params.MarketplaceAdmins {
		if currentAdmin.Equals(admin) {
			params.MarketplaceAdmins = append(params.MarketplaceAdmins[:i], params.MarketplaceAdmins[i+1:]...)
		}
	}

	k.SetParams(ctx, params)

	return
}

func (k Keeper) validateParams(ctx sdk.Context, id, name, description string, creator sdk.AccAddress) (err sdk.Error) {
	params := k.GetParams(ctx)
	if len(id) < params.MinIDLength || len(id) > params.MaxIDLength {
		err = ErrInvalidMarketplaceMsg(
			fmt.Sprintf("ID must be between %d-%d chars in length", params.MinIDLength, params.MaxIDLength),
		)
	}
	if len(name) < params.MinNameLength || len(name) > params.MaxNameLength {
		err = ErrInvalidMarketplaceMsg(
			fmt.Sprintf("Name must be between %d-%d chars in length", params.MinNameLength, params.MaxNameLength),
		)
	}
	if len(description) > params.MaxDescriptionLength {
		err = ErrInvalidMarketplaceyMsg(
			fmt.Sprintf("Description must be less than %d chars in length", params.MaxDescriptionLength),
		)
	}

	if !k.isAdmin(ctx, creator) {
		err = ErrAddressNotAuthorised()
	}

	return
}

func (k Keeper) setMarketplace(ctx sdk.Context, marketplace Marketplace) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(marketplace)
	store.Set(key(marketplace.ID), bz)
}

func (k Keeper) isAdmin(ctx sdk.Context, address sdk.AccAddress) bool {
	for _, admin := range k.GetParams(ctx).MarketplaceAdmins {
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