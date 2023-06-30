package keeper_test

import (
	"testing"

	keepertest "frogchain/testutil/keeper"
	"frogchain/testutil/nullify"
	"frogchain/x/amm/keeper"
	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func createNPool(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Pool {
	items := make([]types.Pool, n)
	for i := range items {
		items[i].Id = keeper.AppendPool(ctx, items[i])
	}
	return items
}

func TestPoolGet(t *testing.T) {
	keeper, ctx := keepertest.AmmKeeper(t)
	items := createNPool(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetPool(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestPoolRemove(t *testing.T) {
	keeper, ctx := keepertest.AmmKeeper(t)
	items := createNPool(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePool(ctx, item.Id)
		_, found := keeper.GetPool(ctx, item.Id)
		require.False(t, found)
	}
}

func TestPoolGetAll(t *testing.T) {
	keeper, ctx := keepertest.AmmKeeper(t)
	items := createNPool(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPool(ctx)),
	)
}

func TestPoolCount(t *testing.T) {
	keeper, ctx := keepertest.AmmKeeper(t)
	items := createNPool(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetPoolCount(ctx))
}
