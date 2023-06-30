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

func setupMsgAddLiquidity(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.AmmKeeper(t)
	amm.InitGenesis(ctx, *k, *types.DefaultGenesis())
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

func TestMsgAddLiquidity(t *testing.T) {
	ms, ctx := setupMsgAddLiquidity(t)
	createResponse, err := ms.AddLiquidity(ctx, &types.MsgAddLiquidity{
		Creator:        Alice,
		PoolId:         1,
		DesiredAmounts: []uint64{10, 10},
		MinAmounts:     []uint64{10, 10},
	})
	require.Nil(t, createResponse)
	require.Equal(t,
		"key 1 doesn't exist: key not found",
		err.Error())

	// shareToken := sdk.NewCoin(types.ShareTokenIndex(1), math.NewInt(100))

	// require.EqualValues(t, types.MsgAddLiquidityResponse{
	// 	ShareToken: &shareToken,
	// }, *createResponse)
}
