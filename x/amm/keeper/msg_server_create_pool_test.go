package keeper_test

import (
	"context"
	"testing"

	keepertest "frogchain/testutil/keeper"
	"frogchain/x/amm"
	"frogchain/x/amm/keeper"
	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func setupMsgCreatePool(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.AmmKeeper(t)
	amm.InitGenesis(ctx, *k, *types.DefaultGenesis())
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

func TestMsgCreatePoolNotEnoughCoin(t *testing.T) {
	ms, ctx := setupMsgCreatePool(t)

	createResponse, err := ms.CreatePool(ctx, &types.MsgCreatePool{
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
	require.Nil(t, createResponse)
	require.Equal(t,
		"123123123",
		err.Error())
}

func TestMsgCreatePool(t *testing.T) {
	ms, ctx := setupMsgCreatePool(t)

	createResponse, _ := ms.CreatePool(ctx, &types.MsgCreatePool{
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
	require.Nil(t, createResponse)
	require.EqualValues(t, types.MsgCreatePoolResponse{
		Id: 1,
	}, *createResponse)
}
