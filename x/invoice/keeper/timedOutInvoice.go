package keeper

import (
	"encoding/binary"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stateset/stateset-blockchain/x/purchaseorder/types"
)

// GetTimedoutInvoiceCount get the total number of timedoutInvoice
func (k Keeper) GetTimedoutInvoiceCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutInvoiceCountKey))
	byteKey := types.KeyPrefix(types.TimedoutInvoiceCountKey)
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

// SetTimedoutInvoiceCount set the total number of timedoutInvoice
func (k Keeper) SetTimedoutInvoiceCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutInvoiceCountKey))
	byteKey := types.KeyPrefix(types.TimedoutInvoiceCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendTimedoutInvoice appends a timedoutInvoice in the store with a new id and update the count
func (k Keeper) AppendTimedoutInvoice(
	ctx sdk.Context,
	timedoutInvoice types.TimedoutInvoice,
) uint64 {
	// Create the timedoutInvoice
	count := k.GetTimedoutInvoiceCount(ctx)

	// Set the ID of the appended value
	timedoutInvoice.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutInvoiceKey))
	appendedValue := k.cdc.MustMarshalBinaryBare(&timedoutInvoice)
	store.Set(GetTimedoutInvoiceIDBytes(timedoutInvoice.Id), appendedValue)

	// Update timedoutInvoice count
	k.SetTimedoutInvoiceCount(ctx, count+1)

	return count
}

// SetTimedoutInvoice set a specific timedoutInvoice in the store
func (k Keeper) SetTimedoutInvoice(ctx sdk.Context, timedoutInvoice types.TimedoutInvoice) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutInvoiceKey))
	b := k.cdc.MustMarshalBinaryBare(&timedoutInvoice)
	store.Set(GetTimedoutInvoiceIDBytes(timedoutInvoice.Id), b)
}

// GetTimedoutInvoice returns a timedoutInvoice from its id
func (k Keeper) GetTimedoutInvoice(ctx sdk.Context, id uint64) types.TimedoutInvoice {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutInvoiceKey))
	var timedoutInvoice types.TimedoutInvoice
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetTimedoutInvoiceIDBytes(id)), &timedoutInvoice)
	return timedoutInvoice
}

// HasTimedoutInvoice checks if the timedoutInvoice exists in the store
func (k Keeper) HasTimedoutInvoice(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutInvoiceKey))
	return store.Has(GetTimedoutInvoiceIDBytes(id))
}

// GetTimedoutInvoiceOwner returns the creator of the timedoutInvoice
func (k Keeper) GetTimedoutInvoiceOwner(ctx sdk.Context, id uint64) string {
	return k.GetTimedoutInvoice(ctx, id).Creator
}

// RemoveTimedoutInvoice removes a timedoutInvoice from the store
func (k Keeper) RemoveTimedoutInvoice(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutInvoiceKey))
	store.Delete(GetTimedoutInvoiceIDBytes(id))
}

// GetAllTimedoutInvoice returns all timedoutInvoice
func (k Keeper) GetAllTimedoutInvoice(ctx sdk.Context) (list []types.TimedoutInvoice) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutInvoiceKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TimedoutInvoice
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetTimedoutInvoiceIDBytes returns the byte representation of the ID
func GetTimedoutInvoiceIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetTimedoutInvoiceIDFromBytes returns ID in uint64 format from a byte array
func GetTimedoutInvoiceIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
