package keeper_test

import (
	"strconv"
	"testing"

	keepertest "frogchain/testutil/keeper"
	"frogchain/testutil/nullify"
	"frogchain/x/investibc/keeper"
	"frogchain/x/investibc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNDepositBalance(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.DepositBalance {
	items := make([]types.DepositBalance, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)

		keeper.SetDepositBalance(ctx, items[i])
	}
	return items
}

func TestDepositBalanceGet(t *testing.T) {
	keeper, ctx := keepertest.InvestibcKeeper(t)
	items := createNDepositBalance(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetDepositBalance(ctx,
			item.Index,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestDepositBalanceRemove(t *testing.T) {
	keeper, ctx := keepertest.InvestibcKeeper(t)
	items := createNDepositBalance(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveDepositBalance(ctx,
			item.Index,
		)
		_, found := keeper.GetDepositBalance(ctx,
			item.Index,
		)
		require.False(t, found)
	}
}

func TestDepositBalanceGetAll(t *testing.T) {
	keeper, ctx := keepertest.InvestibcKeeper(t)
	items := createNDepositBalance(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllDepositBalance(ctx)),
	)
}
