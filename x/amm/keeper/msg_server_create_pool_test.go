package keeper_test

import (
	"context"
	"testing"

	keepertest "frogchain/testutil/keeper"
	"frogchain/x/amm"
	"frogchain/x/amm/keeper"
	"frogchain/x/amm/testutil"
	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func setupMsgCreatePool(t testing.TB) (types.MsgServer, keeper.Keeper, context.Context,
	*gomock.Controller, *testutil.MockBankKeeper) {
	ctrl := gomock.NewController(t)
	bankMock := testutil.NewMockBankKeeper(ctrl)
	k, ctx := keepertest.AmmKeeperWithMocks(t, bankMock)
	amm.InitGenesis(ctx, *k, *types.DefaultGenesis())
	server := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)

	bankMock.ExpectAny(context)
	return server, *k, context, ctrl, bankMock
}

func TestMsgCreatePoolNotEnoughCoin(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgCreatePool(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()

	createResponse, err := ms.CreatePool(ctx, &types.MsgCreatePool{
		Creator: alice,
		PoolParam: &types.PoolParam{
			SwapFee:      1,
			ExitFee:      1,
			FeeCollector: alice,
		},
		PoolAssets: []sdk.Coin{
			sdk.NewCoin(
				"token",
				sdk.NewInt(100),
			),
			sdk.NewCoin(
				"foocoin",
				sdk.NewInt(100),
			),
		},
		AssetWeights: []uint64{1, 1},
	})
	require.Nil(t, createResponse)
	require.Equal(t,
		"not enough coins for minimum liquidity: invalid amount",
		err.Error())
}

func TestMsgCreatePool(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgCreatePool(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()

	createResponse, _ := ms.CreatePool(ctx, &types.MsgCreatePool{
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

	require.EqualValues(t, types.MsgCreatePoolResponse{
		Id: 0,
	}, *createResponse)
}
