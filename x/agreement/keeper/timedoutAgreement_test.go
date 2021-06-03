package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stateset/stateset-blockchain/x/agreement/types"
	"github.com/stretchr/testify/assert"
)

func createNTimedoutAgreement(keeper *Keeper, ctx sdk.Context, n int) []types.TimedoutAgreement {
	items := make([]types.TimedoutAgreement, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Id = keeper.AppendTimedoutAgreement(ctx, items[i])
	}
	return items
}

func TestTimedoutAgreementGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTimedoutAgreement(keeper, ctx, 10)
	for _, item := range items {
		assert.Equal(t, item, keeper.GetTimedoutAgreement(ctx, item.Id))
	}
}

func TestTimedoutAgreementExist(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTimedoutAgreement(keeper, ctx, 10)
	for _, item := range items {
		assert.True(t, keeper.HasTimedoutAgreement(ctx, item.Id))
	}
}

func TestTimedoutAgreementRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTimedoutAgreement(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveTimedoutAgreement(ctx, item.Id)
		assert.False(t, keeper.HasTimedoutAgreement(ctx, item.Id))
	}
}

func TestTimedoutAgreementGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTimedoutAgreement(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllTimedoutAgreement(ctx))
}

func TestTimedoutAgreementCount(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTimedoutAgreement(keeper, ctx, 10)
	count := uint64(len(items))
	assert.Equal(t, count, keeper.GetTimedoutAgreementCount(ctx))
}
