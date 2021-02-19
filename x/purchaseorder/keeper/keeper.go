package keeper

import (
	"net/url"
	"time"

	app "github.com/stateset/stateset-blockchain/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/gaskv"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	log "github.com/tendermint/tendermint/libs/log"
)

// Keeper is the model object for the module
type Keeper struct {
	storeKey   sdk.StoreKey
	codec      *codec.Codec
	paramStore params.Subspace
	accountKeeper   AccountKeeper
	invoiceKeeper InvoiceKeeper
	liquidityKeeper LiquidityKeeper
}

// NewKeeper creates a new account keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec, accountKeeper AccountKeeper, marketKeeper market.Keeper) Keeper {
	return Keeper{
		storeKey,
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
		accountKeeper,
		marketKeeper,
	}
}

// CreatePurchaseOrder creates a new purchaseorder in the purchaseorder key-value store
func (k Keeper) CreatePurchaseOrder(ctx sdk.Context, body, purchaseorderID string,
	merchant sdk.AccAddress, source url.URL) (purchaseorder PurchaseOrder, err sdk.Error) {

	err = k.validateLength(ctx, body)
	if err != nil {
		return
	}

	purchaseorder = NewPurchaseOrder(purchaseorderID, MarketID, body, merchant, source,
		ctx.BlockHeader().Time,
	)

	// persist purchaseorder
	k.setPurchaseOrder(ctx, purchaseorder)
	// increment purchaseorderID (primary key) for next purchaseorder
	k.setPurchaseOrderID(ctx, purchaseorderID+1)

	// persist associations
	k.setPurchaserPurchaseOrder(ctx, purchaseorder.Purchaser, purchaseorderID)
	k.setVendorPurchaseOrder(ctx, purchaseorder.Vendor, purchaseorderId)
	k.setCreatedTimePurchaseOrder(ctx, purchaseorder.CreatedTime, purchaseorderID)

	logger(ctx).Info("Submitted " + purchaseorder.String())

	return purchaseorder, nil
}

// EditPurchaseOrder allows admins to edit the body of an purchaseorder

func (k Keeper) EditPurchaseOrder(ctx sdk.Context, id uint64, body string, editor sdk.AccAddress) (purchaseorder PurchaseOrder, err sdk.Error) {
	if !k.isAdmin(ctx, editor) {
		err = ErrAddressNotAuthorised()
		return
	}

	err = k.validateLength(ctx, body)
	if err != nil {
		return
	}

	purchaseorder, ok := k.PurchaseOrder(ctx, id)
	if !ok {
		err = ErrUnknownPurchaseOrder(id)
		return
	}

	purchaseorder.Body = body
	k.setPurchaseOrder(ctx, purchaseorder)

	return
}

// PurchaseOrder gets a single purchaseorder by its ID
func (k Keeper) PurchaseOrder(ctx sdk.Context, id uint64) (purchaseorder PurchaseOrder, ok bool) {
	store := k.store(ctx)
	purchaseorderBytes := store.Get(key(id))
	if purchaseorderBytes == nil {
		return purchaseorder, false
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(purchaseorderBytes, &purchaseorder)

	return purchaseorder, true
}

// PurchaseOrders gets all the purchaseorders in reverse order
func (k Keeper) PurchaseOrders(ctx sdk.Context) (purchaseorders PurchaseOrders) {
	store := k.store(ctx)
	iterator := sdk.KVStoreReversePrefixIterator(store, PurchaseOrdersKeyPrefix)

	return k.iterate(iterator)
}

// PurchaseOrdersBetweenIDs gets all purchaseorders between startPurchaseOrderID to endPurchaseOrderID
func (k Keeper) PurchaseOrdersBetweenIDs(ctx sdk.Context, startPurchaseOrderID, endPurchaseOrderID uint64) (purchaseorders PurchaseOrders) {
	iterator := k.purchaseordersIterator(ctx, startPurchaseOrderID, endPurchaseOrderID)

	return k.iterate(iterator)
}

// PurchaseOrdersBetweenTimes gets all purchaseorders between startTime and endTime
func (k Keeper) PurchaseOrdersBetweenTimes(ctx sdk.Context, startTime time.Time, endTime time.Time) (purchaseorders PurchaseOrders) {
	iterator := k.createdTimeRangePurchaseOrdersIterator(ctx, startTime, endTime)

	return k.iterateAssociated(ctx, iterator)
}

// PurchaseOrderssBeforeTime gets all purchaseorders after a certain CreatedTime
func (k Keeper) PurchaseOrdersBeforeTime(ctx sdk.Context, createdTime time.Time) (purchaseorders PurchaseOrders) {
	iterator := k.beforeCreatedTimePurchaseOrdersIterator(ctx, createdTime)

	return k.iterateAssociated(ctx, iterator)
}

// PurchaseOrderssAfterTime gets all purchaseorders after a certain CreatedTime
func (k Keeper) PurchaseOrderssAfterTime(ctx sdk.Context, createdTime time.Time) (purchaseorders PurchaseOrders) {
	iterator := k.afterCreatedTimePurchaseOrdersIterator(ctx, createdTime)

	return k.iterateAssociated(ctx, iterator)
}

// MerchantPurchaseOrders gets all the purchaseorders for a given merchant
func (k Keeper) MerchantPurchaseOrders(ctx sdk.Context, creator sdk.AccAddress) (purchaseorders PurchaseOrders) {
	return k.associatedPurchaseOrders(ctx, merchantPurchaseOrdersKey(merchant))
}

// AddBackingStake adds a stake amount to the total backing amount
func (k Keeper) AddFinancingStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	purchaseorder, ok := k.PurchaseOrder(ctx, id)
	if !ok {
		return ErrUnknownPurchaseOrder(id)
	}
	purchaseorder.TotalBacked = purchaseorder.TotalBacked.Add(stake)
	purchaseorder.TotalStakers++
	k.setPurchaseOrder(ctx, purchaseorder)

	return nil
}

// AddFactorStake adds a stake amount to the total factor amount
func (k Keeper) AddFinanceStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	, ok := k.PurchaseOrder(ctx, id)
	if !ok {
		return ErrUnknownPurchaseOrder(id)
	}
	purchaseorder.TotalFinancer = purchaseorder.TotalFinanced.Add(stake)
	purchaseorder.TotalStakers++
	k.setPurchaseOrder(ctx, purchaseorder)

	return nil
}

// SubtractFactorStake adds a stake amount to the total factoring amount
func (k Keeper) SubtractFinanceStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	purchaseorder, ok := k.PurchaseOrder(ctx, id)
	if !ok {
		return ErrUnknownPurchaseOrder(id)
	}
	purchaseorder.TotalFactored = purchaseorder.TotalFactored.Sub(stake)
	k.setIncoice(ctx, purchaseorder)

	return nil
}

// SetFirstArgumentTime sets time when first argument was created on a purchaseorder
func (k Keeper) SetFirstArgumentTime(ctx sdk.Context, id uint64, firstArgumentTime time.Time) sdk.Error {
	purchaseorder, ok := k.PurchaseOrder(ctx, id)
	if !ok {
		return ErrUnknownPurchaseOrder(id)
	}
	purchaseorder.FirstArgumentTime = firstArgumentTime
	k.setPurchaseOrder(ctx, purchaseorder)

	return nil
}

// AddAdmin adds a new admin
func (k Keeper) AddAdmin(ctx sdk.Context, admin, creator sdk.AccAddress) (err sdk.Error) {
	params := k.GetParams(ctx)

	// first admin can be added without any authorisation
	if len(params.PurchaseOrderAdmins) > 0 && !k.isAdmin(ctx, creator) {
		err = ErrAddressNotAuthorised()
	}

	// if already present, don't add again
	for _, currentAdmin := range params.PurchaseOrderAdmins {
		if currentAdmin.Equals(admin) {
			return
		}
	}

	params.PurchaseOrderAdmins = append(params.PurchaseOrderAdmins, admin)

	k.SetParams(ctx, params)

	return
}

// RemoveAdmin removes an admin
func (k Keeper) RemoveAdmin(ctx sdk.Context, admin, remover sdk.AccAddress) (err sdk.Error) {
	if !k.isAdmin(ctx, remover) {
		err = ErrAddressNotAuthorised()
	}

	params := k.GetParams(ctx)
	for i, currentAdmin := range params.PurchaseOrderAdmins {
		if currentAdmin.Equals(admin) {
			params.PurchaseOrderAdmins = append(params.PurchaseOrderAdmins[:i], params.PurchaseOrderAdmins[i+1:]...)
		}
	}

	k.SetParams(ctx, params)

	return
}

func (k Keeper) isAdmin(ctx sdk.Context, address sdk.AccAddress) bool {
	for _, admin := range k.GetParams(ctx).PurchaseOrderAdmins {
		if address.Equals(admin) {
			return true
		}
	}
	return false
}

func (k Keeper) validateLength(ctx sdk.Context, body string) sdk.Error {
	var minPurchaseOrderLength int
	var maxPurchaseOrderLength int

	k.paramStore.Get(ctx, KeyMinPurchaseOrderLength, &minPurchaseOrderLength)
	k.paramStore.Get(ctx, KeyMaxPurchaseOrderLength, &maxPurchaseOrderLength)

	len := len([]rune(body))
	if len < minPurchaseOrderLength {
		return ErrInvalidBodyTooShort(body)
	}
	if len > maxPurchaseOrderLength {
		return ErrInvalidBodyTooLong()
	}

	return nil
}

// purchaseorderID gets the highest purchaseorder ID
func (k Keeper) purchaseorderID(ctx sdk.Context) (purchaseorderID uint64, err sdk.Error) {
	store := k.store(ctx)
	bz := store.Get(PurchaseOrderIDKey)
	if bz == nil {
		return 0, ErrUnknownPurchaseOrder(purchaseorderID)
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &purchaseorderID)
	return purchaseorderID, nil
}

// set the purchaseorder ID
func (k Keeper) setPurchaseOrderID(ctx sdk.Context, purchaseorderID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(purchaseorderID)
	store.Set(PurchaseOrderIDKey, bz)
}

// setPurchaseOrder sets a purchaseorder in store
func (k Keeper) setPurchaseOrder(ctx sdk.Context, purchaseorder PurchaseOrder) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(purchaseorder)
	store.Set(key(purchaseorder.ID), bz)
}

// setMarketPurchaseOrder sets a market <-> purchaseorder association in store
func (k Keeper) setMarketPurchaseOrder(ctx sdk.Context, MarketID string, purchaseorderID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(purchaseorderID)
	store.Set(merchantPurchaseOrderKey(merchantID, purchaseorderID), bz)
}

// setMerchantPurchaseOrder sets a merchant <-> purchaseorder association in store
func (k Keeper) setMerchantPurchaseOrder(ctx sdk.Context, merchant sdk.AccAddress, purchaseorderID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(purchaseorderID)
	store.Set(merchantPurchaseOrderKey(merchant, purchaseorderID), bz)
}

func (k Keeper) setCreatedTimePurchaseOrder(ctx sdk.Context, createdTime time.Time, purchaseorderID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(purchaseorderID)
	store.Set(createdTimePurchaseOrderKey(createdTime, purchaseorderID), bz)
}

// purchaseordersIterator returns an sdk.Iterator for purchaseorders from startPurchaseOrderID to endPurchaseOrderID
func (k Keeper) purchaseordersIterator(ctx sdk.Context, startPurchaseOrderID, endPurchaseOrderID uint64) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(key(startPurchaseOrderID), sdk.PrefixEndBytes(key(endPurchaseOrderID)))
}

func (k Keeper) beforeCreatedTimePurchaseOrdersIterator(ctx sdk.Context, createdTime time.Time) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(CreatedTimePurchaseOrdersPrefix, sdk.PrefixEndBytes(createdTimePurchaseOrdersKey(createdTime)))
}

func (k Keeper) afterCreatedTimePurchaseOrdersIterator(ctx sdk.Context, createdTime time.Time) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(createdTimePurchaseOrdersKey(createdTime), sdk.PrefixEndBytes(CreatedTimePurchaseOrdersPrefix))
}

// createdTimeRangePurchaseOrdersIterator returns an sdk.Iterator for all purchaseorders between startCreatedTime and endCreatedTime
func (k Keeper) createdTimeRangePurchaseOrdersIterator(ctx sdk.Context, startCreatedTime, endCreatedTime time.Time) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(createdTimePurchaseOrdersKey(startCreatedTime), sdk.PrefixEndBytes(createdTimePurchaseOrdersKey(endCreatedTime)))
}

func (k Keeper) associatedPurchaseOrders(ctx sdk.Context, prefix []byte) (purchaseorders PurchaseOrders) {
	store := k.store(ctx)
	iterator := sdk.KVStoreReversePrefixIterator(store, prefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchaseorderID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &purchaseorderID)
		purchaseorder, ok := k.PurchaseOrder(ctx, purchaseorderID)
		if ok {
			purchaseorders = append(purchaseorders, purchaseorder)
		}
	}

	return
}

func (k Keeper) iterate(iterator sdk.Iterator) (purchaseorders PurchaseOrders) {
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchaseorder PurchaseOrder
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &purchaseorder)
		purchaseorders = append(purchaseorders, purchaseorder)
	}

	return
}

func (k Keeper) iterateAssociated(ctx sdk.Context, iterator sdk.Iterator) (purchaseorders PurchaseOrders) {
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchaseorderID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &purchaseorderID)
		purchaseorder, ok := k.PurchaseOrder(ctx, purchaseorderID)
		if ok {
			purchaseorders = append(purchaseorders, purchaseorder)
		}
	}

	return
}

func (k Keeper) UpdatePurchaseorder(ctx sdk.Context, purchaseorder types.Purchaseorder) {
	store :=  prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PurchaseorderKey))
	b := k.cdc.MustMarshalBinaryBare(&purchaseorder)
	store.Set(types.KeyPrefix(types.PurchaseorderKey + purchaseorder.Id), b)
}

func (k Keeper) GetPurchaseorder(ctx sdk.Context, key string) types.Purchaseorder {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PurchaseorderKey))
	var purchaseorder types.Purchaseorder
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.PurchaseorderKey + key)), &purchaseorder)
	return purchaseorder
}

func (k Keeper) HasPurchaseorder(ctx sdk.Context, id string) bool {
	store :=  prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PurchaseorderKey))
	return store.Has(types.KeyPrefix(types.PurchaseorderKey + id))
}

func (k Keeper) GetPurchaseorderOwner(ctx sdk.Context, key string) string {
    return k.GetPurchaseorder(ctx, key).Creator
}

// DeletePurchaseorder deletes a purchaseorder
func (k Keeper) DeletePurchaseorder(ctx sdk.Context, key string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PurchaseorderKey))
	store.Delete(types.KeyPrefix(types.PurchaseorderKey + key))
}

func (k Keeper) GetAllPurchaseorder(ctx sdk.Context) (msgs []types.Purchaseorder) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PurchaseorderKey))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.PurchaseorderKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.Purchaseorder
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
        msgs = append(msgs, msg)
	}

    return
}


func (k Keeper) store(ctx sdk.Context) sdk.KVStore {
	return gaskv.NewStore(ctx.MultiStore().GetKVStore(k.storeKey), ctx.GasMeter(), app.KVGasConfig())
}

func logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}