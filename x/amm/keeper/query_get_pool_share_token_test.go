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

func TestGetPoolShareToken(t *testing.T) {
	keeper, ctx := keepertest.AmmKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNPool(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetPoolShareTokenRequest
		response *types.QueryGetPoolShareTokenResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetPoolShareTokenRequest{
				PoolId: msgs[0].Id,
			},
			response: &types.QueryGetPoolShareTokenResponse{ShareToken: msgs[0].ShareToken},
		},
		{
			desc: "Second",
			request: &types.QueryGetPoolShareTokenRequest{
				PoolId: msgs[1].Id,
			},
			response: &types.QueryGetPoolShareTokenResponse{ShareToken: msgs[1].ShareToken},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetPoolShareTokenRequest{
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
			response, err := keeper.GetPoolShareToken(wctx, tc.request)
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
