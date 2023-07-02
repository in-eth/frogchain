package keeper_test

import (
	"context"
	"testing"

	keepertest "frogchain/testutil/keeper"
	"frogchain/x/amm"
	"frogchain/x/amm/keeper"
	"frogchain/x/amm/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func setupMsgAddLiquidity(t testing.TB) (types.MsgServer, keeper.Keeper, context.Context) {
	// *gomock.Controller, *testutil.MockBankKeeper) {
	// ctrl := gomock.NewController(t)
	// bankMock := testutil.NewMockBankKeeper(ctrl)
	k, ctx := keepertest.AmmKeeperWithMocks(t, nil)
	amm.InitGenesis(ctx, *k, *types.DefaultGenesis())
	server := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)

	// bankMock.EXPECT().SendCoinsFromAccountToModule(ctx, alice, types.ModuleName, sdk.NewCoins(
	// 	sdk.NewCoin(
	// 		"token",
	// 		sdk.NewInt(10),
	// 	),
	// 	sdk.NewCoin(
	// 		"foocoin",
	// 		sdk.NewInt(10),
	// 	),
	// ))

	// bankMock.EXPECT().SendCoinsFromModuleToAccount(ctx, types.ModuleName, alice, sdk.NewCoins(
	// 	sdk.NewCoin(
	// 		types.ShareTokenIndex(1),
	// 		sdk.NewInt(10),
	// 	),
	// ))

	server.CreatePool(context, &types.MsgCreatePool{
		Creator: alice,
		PoolParam: &types.PoolParam{
			SwapFee:      1,
			ExitFee:      1,
			FeeCollector: alice,
		},
		PoolAssets: []types.PoolToken{
			types.PoolToken{
				TokenDenom:   "token",
				TokenWeight:  1,
				TokenReserve: 0,
			},
			types.PoolToken{
				TokenDenom:   "foocoin",
				TokenWeight:  1,
				TokenReserve: 0,
			},
		},
		AssetAmounts: []uint64{10, 10},
	})
	return server, *k, context //ctrl, bankMock
}

func TestMsgAddLiquidityNoKey(t *testing.T) {
	// ms, _, context, ctrl, _ := setupMsgAddLiquidity(t)
	ms, _, context := setupMsgAddLiquidity(t)
	ctx := sdk.UnwrapSDKContext(context)
	// defer ctrl.Finish()
	createResponse, err := ms.AddLiquidity(ctx, &types.MsgAddLiquidity{
		Creator:        alice,
		PoolId:         2,
		DesiredAmounts: []uint64{10, 10},
		MinAmounts:     []uint64{10, 10},
	})
	require.Nil(t, createResponse)
	require.Equal(t,
		"key 1 doesn't exist: key not found",
		err.Error())
}

func TestMsgAddLiquidity(t *testing.T) {
	// ms, keeper, context, ctrl, bank := setupMsgAddLiquidity(t)
	// ms, _, context, ctrl, _ := setupMsgAddLiquidity(t)
	ms, _, context := setupMsgAddLiquidity(t)
	ctx := sdk.UnwrapSDKContext(context)
	// defer ctrl.Finish()

	addResponse, _ := ms.AddLiquidity(ctx, &types.MsgAddLiquidity{
		Creator:        alice,
		PoolId:         1,
		DesiredAmounts: []uint64{10, 10},
		MinAmounts:     []uint64{10, 10},
	})
	require.Nil(t, addResponse)

	shareToken := sdk.NewCoin(types.ShareTokenIndex(1), math.NewInt(10))

	require.EqualValues(t, types.MsgAddLiquidityResponse{
		ShareToken: shareToken,
	}, *addResponse)
}
