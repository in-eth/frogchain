package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "frogchain/testutil/keeper"
	"frogchain/testutil/nullify"
	"frogchain/x/investibc/types"
)

func TestDepositDenomQuery(t *testing.T) {
	keeper, ctx := keepertest.InvestibcKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	item := createTestDepositDenom(keeper, ctx)
	tests := []struct {
		desc     string
		request  *types.QueryGetDepositDenomRequest
		response *types.QueryGetDepositDenomResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetDepositDenomRequest{},
			response: &types.QueryGetDepositDenomResponse{DepositDenom: item},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.DepositDenom(wctx, tc.request)
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
