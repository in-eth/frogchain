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

func TestGetPoolAsset(t *testing.T) {
	keeper, ctx := keepertest.AmmKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNPool(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetPoolAssetRequest
		response *types.QueryGetPoolAssetResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetPoolAssetRequest{
				PoolId:  msgs[0].Id,
				AssetId: 0,
			},
			response: &types.QueryGetPoolAssetResponse{PoolAsset: msgs[0].PoolAssets[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetPoolAssetRequest{
				PoolId:  msgs[1].Id,
				AssetId: 1,
			},
			response: &types.QueryGetPoolAssetResponse{PoolAsset: msgs[1].PoolAssets[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetPoolAssetRequest{
				PoolId:  0,
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
			response, err := keeper.GetPoolAsset(wctx, tc.request)
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
