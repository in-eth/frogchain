package keeper_test

import (
	keepertest "frogchain/testutil/keeper"
	"frogchain/testutil/nullify"
	"frogchain/x/amm/types"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetSwapTokensForExactTokens(t *testing.T) {
	keeper, ctx := keepertest.AmmKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNPool(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetSwapTokensForExactTokensRequest
		response *types.QueryGetSwapTokensForExactTokensResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetSwapTokensForExactTokensRequest{
				PoolId:    msgs[0].Id,
				AmountOut: 5,
				Path: []string{
					"123",
					"124",
				},
			},
			response: &types.QueryGetSwapTokensForExactTokensResponse{AmountIn: 1},
		},
		{
			desc: "Second",
			request: &types.QueryGetSwapTokensForExactTokensRequest{
				PoolId:    msgs[1].Id,
				AmountOut: 5,
				Path: []string{
					"123",
					"124",
				},
			},
			response: &types.QueryGetSwapTokensForExactTokensResponse{AmountIn: 1},
		},
		{
			desc: "Invalidpatth",
			request: &types.QueryGetSwapTokensForExactTokensRequest{
				PoolId:    msgs[1].Id,
				AmountOut: 5,
				Path: []string{
					"123",
					"125",
				},
			},
			err: types.ErrInvalidPath,
		},
		{
			desc: "ExceedBalance",
			request: &types.QueryGetSwapTokensForExactTokensRequest{
				PoolId:    msgs[1].Id,
				AmountOut: 11,
				Path: []string{
					"124",
					"123",
				},
			},
			err: types.ErrInvalidSwapAmount,
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetSwapTokensForExactTokensRequest{
				PoolId: 10000,
			},
			err: sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key 10000 doesn't exist"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.GetSwapTokensForExactTokens(wctx, tc.request)
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
