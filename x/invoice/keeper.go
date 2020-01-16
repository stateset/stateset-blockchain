package invoice

import (
	"net/url"
	"time"

	app "github.com/stateset/stateset/types"
	"github.com/stateset/stateset/x/community"
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
	communityKeeper community.Keeper
}

// NewKeeper creates a new account keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec, accountKeeper AccountKeeper, communityKeeper community.Keeper) Keeper {
	return Keeper{
		storeKey,
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
		accountKeeper,
		communityKeeper,
	}
}

// SubmitInvoice creates a new invoice in the invoice key-value store
func (k Keeper) SubmitInvoice(ctx sdk.Context, body, invoiceID string,
	creator sdk.AccAddress, source url.URL) (invoice Invoice, err sdk.Error) {

	err = k.validateLength(ctx, body)
	if err != nil {
		return
	}
	jailed, err := k.invoiceKeeper.IsJailed(ctx, creator)
	if err != nil {
		return
	}
	if jailed {
		return claim, ErrCreatorJailed(creator)
	}
	community, err := k.communityKeeper.Community(ctx, communityID)
	if err != nil {
		return claim, ErrInvalidCommunityID(community.ID)
	}

	invoiceID, err := k.invoiceID(ctx)
	if err != nil {
		return
	}
	invoice = NewInvoice(invoiceID, communityID, body, creator, source,
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

// ClaimsBeforeTime gets all claims after a certain CreatedTime
func (k Keeper) ClaimsBeforeTime(ctx sdk.Context, createdTime time.Time) (claims Claims) {
	iterator := k.beforeCreatedTimeClaimsIterator(ctx, createdTime)

	return k.iterateAssociated(ctx, iterator)
}

// ClaimsAfterTime gets all claims after a certain CreatedTime
func (k Keeper) ClaimsAfterTime(ctx sdk.Context, createdTime time.Time) (claims Claims) {
	iterator := k.afterCreatedTimeClaimsIterator(ctx, createdTime)

	return k.iterateAssociated(ctx, iterator)
}

// CommunityClaims gets all the claims for a given community
func (k Keeper) CommunityClaims(ctx sdk.Context, communityID string) (claims Claims) {
	return k.associatedClaims(ctx, communityClaimsKey(communityID))
}

// CreatorClaims gets all the claims for a given creator
func (k Keeper) CreatorClaims(ctx sdk.Context, creator sdk.AccAddress) (claims Claims) {
	return k.associatedClaims(ctx, creatorClaimsKey(creator))
}

// AddBackingStake adds a stake amount to the total backing amount
func (k Keeper) AddBackingStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	claim, ok := k.Claim(ctx, id)
	if !ok {
		return ErrUnknownClaim(id)
	}
	claim.TotalBacked = claim.TotalBacked.Add(stake)
	claim.TotalStakers++
	k.setClaim(ctx, claim)

	return nil
}

// AddChallengeStake adds a stake amount to the total challenge amount
func (k Keeper) AddChallengeStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	claim, ok := k.Claim(ctx, id)
	if !ok {
		return ErrUnknownClaim(id)
	}
	claim.TotalChallenged = claim.TotalChallenged.Add(stake)
	claim.TotalStakers++
	k.setClaim(ctx, claim)

	return nil
}

// SubtractBackingStake adds a stake amount to the total backing amount
func (k Keeper) SubtractBackingStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	claim, ok := k.Claim(ctx, id)
	if !ok {
		return ErrUnknownClaim(id)
	}
	claim.TotalBacked = claim.TotalBacked.Sub(stake)
	k.setClaim(ctx, claim)

	return nil
}

// SubtractChallengeStake adds a stake amount to the total challenge amount
func (k Keeper) SubtractChallengeStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	claim, ok := k.Claim(ctx, id)
	if !ok {
		return ErrUnknownClaim(id)
	}
	claim.TotalChallenged = claim.TotalChallenged.Sub(stake)
	k.setClaim(ctx, claim)

	return nil
}

// SetFirstArgumentTime sets time when first argument was created on a claim
func (k Keeper) SetFirstArgumentTime(ctx sdk.Context, id uint64, firstArgumentTime time.Time) sdk.Error {
	claim, ok := k.Claim(ctx, id)
	if !ok {
		return ErrUnknownClaim(id)
	}
	claim.FirstArgumentTime = firstArgumentTime
	k.setClaim(ctx, claim)

	return nil
}

// AddAdmin adds a new admin
func (k Keeper) AddAdmin(ctx sdk.Context, admin, creator sdk.AccAddress) (err sdk.Error) {
	params := k.GetParams(ctx)

	// first admin can be added without any authorisation
	if len(params.ClaimAdmins) > 0 && !k.isAdmin(ctx, creator) {
		err = ErrAddressNotAuthorised()
	}

	// if already present, don't add again
	for _, currentAdmin := range params.ClaimAdmins {
		if currentAdmin.Equals(admin) {
			return
		}
	}

	params.ClaimAdmins = append(params.ClaimAdmins, admin)

	k.SetParams(ctx, params)

	return
}

// RemoveAdmin removes an admin
func (k Keeper) RemoveAdmin(ctx sdk.Context, admin, remover sdk.AccAddress) (err sdk.Error) {
	if !k.isAdmin(ctx, remover) {
		err = ErrAddressNotAuthorised()
	}

	params := k.GetParams(ctx)
	for i, currentAdmin := range params.ClaimAdmins {
		if currentAdmin.Equals(admin) {
			params.ClaimAdmins = append(params.ClaimAdmins[:i], params.ClaimAdmins[i+1:]...)
		}
	}

	k.SetParams(ctx, params)

	return
}

func (k Keeper) isAdmin(ctx sdk.Context, address sdk.AccAddress) bool {
	for _, admin := range k.GetParams(ctx).ClaimAdmins {
		if address.Equals(admin) {
			return true
		}
	}
	return false
}

func (k Keeper) validateLength(ctx sdk.Context, body string) sdk.Error {
	var minClaimLength int
	var maxClaimLength int

	k.paramStore.Get(ctx, KeyMinClaimLength, &minClaimLength)
	k.paramStore.Get(ctx, KeyMaxClaimLength, &maxClaimLength)

	len := len([]rune(body))
	if len < minClaimLength {
		return ErrInvalidBodyTooShort(body)
	}
	if len > maxClaimLength {
		return ErrInvalidBodyTooLong()
	}

	return nil
}

// claimID gets the highest claim ID
func (k Keeper) claimID(ctx sdk.Context) (claimID uint64, err sdk.Error) {
	store := k.store(ctx)
	bz := store.Get(ClaimIDKey)
	if bz == nil {
		return 0, ErrUnknownClaim(claimID)
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &claimID)
	return claimID, nil
}

// set the claim ID
func (k Keeper) setClaimID(ctx sdk.Context, claimID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(claimID)
	store.Set(ClaimIDKey, bz)
}

// setClaim sets a claim in store
func (k Keeper) setClaim(ctx sdk.Context, claim Claim) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(claim)
	store.Set(key(claim.ID), bz)
}

// setCommunityClaim sets a community <-> claim association in store
func (k Keeper) setCommunityClaim(ctx sdk.Context, communityID string, claimID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(claimID)
	store.Set(communityClaimKey(communityID, claimID), bz)
}

func (k Keeper) setCreatorClaim(ctx sdk.Context, creator sdk.AccAddress, claimID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(claimID)
	store.Set(creatorClaimKey(creator, claimID), bz)
}

func (k Keeper) setCreatedTimeClaim(ctx sdk.Context, createdTime time.Time, claimID uint64) {
	store := k.store(ctx)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(claimID)
	store.Set(createdTimeClaimKey(createdTime, claimID), bz)
}

// claimsIterator returns an sdk.Iterator for claims from startClaimID to endClaimID
func (k Keeper) claimsIterator(ctx sdk.Context, startClaimID, endClaimID uint64) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(key(startClaimID), sdk.PrefixEndBytes(key(endClaimID)))
}

func (k Keeper) beforeCreatedTimeClaimsIterator(ctx sdk.Context, createdTime time.Time) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(CreatedTimeClaimsPrefix, sdk.PrefixEndBytes(createdTimeClaimsKey(createdTime)))
}

func (k Keeper) afterCreatedTimeClaimsIterator(ctx sdk.Context, createdTime time.Time) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(createdTimeClaimsKey(createdTime), sdk.PrefixEndBytes(CreatedTimeClaimsPrefix))
}

// createdTimeRangeClaimsIterator returns an sdk.Iterator for all claims between startCreatedTime and endCreatedTime
func (k Keeper) createdTimeRangeClaimsIterator(ctx sdk.Context, startCreatedTime, endCreatedTime time.Time) sdk.Iterator {
	store := k.store(ctx)
	return store.Iterator(createdTimeClaimsKey(startCreatedTime), sdk.PrefixEndBytes(createdTimeClaimsKey(endCreatedTime)))
}

func (k Keeper) associatedClaims(ctx sdk.Context, prefix []byte) (claims Claims) {
	store := k.store(ctx)
	iterator := sdk.KVStoreReversePrefixIterator(store, prefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var claimID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &claimID)
		claim, ok := k.Claim(ctx, claimID)
		if ok {
			claims = append(claims, claim)
		}
	}

	return
}

func (k Keeper) iterate(iterator sdk.Iterator) (claims Claims) {
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var claim Claim
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &claim)
		claims = append(claims, claim)
	}

	return
}

func (k Keeper) iterateAssociated(ctx sdk.Context, iterator sdk.Iterator) (claims Claims) {
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var claimID uint64
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &claimID)
		claim, ok := k.Claim(ctx, claimID)
		if ok {
			claims = append(claims, claim)
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