package keeper_test

import (
	keepertest "frogchain/testutil/keeper"
	"frogchain/testutil/nullify"
	"frogchain/x/amm/types"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetPoolTokenDenom(t *testing.T) {
	keeper, ctx := keepertest.AmmKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNPool(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetPoolTokenDenomRequest
		response *types.QueryGetPoolTokenDenomResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetPoolTokenDenomRequest{
				Id:      msgs[0].Id,
				AssetId: 0,
			},
			response: &types.QueryGetPoolTokenDenomResponse{TokenDenom: msgs[0].PoolAssets[0].TokenDenom},
		},
		{
			desc: "Second",
			request: &types.QueryGetPoolTokenDenomRequest{
				Id:      msgs[1].Id,
				AssetId: 1,
			},
			response: &types.QueryGetPoolTokenDenomResponse{TokenDenom: msgs[1].PoolAssets[1].TokenDenom},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetPoolTokenDenomRequest{
				Id:      0,
				AssetId: 10000,
			},
			err: types.ErrInvalidLength,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.GetPoolTokenDenom(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}
