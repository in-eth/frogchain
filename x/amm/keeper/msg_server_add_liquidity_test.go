package keeper_test

import (
	"context"
	"testing"

	keepertest "frogchain/testutil/keeper"
	"frogchain/x/amm"
	"frogchain/x/amm/keeper"
	"frogchain/x/amm/types"

	"frogchain/x/amm/testutil"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func setupMsgAddLiquidity(t testing.TB) (types.MsgServer, keeper.Keeper, context.Context,
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
	return server, *k, context, ctrl, bankMock
}

func TestMsgAddLiquidityNoKey(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgAddLiquidity(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	createResponse, err := ms.AddLiquidity(ctx, &types.MsgAddLiquidity{
		Creator:        alice,
		PoolId:         2,
		DesiredAmounts: []uint64{10, 10},
		MinAmounts:     []uint64{10, 10},
	})
	require.Nil(t, createResponse)
	require.Equal(t,
		"key 2 doesn't exist: key not found",
		err.Error())
}

func TestMsgAddLiquidityNotCorrectAmountLength(t *testing.T) {
	// ms, keeper, context, ctrl, bank := setupMsgAddLiquidity(t)
	ms, _, context, ctrl, _ := setupMsgAddLiquidity(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()

	addResponse, err := ms.AddLiquidity(ctx, &types.MsgAddLiquidity{
		Creator:        alice,
		PoolId:         0,
		DesiredAmounts: []uint64{10, 10, 10},
		MinAmounts:     []uint64{10, 10, 10},
	})
	require.Nil(t, addResponse)
	require.Equal(t,
		"invalid assets length",
		err.Error())
}

func TestMsgAddLiquidityNoLiquidity(t *testing.T) {
	// ms, keeper, context, ctrl, bank := setupMsgAddLiquidity(t)
	ms, _, context, ctrl, _ := setupMsgAddLiquidity(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()

	addResponse, err := ms.AddLiquidity(ctx, &types.MsgAddLiquidity{
		Creator:        alice,
		PoolId:         0,
		DesiredAmounts: []uint64{0, 10},
		MinAmounts:     []uint64{10, 10},
	})
	require.Nil(t, addResponse)
	require.Equal(t,
		"no liquidity with amounts you deposit: invalid amount",
		err.Error())
}

func TestMsgAddLiquidityMinErr(t *testing.T) {
	// ms, keeper, context, ctrl, bank := setupMsgAddLiquidity(t)
	ms, _, context, ctrl, _ := setupMsgAddLiquidity(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()

	addResponse, err := ms.AddLiquidity(ctx, &types.MsgAddLiquidity{
		Creator:        alice,
		PoolId:         0,
		DesiredAmounts: []uint64{10, 10},
		MinAmounts:     []uint64{11, 10},
	})
	require.Nil(t, addResponse)
	require.Equal(t,
		"calculated amount is below minimum, 0, 10, 11: invalid amount",
		err.Error())
}

func TestMsgAddLiquidity(t *testing.T) {
	// ms, keeper, context, ctrl, bank := setupMsgAddLiquidity(t)
	ms, _, context, ctrl, _ := setupMsgAddLiquidity(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()

	addResponse, _ := ms.AddLiquidity(ctx, &types.MsgAddLiquidity{
		Creator:        alice,
		PoolId:         0,
		DesiredAmounts: []uint64{10, 10},
		MinAmounts:     []uint64{10, 10},
	})

	shareToken := sdk.NewCoin(types.ShareTokenIndex(0), math.NewInt(10))

	response := &types.MsgAddLiquidityResponse{
		ShareToken: shareToken,
	}

	require.EqualValues(t, response, addResponse)
}
