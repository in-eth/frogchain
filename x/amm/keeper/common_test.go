package keeper_test

import (
	"context"
	keepertest "frogchain/testutil/keeper"
	"frogchain/x/amm"
	"frogchain/x/amm/keeper"
	"frogchain/x/amm/testutil"
	"frogchain/x/amm/types"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/golang/mock/gomock"
)

const (
	alice = testutil.Alice
	bob   = testutil.Bob
	// carol = testutil.Carol
)

func createPool(ms types.MsgServer, ctx context.Context) {
	ms.CreatePool(ctx, &types.MsgCreatePool{
		Creator: alice,
		PoolParam: &types.PoolParam{
			SwapFee:      sdk.NewDec(1),
			ExitFee:      sdk.NewDec(1),
			FeeCollector: alice,
		},
		PoolAssets: []sdk.Coin{
			sdk.NewCoin(
				"foocoin",
				sdk.NewInt(30000),
			),
			sdk.NewCoin(
				"token",
				sdk.NewInt(300),
			),
		},
		AssetWeights: []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10)},
	})
}

func ErrorWrap(err error, format string, args ...interface{}) error {
	return sdkerrors.Wrapf(
		err,
		format,
		args,
	)
}

func setupMsgManagePools(t testing.TB, isPoolCreated bool) (types.MsgServer, keeper.Keeper, context.Context,
	*gomock.Controller, *testutil.MockBankKeeper) {
	ctrl := gomock.NewController(t)
	bankMock := testutil.NewMockBankKeeper(ctrl)
	k, ctx := keepertest.AmmKeeperWithMocks(t, bankMock)
	amm.InitGenesis(ctx, *k, *types.DefaultGenesis())
	server := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)

	bankMock.ExpectAny(context)

	if isPoolCreated {
		createPool(server, context)
	}

	return server, *k, context, ctrl, bankMock
}
