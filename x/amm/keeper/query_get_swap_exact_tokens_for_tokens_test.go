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

func TestGetSwapExactTokensForTokens(t *testing.T) {
	keeper, ctx := keepertest.AmmKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNPool(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetSwapExactTokensForTokensRequest
		response *types.QueryGetSwapExactTokensForTokensResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetSwapExactTokensForTokensRequest{
				PoolId:   msgs[0].Id,
				AmountIn: 10,
				Path: []string{
					"foocoin",
					"token",
				},
			},
			response: &types.QueryGetSwapExactTokensForTokensResponse{AmountOut: 25},
		},
		{
			desc: "Second",
			request: &types.QueryGetSwapExactTokensForTokensRequest{
				PoolId:   msgs[1].Id,
				AmountIn: 10,
				Path: []string{
					"foocoin",
					"token",
				},
			},
			response: &types.QueryGetSwapExactTokensForTokensResponse{AmountOut: 25},
		},
		{
			desc: "Invalidpatth",
			request: &types.QueryGetSwapExactTokensForTokensRequest{
				PoolId:   msgs[1].Id,
				AmountIn: 10,
				Path: []string{
					"foocoin",
					"token1",
				},
			},
			err: types.ErrInvalidPath,
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetSwapExactTokensForTokensRequest{
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
			response, err := keeper.GetSwapExactTokensForTokens(wctx, tc.request)
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
