package keeper_test

import (
	"testing"

	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgSwapTokensForExactTokens(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgManagePools(t, true)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()

	tests := []struct {
		desc             string
		msg              types.MsgSwapTokensForExactTokens
		expectedResponse types.MsgSwapTokensForExactTokensResponse
		err              error
	}{
		{
			desc: "invalid pool id",
			msg: types.MsgSwapTokensForExactTokens{
				Creator: alice,
				PoolId:  1,
			},
			err: sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key 1 doesn't exist"),
		}, {
			desc: "deadline passed",
			msg: types.MsgSwapTokensForExactTokens{
				Creator:  alice,
				PoolId:   0,
				Deadline: ctx.BlockTime().Add(-100),
			},
			err: sdkerrors.Wrapf(
				types.ErrDeadlinePassed,
				ctx.BlockTime().Add(-100).String(),
			),
		}, {
			desc: "in and out token denom is equal",
			msg: types.MsgSwapTokensForExactTokens{
				Creator:   alice,
				PoolId:    0,
				AmountOut: sdk.NewDec(10),
				Path:      []string{"foocoin", "foocoin"},
				To:        alice,
				Deadline:  ctx.BlockTime(),
			},
			err: types.ErrInvalidSwapDenom,
		}, {
			desc: "invalid path",
			msg: types.MsgSwapTokensForExactTokens{
				Creator:   alice,
				PoolId:    0,
				AmountOut: sdk.NewDec(10),
				Path:      []string{"foocoin", "token1"},
				To:        alice,
				Deadline:  ctx.BlockTime(),
			},
			err: types.ErrInvalidPath,
		}, {
			desc: "out amount exceed pool balance",
			msg: types.MsgSwapTokensForExactTokens{
				Creator:   alice,
				PoolId:    0,
				AmountOut: sdk.NewDec(1000),
				Path:      []string{"foocoin", "token"},
				To:        alice,
				Deadline:  ctx.BlockTime(),
			},
			err: types.ErrInvalidSwapAmount,
		}, {
			desc: "valid params",
			msg: types.MsgSwapTokensForExactTokens{
				Creator:   alice,
				PoolId:    0,
				AmountOut: sdk.NewDec(100),
				Path:      []string{"foocoin", "token"},
				To:        alice,
				Deadline:  ctx.BlockTime(),
			},
			err:              nil,
			expectedResponse: types.MsgSwapTokensForExactTokensResponse{AmountIn: 15000},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			response, err := ms.SwapTokensForExactTokens(ctx, &tt.msg)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				require.EqualValues(t, &tt.expectedResponse, response)
			}
		})
	}
}
