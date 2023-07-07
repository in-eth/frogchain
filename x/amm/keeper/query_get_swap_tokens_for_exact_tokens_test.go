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
				AmountOut: sdk.NewDec(5),
				Path: []string{
					"foocoin",
					"token",
				},
			},
			response: &types.QueryGetSwapTokensForExactTokensResponse{AmountIn: 5},
		},
		{
			desc: "Second",
			request: &types.QueryGetSwapTokensForExactTokensRequest{
				PoolId:    msgs[1].Id,
				AmountOut: sdk.NewDec(5),
				Path: []string{
					"foocoin",
					"token",
				},
			},
			response: &types.QueryGetSwapTokensForExactTokensResponse{AmountIn: 5},
		},
		{
			desc: "Invalidpatth",
			request: &types.QueryGetSwapTokensForExactTokensRequest{
				PoolId:    msgs[1].Id,
				AmountOut: sdk.NewDec(5),
				Path: []string{
					"foocoin",
					"125",
				},
			},
			err: types.ErrInvalidPath,
		},
		{
			desc: "ExceedBalance",
			request: &types.QueryGetSwapTokensForExactTokensRequest{
				PoolId:    msgs[1].Id,
				AmountOut: sdk.NewDec(100000),
				Path: []string{
					"token",
					"foocoin",
				},
			},
			err: types.ErrInvalidSwapAmount,
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetSwapTokensForExactTokensRequest{
				PoolId: 10000,
			},
			err: ErrorWrap(sdkerrors.ErrKeyNotFound, "key 10000 doesn't exist"),
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
