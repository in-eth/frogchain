package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "frogchain/testutil/keeper"
	"frogchain/testutil/nullify"
	"frogchain/x/investibc/keeper"
	"frogchain/x/investibc/types"
)

func createTestDepositDenom(keeper *keeper.Keeper, ctx sdk.Context) types.DepositDenom {
	item := types.DepositDenom{}
	keeper.SetDepositDenomStore(ctx, item)
	return item
}

func TestDepositDenomGet(t *testing.T) {
	keeper, ctx := keepertest.InvestibcKeeper(t)
	item := createTestDepositDenom(keeper, ctx)
	rst, found := keeper.GetDepositDenom(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestDepositDenomRemove(t *testing.T) {
	keeper, ctx := keepertest.InvestibcKeeper(t)
	createTestDepositDenom(keeper, ctx)
	keeper.RemoveDepositDenom(ctx)
	_, found := keeper.GetDepositDenom(ctx)
	require.False(t, found)
}
