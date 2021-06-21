package invoice

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
	codec      *codec.LegacyAmino
	paramStore params.Subspace
	accountKeeper   AccountKeeper
}

// NewKeeper creates a new account keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.LegacyAmino, accountKeeper AccountKeeper) Keeper {
	return Keeper{
		storeKey,
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
		accountKeeper
	}
}

// CreateInvoice creates a new invoice in the invoice key-value store
func (k Keeper) CreateInvoice(ctx sdk.Context, body, invoiceID string,
	merchant sdk.AccAddress, source url.URL) (invoice Invoice, err sdk.Error) {

	err = k.validateLength(ctx, body)
	if err != nil {
		return
	}

	invoiceID, err := k.invoiceID(ctx)
	if err != nil {
		return
	}
	invoice = NewInvoice(invoiceID, body, merchant, source,
		ctx.BlockHeader().Time,
	)

	// persist invoice
	k.setInvoice(ctx, invoice)
	// increment invoiceID (primary key) for next invoice
	k.setInvoiceID(ctx, invoiceID+1)

	// persist associations
	k.setControllerInvoice(ctx, invoice.Controller, invoiceID)
	k.setProcessorInvoice(ctx, invoice.Processor, invoiceId)
	k.setCreatedTimeInvoice(ctx, invoice.CreatedTime, invoiceID)

	logger(ctx).Info("Submitted " + invoice.String())

	return invoice, nil
}

// EditInvoice allows admins to edit the body of an invoice

func (k Keeper) EditInvoice(ctx sdk.Context, id uint64, body string, editor sdk.AccAddress) (invoice Invoice, err sdk.Error) {
	if !k.isAdmin(ctx, editor) {
		err = ErrAddressNotAuthorised()
		return
	}

	err = k.validateLength(ctx, body)
	if err != nil {
		return
	}

	invoice, ok := k.Invoice(ctx, id)
	if !ok {
		err = ErrUnknownInvoice(id)
		return
	}

	invoice.Body = body
	k.setInvoice(ctx, invoice)

	return
}

// Invoice gets a single invoice by its ID
func (k Keeper) Invoice(ctx sdk.Context, id uint64) (invoice Invoice, ok bool) {
	store := k.store(ctx)
	invoiceBytes := store.Get(key(id))
	if invoiceBytes == nil {
		return invoice, false
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(invoiceBytes, &invoice)

	return invoice, true
}

// Invoices gets all the invoices in reverse order
func (k Keeper) Invoices(ctx sdk.Context) (invoices Invoices) {
	store := k.store(ctx)
	iterator := sdk.KVStoreReversePrefixIterator(store, InvoicesKeyPrefix)

	return k.iterate(iterator)
}

// InvoicesBetweenIDs gets all invoices between startInvoiceID to endInvoiceID
func (k Keeper) InvoicesBetweenIDs(ctx sdk.Context, startInvoiceID, endInvoiceID uint64) (invoices Invoices) {
	iterator := k.invoicesIterator(ctx, startInvoiceID, endInvoiceID)

	return k.iterate(iterator)
}

// InvoicesBetweenTimes gets all invoices between startTime and endTime
func (k Keeper) InvoicesBetweenTimes(ctx sdk.Context, startTime time.Time, endTime time.Time) (invoices Invoices) {
	iterator := k.createdTimeRangeInvoicesIterator(ctx, startTime, endTime)

	return k.iterateAssociated(ctx, iterator)
}

// InvoicessBeforeTime gets all invoices after a certain CreatedTime
func (k Keeper) InvoicesBeforeTime(ctx sdk.Context, createdTime time.Time) (invoices Invoices) {
	iterator := k.beforeCreatedTimeInvoicesIterator(ctx, createdTime)

	return k.iterateAssociated(ctx, iterator)
}

// InvoicessAfterTime gets all invoices after a certain CreatedTime
func (k Keeper) InvoicessAfterTime(ctx sdk.Context, createdTime time.Time) (invoices Invoices) {
	iterator := k.afterCreatedTimeInvoicesIterator(ctx, createdTime)

	return k.iterateAssociated(ctx, iterator)
}

// MerchantInvoices gets all the invoices for a given merchant
func (k Keeper) CreatorInvoices(ctx sdk.Context, creator sdk.AccAddress) (invoices Invoices) {
	return k.associatedInvoices(ctx, creatorInvoicesKey(creator))
}

// AddBackingStake adds a stake amount to the total backing amount
func (k Keeper) AddBackingStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	invoice, ok := k.Invoice(ctx, id)
	if !ok {
		return ErrUnknownInvoice(id)
	}
	invoice.TotalBacked = invoice.TotalBacked.Add(stake)
	invoice.TotalStakers++
	k.setInvoice(ctx, invoice)

	return nil
}

// AddFactorStake adds a stake amount to the total factor amount
func (k Keeper) AddFactorStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	, ok := k.Invoice(ctx, id)
	if !ok {
		return ErrUnknownInvoice(id)
	}
	invoice.TotalFactored = invoice.TotalFactored.Add(stake)
	invoice.TotalStakers++
	k.setInvoice(ctx, invoice)

	return nil
}

// SubtractFactorStake adds a stake amount to the total factoring amount
func (k Keeper) SubtractFactorStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	invoice, ok := k.Invoice(ctx, id)
	if !ok {
		return ErrUnknownInvoice(id)
	}
	invoice.TotalFactored = invoice.TotalFactored.Sub(stake)
	k.setIncoice(ctx, invoice)

	return nil
}

// SubtractChallengeStake adds a stake amount to the total factoring amount
func (k Keeper) SubtractFactorStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	invoice, ok := k.Invoice(ctx, id)
	if !ok {
		return ErrUnknownInvoice(id)
	}
	invoice.TotalFactored = invoice.TotalFactored.Sub(stake)
	k.setInvoice(ctx, invoice)

	return nil
}

// SetFirstArgumentTime sets time when first argument was created on a invoice
func (k Keeper) SetFirstArgumentTime(ctx sdk.Context, id uint64, firstArgumentTime time.Time) sdk.Error {
	invoice, ok := k.Invoice(ctx, id)
	if !ok {
		return ErrUnknownInvoice(id)
	}
	invoice.FirstArgumentTime = firstArgumentTime
	k.setInvoice(ctx, invoice)

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

// invoiceID gets the highest invoice ID
func (k Keeper) invoiceID(ctx sdk.Context) (invoiceID uint64, err sdk.Error) {
	store := k.store(ctx)
	bz := store.Get(InvoiceIDKey)
	if bz == nil {
		return 0, ErrUnknownInvoice(invoiceID)
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &invoiceID)
	return invoiceID, nil
}

// set the invoice ID
func (k Keeper) setInvoiceID(ctx sdk.Context, invoiceID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(invoiceID)
	store.Set(InvoiceIDKey, bz)
}

// setInvoice sets a invoice in store
func (k Keeper) setInvoice(ctx sdk.Context, invoice Invoice) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(invoice)
	store.Set(key(invoice.ID), bz)
}

// setMerchantInvoice sets a merchant <-> invoice association in store
func (k Keeper) setMerchantInvoice(ctx sdk.Context, merchant sdk.AccAddress, invoiceID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(invoiceID)
	store.Set(merchantInvoiceKey(merchant, invoiceID), bz)
}

func (k Keeper) setCreatedTimeInvoice(ctx sdk.Context, createdTime time.Time, invoiceID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(invoiceID)
	store.Set(createdTimeInvoiceKey(createdTime, invoiceID), bz)
}

// invoicesIterator returns an sdk.Iterator for invoices from startInvoiceID to endInvoiceID
func (k Keeper) invoicesIterator(ctx sdk.Context, startInvoiceID, endInvoiceID uint64) sdk.Iterator {
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

// createdTimeRangeInvoicesIterator returns an sdk.Iterator for all invoices between startCreatedTime and endCreatedTime
func (k Keeper) createdTimeRangeInvoicesIterator(ctx sdk.Context, startCreatedTime, endCreatedTime time.Time) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(createdTimeInvoicesKey(startCreatedTime), sdk.PrefixEndBytes(createdTimeInvoicesKey(endCreatedTime)))
}

func (k Keeper) associatedInvoices(ctx sdk.Context, prefix []byte) (invoices Invoices) {
	store := k.store(ctx)
	iterator := sdk.KVStoreReversePrefixIterator(store, prefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var invoiceID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &invoiceID)
		invoice, ok := k.Invoice(ctx, invoiceID)
		if ok {
			invoices = append(invoices, invoice)
		}
	}

	return
}

func (k Keeper) iterate(iterator sdk.Iterator) (invoices Invoices) {
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var invoice Invoice
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &invoice)
		invoices = append(invoices, invoice)
	}

	return
}

func (k Keeper) iterateAssociated(ctx sdk.Context, iterator sdk.Iterator) (invoices Invoices) {
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var invoiceID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &invoiceID)
		invoice, ok := k.Invoice(ctx, invoiceID)
		if ok {
			invoices = append(invoices, invoice)
		}
	}

	return
}

func (k Keeper) UpdateInvoice(ctx sdk.Context, invoice types.Invoice) {
	store :=  prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InvoiceKey))
	b := k.cdc.MustMarshalBinaryBare(&invoice)
	store.Set(types.KeyPrefix(types.InvoiceKey + invoice.Id), b)
}

func (k Keeper) GetInvoice(ctx sdk.Context, key string) types.Invoice {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InvoiceKey))
	var invoice types.Invoice
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.InvoiceKey + key)), &invoice)
	return invoice
}

func (k Keeper) HasInvoice(ctx sdk.Context, id string) bool {
	store :=  prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InvoiceKey))
	return store.Has(types.KeyPrefix(types.InvoiceKey + id))
}

func (k Keeper) GetInvoiceOwner(ctx sdk.Context, key string) string {
    return k.GetInvoice(ctx, key).Creator
}

// DeleteInvoice deletes a invoice
func (k Keeper) DeleteInvoice(ctx sdk.Context, key string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InvoiceKey))
	store.Delete(types.KeyPrefix(types.InvoiceKey + key))
}

func (k Keeper) GetAllInvoice(ctx sdk.Context) (msgs []types.Invoice) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InvoiceKey))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.InvoiceKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.Invoice
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