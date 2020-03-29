package agreement

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

// CreateAgreement creates a new agreement in the agreement key-value store
func (k Keeper) CreateAgreement(ctx sdk.Context, body, agreementID string,
	merchant sdk.AccAddress, source url.URL) (agreement Agreement, err sdk.Error) {

	err = k.validateLength(ctx, body)
	if err != nil {
		return
	}
	jailed, err := k.agreementKeeper.IsJailed(ctx, merchant)
	if err != nil {
		return
	}
	if jailed {
		return agreement, ErrMerchantJailed(merchant)
	}
	marketplace, err := k.marketplaceKeeper.Marketplace(ctx, marketplaceID)
	if err != nil {
		return invoice, ErrInvalidMarketplaceID(marketplace.ID)
	}

	agreementID, err := k.agreementID(ctx)
	if err != nil {
		return
	}
	agreement = NewAgreement(agreementID, marketplaceID, body, merchant, source,
		ctx.BlockHeader().Time,
	)

	// persist agreement
	k.setAgreement(ctx, agreement)
	// increment agreementID (primary key) for next agreement
	k.setAgreementID(ctx, agreementID+1)

	// persist associations
	k.setPartyAgreement(ctx, agreement.Controller, agreementID)
	k.setProcessorAgreement(ctx, agreement.Processor, agreementId)
	k.setCreatedTimeAgreement(ctx, agreement.CreatedTime, agreementID)

	logger(ctx).Info("Submitted " + agreement.String())

	return agreement, nil
}

// EditInvoice allows admins to edit the body of an invoice

func (k Keeper) EditAgreement(ctx sdk.Context, id uint64, body string, editor sdk.AccAddress) (agreement Agreement, err sdk.Error) {
	if !k.isAdmin(ctx, editor) {
		err = ErrAddressNotAuthorised()
		return
	}

	err = k.validateLength(ctx, body)
	if err != nil {
		return
	}

	agreement, ok := k.Agreement(ctx, id)
	if !ok {
		err = ErrUnknownAgreement(id)
		return
	}

	agreement.Body = body
	k.setAgreement(ctx, agreement)

	return
}

// Invoice gets a single invoice by its ID
func (k Keeper) Agreement(ctx sdk.Context, id uint64) (agreement Agreement, ok bool) {
	store := k.store(ctx)
	agreementBytes := store.Get(key(id))
	if agreementBytes == nil {
		return agreement, false
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(agreementBytes, &agreement)

	return agreement, true
}

// Agreements gets all the agreements in reverse order
func (k Keeper) Agreements(ctx sdk.Context) (agreements Agreements) {
	store := k.store(ctx)
	iterator := sdk.KVStoreReversePrefixIterator(store, AgreementsKeyPrefix)

	return k.iterate(iterator)
}

// AgreementsBetweenIDs gets all invoices between startAgreementID to endAgreementID
func (k Keeper) AgreementsBetweenIDs(ctx sdk.Context, startAgreementID, endAgreementID uint64) (agreements Agreements) {
	iterator := k.agreementsIterator(ctx, startAgreementID, endAgreementID)

	return k.iterate(iterator)
}

// AgreementsBetweenTimes gets all agreements between startTime and endTime
func (k Keeper) AgreementsBetweenTimes(ctx sdk.Context, startTime time.Time, endTime time.Time) (agreements Agreements) {
	iterator := k.createdTimeRangeAgreementsIterator(ctx, startTime, endTime)

	return k.iterateAssociated(ctx, iterator)
}

// AgreementsBeforeTime gets all invoices after a certain CreatedTime
func (k Keeper) AgreementsBeforeTime(ctx sdk.Context, createdTime time.Time) (agreements Agreements) {
	iterator := k.beforeCreatedTimeInvoicesIterator(ctx, createdTime)

	return k.iterateAssociated(ctx, iterator)
}

// AgreementsAfterTime gets all agreements after a certain CreatedTime
func (k Keeper) AgreementsAfterTime(ctx sdk.Context, createdTime time.Time) (agreements Agreements) {
	iterator := k.afterCreatedTimeAgreementsIterator(ctx, createdTime)

	return k.iterateAssociated(ctx, iterator)
}

// MarketplaceAgreements gets all the agreements for a given markerplace
func (k Keeper) MarketplaceAgreements(ctx sdk.Context, marketplaceID string) (agreements Agreements) {
	return k.associatedAgreements(ctx, marketplaceAgreementsKey(marketplaceID))
}

// MerchantInvoices gets all the agreements for a given merchant
func (k Keeper) MerchantAgreements(ctx sdk.Context, creator sdk.AccAddress) (agreements Agreements) {
	return k.associatedAgreements(ctx, merchantAgreementsKey(merchant))
}


// agreememtID gets the highest agreement ID
func (k Keeper) invoiceID(ctx sdk.Context) (agreementID uint64, err sdk.Error) {
	store := k.store(ctx)
	bz := store.Get(AgreementIDKey)
	if bz == nil {
		return 0, ErrUnknownAgreement(agreementID)
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &agreementID)
	return agreementID, nil
}

// set the agreement ID
func (k Keeper) setAgreementID(ctx sdk.Context, agreementID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(agreementID)
	store.Set(AgreementIDKey, bz)
}

// setAgreement sets a agreement in store
func (k Keeper) setAgreement(ctx sdk.Context, agreement Agreement) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(agreement)
	store.Set(key(agreement.ID), bz)
}

// setMarketplaceInvoice sets a marketplace <-> invoice association in store
func (k Keeper) setMarketplaceAgreement(ctx sdk.Context, agreementID string, agreementID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(agreementID)
	store.Set(merchantAgreementKey(agreementID, agreementID), bz)
}

// setMerchantInvoice sets a merchant <-> invoice association in store
func (k Keeper) setMerchantAgreement(ctx sdk.Context, merchant sdk.AccAddress, agreementID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(agreementID)
	store.Set(merchantInvoiceKey(merchant, invoiceID), bz)
}


func (k Keeper) setCreatedTimeAgreement(ctx sdk.Context, createdTime time.Time, agreementID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(agreementID)
	store.Set(createdTimeAgreementKey(createdTime, agreementID), bz)
}


// agreementsIterator returns an sdk.Iterator for agreement from startAgreementID to endAgreementID
func (k Keeper) agreementsIterator(ctx sdk.Context, startAgreementID, endAgreementID uint64) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(key(startAgreementID), sdk.PrefixEndBytes(key(endAgreementID)))
}

func (k Keeper) beforeCreatedTimeAgreementsIterator(ctx sdk.Context, createdTime time.Time) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(CreatedTimeAgreementsPrefix, sdk.PrefixEndBytes(createdTimeAgreementsKey(createdTime)))
}

func (k Keeper) afterCreatedTimeAgreementsIterator(ctx sdk.Context, createdTime time.Time) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(createdTimeAgreementsKey(createdTime), sdk.PrefixEndBytes(CreatedTimeAgreementsPrefix))
}

// createdTimeRangeInvoicesIterator returns an sdk.Iterator for all invoices between startCreatedTime and endCreatedTime
func (k Keeper) createdTimeRangeInvoicesIterator(ctx sdk.Context, startCreatedTime, endCreatedTime time.Time) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(createdTimeAgreementsKey(startCreatedTime), sdk.PrefixEndBytes(createdTimeAgreementsKey(endCreatedTime)))
}

func (k Keeper) associatedAgreements(ctx sdk.Context, prefix []byte) (agreements Agreements) {
	store := k.store(ctx)
	iterator := sdk.KVStoreReversePrefixIterator(store, prefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var invoiceID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &agreementID)
		agreement, ok := k.Agreement(ctx, agreementID)
		if ok {
			agreements = append(agreements, agreement)
		}
	}

	return
}

func (k Keeper) iterate(iterator sdk.Iterator) (agreements Agreements) {
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var agreement Agreement
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &agreement)
		agreements = append(agreements, agreement)
	}

	return
}

func (k Keeper) iterateAssociated(ctx sdk.Context, iterator sdk.Iterator) (agreements Agreements) {
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var agreementID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &invoiceID)
		invoice, ok := k.Invoice(ctx, invoiceID)
		if ok {
			agreements = append(agreements, agreement)
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