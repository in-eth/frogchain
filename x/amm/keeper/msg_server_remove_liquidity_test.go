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

func setupMsgRemoveLiquidity(t testing.TB) (types.MsgServer, keeper.Keeper, context.Context,
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
		AssetAmounts: []uint64{10000, 10000},
	})

	server.AddLiquidity(ctx, &types.MsgAddLiquidity{
		Creator:        bob,
		PoolId:         0,
		DesiredAmounts: []uint64{10, 10},
		MinAmounts:     []uint64{10, 10},
	})
	return server, *k, context, ctrl, bankMock
}

func TestMsgRemoveLiquidityNoKey(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgAddLiquidity(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	removeResponse, err := ms.RemoveLiquidity(ctx, &types.MsgRemoveLiquidity{
		Creator:    alice,
		PoolId:     1,
		Liquidity:  10,
		MinAmounts: []uint64{10, 10},
	})

	require.Nil(t, removeResponse)
	require.Equal(t,
		"key 1 doesn't exist: key not found",
		err.Error())
}

func TestMsgRemoveLiquidity(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgAddLiquidity(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	removeResponse, _ := ms.RemoveLiquidity(ctx, &types.MsgRemoveLiquidity{
		Creator:    alice,
		PoolId:     0,
		Liquidity:  10,
		MinAmounts: []uint64{10, 10},
	})

	response := &types.MsgRemoveLiquidityResponse{
		ReceivedTokens: []sdk.Coin{
			sdk.NewCoin("token", sdk.NewInt(10)),
			sdk.NewCoin("foocoin", sdk.NewInt(10)),
		},
	}

	require.EqualValues(t, response, removeResponse)
}
