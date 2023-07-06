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

func setupMsgExactTokensForTokens(t testing.TB) (types.MsgServer, keeper.Keeper, context.Context,
	*gomock.Controller, *testutil.MockBankKeeper) {
	ctrl := gomock.NewController(t)
	bankMock := testutil.NewMockBankKeeper(ctrl)
	k, ctx := keepertest.AmmKeeperWithMocks(t, bankMock)
	amm.InitGenesis(ctx, *k, *types.DefaultGenesis())
	server := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)

	bankMock.ExpectAny(context)

	createNPool(k, ctx, 2)

	return server, *k, context, ctrl, bankMock
}

func TestMsgSwapExactTokensForTokensNoKey(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgExactTokensForTokens(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	swapResponse, err := ms.SwapExactTokensForTokens(ctx, &types.MsgSwapExactTokensForTokens{
		Creator:      alice,
		PoolId:       1,
		AmountIn:     10,
		AmountOutMin: 10,
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

func TestMsgSwapExactTokensForTokensAfterDeadline(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgExactTokensForTokens(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	swapResponse, err := ms.SwapExactTokensForTokens(ctx, &types.MsgSwapExactTokensForTokens{
		Creator:      alice,
		PoolId:       0,
		AmountIn:     10,
		AmountOutMin: 10,
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

func TestMsgSwapExactTokensForTokensPathIncorrect(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgExactTokensForTokens(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	swapResponse, err := ms.SwapExactTokensForTokens(ctx, &types.MsgSwapExactTokensForTokens{
		Creator:      alice,
		PoolId:       0,
		AmountIn:     10,
		AmountOutMin: 10,
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

func TestMsgSwapExactTokensForTokensUnderMinAmount(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgExactTokensForTokens(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	swapResponse, err := ms.SwapExactTokensForTokens(ctx, &types.MsgSwapExactTokensForTokens{
		Creator:      alice,
		PoolId:       0,
		AmountIn:     1000,
		AmountOutMin: 100,
		Path: []string{
			"token",
			"foocoin",
		},
		To:       alice,
		Deadline: ctx.BlockTime().UTC().Format(types.DeadlineLayout),
	})

	require.Nil(t, swapResponse)
	require.Equal(t,
		"swaped value is under min amount",
		err.Error())
}

func TestMsgSwapExactTokensForTokens(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgExactTokensForTokens(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	swapResponse, _ := ms.SwapExactTokensForTokens(ctx, &types.MsgSwapExactTokensForTokens{
		Creator:      alice,
		PoolId:       0,
		AmountIn:     1000,
		AmountOutMin: 90,
		Path: []string{
			"token",
			"foocoin",
		},
		To:       alice,
		Deadline: ctx.BlockTime().UTC().Format(types.DeadlineLayout),
	})

	require.EqualValues(t, types.MsgSwapExactTokensForTokensResponse{
		AmountOut: 90,
	}, *swapResponse)
}
