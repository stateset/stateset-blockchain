package purchaseorder

import (
	"net/url"
	"time"

	app "github.com/stateset/stateset-blockchain/types"
	"github.com/stateset/stateset-blockchain/x/marketplace"
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
	marketplaceKeeper marketplace.Keeper
}

// NewKeeper creates a new account keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec, accountKeeper AccountKeeper, marketplaceKeeper marketplace.Keeper) Keeper {
	return Keeper{
		storeKey,
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
		accountKeeper,
		marketplaceKeeper,
	}
}

// CreatePurchaseOrder creates a new purchaseorder in the purchaseorder key-value store
func (k Keeper) CreatePurchaseOrder(ctx sdk.Context, body, purchaseorderID string,
	merchant sdk.AccAddress, source url.URL) (purchaseorder PurchaseOrder, err sdk.Error) {

	err = k.validateLength(ctx, body)
	if err != nil {
		return
	}
	jailed, err := k.purchaseorderKeeper.IsJailed(ctx, merchant)
	if err != nil {
		return
	}
	if jailed {
		return purchaseorder, ErrMerchantJailed(merchant)
	}
	marketplace, err := k.marketplaceKeeper.Marketplace(ctx, marketplaceID)
	if err != nil {
		return purchaseorder, ErrInvalidMarketplaceID(marketplace.ID)
	}

	purchaseorderID, err := k.purchaseorderID(ctx)
	if err != nil {
		return
	}
	purchaseorder = NewPurchaseOrder(purchaseorderID, marketplaceID, body, merchant, source,
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

// MarketplacePurchaseOrders gets all the purchaseorders for a given marketplace
func (k Keeper) MarketplacePurchaseOrders(ctx sdk.Context, marketplaceID string) (purchaseorders PurchaseOrders) {
	return k.associatedPurchaseOrders(ctx, marketplacePurchaseOrdersKey(marketplaceID))
}

// MerchantPurchaseOrders gets all the purchaseorders for a given merchant
func (k Keeper) MerchantPurchaseOrders(ctx sdk.Context, creator sdk.AccAddress) (purchaseorders PurchaseOrders) {
	return k.associatedPurchaseOrders(ctx, creatorPurchaseOrdersKey(creator))
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
func (k Keeper) SubtractFactorStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	purchaseorder, ok := k.Invoice(ctx, id)
	if !ok {
		return ErrUnknownInvoice(id)
	}
	purchaseorder.TotalFactored = purchaseorder.TotalFactored.Sub(stake)
	k.setIncoice(ctx, purchaseorder)

	return nil
}

// SubtractChallengeStake adds a stake amount to the total factoring amount
func (k Keeper) SubtractFactorStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	purchaseorder, ok := k.Invoice(ctx, id)
	if !ok {
		return ErrUnknownInvoice(id)
	}
	purchaseorder.TotalFactored = purchaseorder.TotalFactored.Sub(stake)
	k.setInvoice(ctx, purchaseorder)

	return nil
}

// SetFirstArgumentTime sets time when first argument was created on a purchaseorder
func (k Keeper) SetFirstArgumentTime(ctx sdk.Context, id uint64, firstArgumentTime time.Time) sdk.Error {
	purchaseorder, ok := k.Invoice(ctx, id)
	if !ok {
		return ErrUnknownInvoice(id)
	}
	purchaseorder.FirstArgumentTime = firstArgumentTime
	k.setInvoice(ctx, purchaseorder)

	return nil
}

// AddAdmin adds a new admin
func (k Keeper) AddAdmin(ctx sdk.Context, admin, creator sdk.AccAddress) (err sdk.Error) {
	params := k.GetParams(ctx)

	// first admin can be added without any authorisation
	if len(params.InvoiceAdmins) > 0 && !k.isAdmin(ctx, creator) {
		err = ErrAddressNotAuthorised()
	}

	// if already present, don't add again
	for _, currentAdmin := range params.InvoiceAdmins {
		if currentAdmin.Equals(admin) {
			return
		}
	}

	params.InvoiceAdmins = append(params.InvoiceAdmins, admin)

	k.SetParams(ctx, params)

	return
}

// RemoveAdmin removes an admin
func (k Keeper) RemoveAdmin(ctx sdk.Context, admin, remover sdk.AccAddress) (err sdk.Error) {
	if !k.isAdmin(ctx, remover) {
		err = ErrAddressNotAuthorised()
	}

	params := k.GetParams(ctx)
	for i, currentAdmin := range params.InvoiceAdmins {
		if currentAdmin.Equals(admin) {
			params.InvoiceAdmins = append(params.InvoiceAdmins[:i], params.InvoiceAdmins[i+1:]...)
		}
	}

	k.SetParams(ctx, params)

	return
}

func (k Keeper) isAdmin(ctx sdk.Context, address sdk.AccAddress) bool {
	for _, admin := range k.GetParams(ctx).InvoiceAdmins {
		if address.Equals(admin) {
			return true
		}
	}
	return false
}

func (k Keeper) validateLength(ctx sdk.Context, body string) sdk.Error {
	var minInvoiceLength int
	var maxInvoiceLength int

	k.paramStore.Get(ctx, KeyMinInvoiceLength, &minInvoiceLength)
	k.paramStore.Get(ctx, KeyMaxInvoiceLength, &maxInvoiceLength)

	len := len([]rune(body))
	if len < minInvoiceLength {
		return ErrInvalidBodyTooShort(body)
	}
	if len > maxInvoiceLength {
		return ErrInvalidBodyTooLong()
	}

	return nil
}

// purchaseorderID gets the highest purchaseorder ID
func (k Keeper) purchaseorderID(ctx sdk.Context) (purchaseorderID uint64, err sdk.Error) {
	store := k.store(ctx)
	bz := store.Get(InvoiceIDKey)
	if bz == nil {
		return 0, ErrUnknownInvoice(purchaseorderID)
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &purchaseorderID)
	return purchaseorderID, nil
}

// set the purchaseorder ID
func (k Keeper) setInvoiceID(ctx sdk.Context, purchaseorderID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(purchaseorderID)
	store.Set(InvoiceIDKey, bz)
}

// setInvoice sets a purchaseorder in store
func (k Keeper) setInvoice(ctx sdk.Context, purchaseorder Invoice) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(purchaseorder)
	store.Set(key(purchaseorder.ID), bz)
}

// setMarketplaceInvoice sets a marketplace <-> purchaseorder association in store
func (k Keeper) setMarketplaceInvoice(ctx sdk.Context, marketplaceID string, purchaseorderID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(purchaseorderID)
	store.Set(merchantInvoiceKey(merchantID, purchaseorderID), bz)
}

// setMerchantInvoice sets a merchant <-> purchaseorder association in store
func (k Keeper) setMerchantInvoice(ctx sdk.Context, merchant sdk.AccAddress, purchaseorderID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(purchaseorderID)
	store.Set(merchantInvoiceKey(merchant, purchaseorderID), bz)
}

func (k Keeper) setCreatedTimeInvoice(ctx sdk.Context, createdTime time.Time, purchaseorderID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(purchaseorderID)
	store.Set(createdTimeInvoiceKey(createdTime, purchaseorderID), bz)
}

// purchaseordersIterator returns an sdk.Iterator for purchaseorders from startInvoiceID to endInvoiceID
func (k Keeper) purchaseordersIterator(ctx sdk.Context, startInvoiceID, endInvoiceID uint64) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(key(startInvoiceID), sdk.PrefixEndBytes(key(endInvoiceID)))
}

func (k Keeper) beforeCreatedTimeInvoicesIterator(ctx sdk.Context, createdTime time.Time) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(CreatedTimeInvoicesPrefix, sdk.PrefixEndBytes(createdTimeInvoicesKey(createdTime)))
}

func (k Keeper) afterCreatedTimeInvoicesIterator(ctx sdk.Context, createdTime time.Time) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(createdTimeInvoicesKey(createdTime), sdk.PrefixEndBytes(CreatedTimeInvoicesPrefix))
}

// createdTimeRangeInvoicesIterator returns an sdk.Iterator for all purchaseorders between startCreatedTime and endCreatedTime
func (k Keeper) createdTimeRangeInvoicesIterator(ctx sdk.Context, startCreatedTime, endCreatedTime time.Time) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(createdTimeInvoicesKey(startCreatedTime), sdk.PrefixEndBytes(createdTimeInvoicesKey(endCreatedTime)))
}

func (k Keeper) associatedInvoices(ctx sdk.Context, prefix []byte) (purchaseorders Invoices) {
	store := k.store(ctx)
	iterator := sdk.KVStoreReversePrefixIterator(store, prefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchaseorderID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &purchaseorderID)
		purchaseorder, ok := k.Invoice(ctx, purchaseorderID)
		if ok {
			purchaseorders = append(purchaseorders, purchaseorder)
		}
	}

	return
}

func (k Keeper) iterate(iterator sdk.Iterator) (purchaseorders Invoices) {
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchaseorder Invoice
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &purchaseorder)
		purchaseorders = append(purchaseorders, purchaseorder)
	}

	return
}

func (k Keeper) iterateAssociated(ctx sdk.Context, iterator sdk.Iterator) (purchaseorders Invoices) {
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchaseorderID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &purchaseorderID)
		purchaseorder, ok := k.Invoice(ctx, purchaseorderID)
		if ok {
			purchaseorders = append(purchaseorders, purchaseorder)
		}
	}

	return
}

func (k Keeper) store(ctx sdk.Context) sdk.KVStore {
	return gaskv.NewStore(ctx.MultiStore().GetKVStore(k.storeKey), ctx.GasMeter(), app.KVGasConfig())
}

func logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}