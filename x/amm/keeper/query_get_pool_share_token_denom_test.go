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

func TestGetPoolShareTokenDenom(t *testing.T) {
	keeper, ctx := keepertest.AmmKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNPool(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetPoolShareTokenDenomRequest
		response *types.QueryGetPoolShareTokenDenomResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetPoolShareTokenDenomRequest{
				Id: msgs[0].Id,
			},
			response: &types.QueryGetPoolShareTokenDenomResponse{ShareDenom: msgs[0].ShareToken.Denom},
		},
		{
			desc: "Second",
			request: &types.QueryGetPoolShareTokenDenomRequest{
				Id: msgs[1].Id,
			},
			response: &types.QueryGetPoolShareTokenDenomResponse{ShareDenom: msgs[1].ShareToken.Denom},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetPoolShareTokenDenomRequest{
				Id: 10000,
			},
			err: sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key 10000 doesn't exist"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.GetPoolShareTokenDenom(wctx, tc.request)
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
