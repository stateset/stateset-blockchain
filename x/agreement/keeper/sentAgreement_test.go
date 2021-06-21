package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stateset/stateset-blockchain/x/agreement/types"
	"github.com/stretchr/testify/assert"
)

func createNSentAgreement(keeper *Keeper, ctx sdk.Context, n int) []types.SentAgreement {
	items := make([]types.SentAgreement, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Id = keeper.AppendSentAgreement(ctx, items[i])
	}
	return items
}

func TestSentAgreementGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNSentAgreement(keeper, ctx, 10)
	for _, item := range items {
		assert.Equal(t, item, keeper.GetSentAgreement(ctx, item.Id))
	}
}

func TestSentAgreementExist(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNSentAgreement(keeper, ctx, 10)
	for _, item := range items {
		assert.True(t, keeper.HasSentAgreement(ctx, item.Id))
	}
}

func TestSentAgreementRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNSentAgreement(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveSentAgreement(ctx, item.Id)
		assert.False(t, keeper.HasSentAgreement(ctx, item.Id))
	}
}

func TestSentAgreementGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNSentAgreement(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllSentAgreement(ctx))
}

func TestSentAgreementCount(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNSentAgreement(keeper, ctx, 10)
	count := uint64(len(items))
	assert.Equal(t, count, keeper.GetSentAgreementCount(ctx))
}
