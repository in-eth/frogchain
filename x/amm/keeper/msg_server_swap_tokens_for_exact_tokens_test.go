package keeper_test

import (
	"context"
	"testing"

	keepertest "frogchain/testutil/keeper"
	"frogchain/x/amm"
	"frogchain/x/amm/keeper"
	"frogchain/x/amm/types"

	"frogchain/x/amm/testutil"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func setupMsgTokensForExactTokens(t testing.TB) (types.MsgServer, keeper.Keeper, context.Context,
	*gomock.Controller, *testutil.MockBankKeeper) {
	ctrl := gomock.NewController(t)
	bankMock := testutil.NewMockBankKeeper(ctrl)
	k, ctx := keepertest.AmmKeeperWithMocks(t, bankMock)
	amm.InitGenesis(ctx, *k, *types.DefaultGenesis())
	server := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)

	bankMock.ExpectAny(context)

	server.CreatePool(context, &types.MsgCreatePool{
		Creator: alice,
		PoolParam: &types.PoolParam{
			SwapFee:      1,
			ExitFee:      1,
			FeeCollector: alice,
		},
		PoolAssets: []sdk.Coin{
			sdk.NewCoin(
				"token",
				sdk.NewInt(10000),
			),
			sdk.NewCoin(
				"foocoin",
				sdk.NewInt(10000),
			),
		},
		AssetWeights: []uint64{10, 10},
	})

	return server, *k, context, ctrl, bankMock
}

func TestMsgSwapTokensForExactTokensNoKey(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgTokensForExactTokens(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	swapResponse, err := ms.SwapTokensForExactTokens(ctx, &types.MsgSwapTokensForExactTokens{
		Creator:   alice,
		PoolId:    1,
		AmountOut: 10,
		Path: []string{
			"token",
			"foocoin",
		},
		To:       alice,
		Deadline: ctx.BlockTime().UTC().Format(types.DeadlineLayout),
	})

	require.Nil(t, swapResponse)
	require.Equal(t,
		"key 1 doesn't exist: key not found",
		err.Error())
}

func TestMsgSwapTokensForExactTokensAfterDeadline(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgTokensForExactTokens(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	swapResponse, err := ms.SwapTokensForExactTokens(ctx, &types.MsgSwapTokensForExactTokens{
		Creator:   alice,
		PoolId:    0,
		AmountOut: 10,
		Path: []string{
			"token",
			"foocoin",
		},
		To:       alice,
		Deadline: ctx.BlockTime().Add(-100).UTC().Format(types.DeadlineLayout),
	})

	require.Nil(t, swapResponse)
	require.Equal(t,
		"0000-12-31 23:59:59.9999999 +0000 UTC: deadline is passed",
		err.Error())
}

func TestMsgSwapTokensForExactTokensPathIncorrect(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgTokensForExactTokens(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	swapResponse, err := ms.SwapTokensForExactTokens(ctx, &types.MsgSwapTokensForExactTokens{
		Creator:   alice,
		PoolId:    0,
		AmountOut: 10,
		Path: []string{
			"token1",
			"foocoin",
		},
		To:       alice,
		Deadline: ctx.BlockTime().UTC().Format(types.DeadlineLayout),
	})

	require.Nil(t, swapResponse)
	require.Equal(t,
		"invalid path for swap",
		err.Error())
}

func TestMsgSwapTokensForExactTokens(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgTokensForExactTokens(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	swapResponse, _ := ms.SwapTokensForExactTokens(ctx, &types.MsgSwapTokensForExactTokens{
		Creator:   alice,
		PoolId:    0,
		AmountOut: 20,
		Path: []string{
			"token",
			"foocoin",
		},
		To:       alice,
		Deadline: ctx.BlockTime().UTC().Format(types.DeadlineLayout),
	})

	require.EqualValues(t, types.MsgSwapTokensForExactTokensResponse{
		AmountIn: 20,
	}, *swapResponse)
}
