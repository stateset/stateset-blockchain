package account

import (
	"net/url"
	"time"

	app "github.com/stateset/stateset/types"
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

// SubmitAccount creates a new account in the account key-value store
func (k Keeper) SubmitAccount(ctx sdk.Context, body, marketplaceID string,
	creator sdk.AccAddress, source url.URL) (account Account, err sdk.Error) {

	err = k.validateLength(ctx, body)
	if err != nil {
		return
	}
	jailed, err := k.accountKeeper.IsJailed(ctx, creator)
	if err != nil {
		return
	}
	if jailed {
		return claim, ErrCreatorJailed(creator)
	}
	marketplace, err := k.marketplaceKeeper.Marketplace(ctx, marketplaceID)
	if err != nil {
		return claim, ErrInvalidMarketplaceID(marketplace.ID)
	}

	accountID, err := k.accountID(ctx)
	if err != nil {
		return
	}
	account = NewAccount(accountID, marketplaceID, body, creator, source,
		ctx.BlockHeader().Time,
	)

	// persist account
	k.setAccount(ctx, account)
	// increment accountID (primary key) for next account
	k.setAccountID(ctx, accountID+1)

	// persist associations
	k.setControllerAccount(ctx, account.Controller, accountID)
	k.setProcessorAccount(ctx, account.Processor, accountId)
	k.setCreatedTimeAccount(ctx, account.CreatedTime, accountID)

	logger(ctx).Info("Submitted " + account.String())

	return account, nil
}

// EditAccount allows admins to edit the body of an account

func (k Keeper) EditAccount(ctx sdk.Context, id uint64, body string, editor sdk.AccAddress) (account Account, err sdk.Error) {
	if !k.isAdmin(ctx, editor) {
		err = ErrAddressNotAuthorised()
		return
	}

	err = k.validateLength(ctx, body)
	if err != nil {
		return
	}

	account, ok := k.Account(ctx, id)
	if !ok {
		err = ErrUnknownAccount(id)
		return
	}

	account.Body = body
	k.setAccount(ctx, account)

	return
}

// Account gets a single account by its ID
func (k Keeper) Account(ctx sdk.Context, id uint64) (account Account, ok bool) {
	store := k.store(ctx)
	accountBytes := store.Get(key(id))
	if accountBytes == nil {
		return account, false
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(accountBytes, &account)

	return account, true
}

// Accounts gets all the accounts in reverse order
func (k Keeper) Accounts(ctx sdk.Context) (accounts Accounts) {
	store := k.store(ctx)
	iterator := sdk.KVStoreReversePrefixIterator(store, AccountsKeyPrefix)

	return k.iterate(iterator)
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


func (k Keeper) store(ctx sdk.Context) sdk.KVStore {
	return gaskv.NewStore(ctx.MultiStore().GetKVStore(k.storeKey), ctx.GasMeter(), app.KVGasConfig())
}

func logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}