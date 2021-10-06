package keeper

import (
	"encoding/binary"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stateset/stateset-blockchain/x/invoice/types"
)

// GetSentInvoiceCount get the total number of sentInvoice
func (k Keeper) GetSentInvoiceCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentInvoiceCountKey))
	byteKey := types.KeyPrefix(types.SentInvoiceCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to iint64
		panic("cannot decode count")
	}

	return count
}

// SetSentInvoiceCount set the total number of sentInvoice
func (k Keeper) SetSentInvoiceCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentInvoiceCountKey))
	byteKey := types.KeyPrefix(types.SentInvoiceCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendSentInvoice appends a sentInvoice in the store with a new id and update the count
func (k Keeper) AppendSentInvoice(
	ctx sdk.Context,
	sentInvoice types.SentInvoice,
) uint64 {
	// Create the sentInvoice
	count := k.GetSentInvoiceCount(ctx)

	// Set the ID of the appended value
	sentInvoice.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentInvoiceKey))
	appendedValue := k.cdc.MustMarshalBinaryBare(&sentInvoice)
	store.Set(GetSentInvoiceIDBytes(sentInvoice.Id), appendedValue)

	// Update sentInvoice count
	k.SetSentInvoiceCount(ctx, count+1)

	return count
}

// SetSentInvoice set a specific sentInvoice in the store
func (k Keeper) SetSentInvoice(ctx sdk.Context, sentInvoice types.SentInvoice) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentInvoiceKey))
	b := k.cdc.MustMarshalBinaryBare(&sentInvoice)
	store.Set(GetSentInvoiceIDBytes(sentInvoice.Id), b)
}

// GetSentInvoice returns a sentInvoice from its id
func (k Keeper) GetSentInvoice(ctx sdk.Context, id uint64) types.SentInvoice {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentInvoiceKey))
	var sentInvoice types.SentInvoice
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetSentInvoiceIDBytes(id)), &sentInvoice)
	return sentInvoice
}

// HasSentInvoice checks if the sentInvoice exists in the store
func (k Keeper) HasSentInvoice(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentInvoiceKey))
	return store.Has(GetSentInvoiceIDBytes(id))
}

// GetSentInvoiceOwner returns the creator of the sentInvoice
func (k Keeper) GetSentInvoiceOwner(ctx sdk.Context, id uint64) string {
	return k.GetSentInvoice(ctx, id).Creator
}

// RemoveSentInvoice removes a sentInvoice from the store
func (k Keeper) RemoveSentInvoice(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentInvoiceKey))
	store.Delete(GetSentInvoiceIDBytes(id))
}

// GetAllSentInvoice returns all sentInvoice
func (k Keeper) GetAllSentInvoice(ctx sdk.Context) (list []types.SentInvoice) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentInvoiceKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.SentInvoice
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetSentInvoiceIDBytes returns the byte representation of the ID
func GetSentInvoiceIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetSentInvoiceIDFromBytes returns ID in uint64 format from a byte array
func GetSentInvoiceIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
