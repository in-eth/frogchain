package keeper_test

import (
	"testing"

	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgSwapExactTokensForTokens(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgManagePools(t, true)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()

	tests := []struct {
		desc             string
		msg              types.MsgSwapExactTokensForTokens
		expectedResponse types.MsgSwapExactTokensForTokensResponse
		err              error
	}{
		{
			desc: "invalid pool id",
			msg: types.MsgSwapExactTokensForTokens{
				Creator: alice,
				PoolId:  1,
			},
			err: sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key 1 doesn't exist"),
		}, {
			desc: "deadline passed",
			msg: types.MsgSwapExactTokensForTokens{
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
			msg: types.MsgSwapExactTokensForTokens{
				Creator:  alice,
				PoolId:   0,
				AmountIn: sdk.NewDec(1000),
				Path:     []string{"foocoin", "foocoin"},
				To:       alice,
				Deadline: ctx.BlockTime(),
			},
			err: types.ErrInvalidSwapDenom,
		}, {
			desc: "invalid path",
			msg: types.MsgSwapExactTokensForTokens{
				Creator:  alice,
				PoolId:   0,
				AmountIn: sdk.NewDec(1000),
				Path:     []string{"foocoin", "token1"},
				To:       alice,
				Deadline: ctx.BlockTime(),
			},
			err: types.ErrInvalidPath,
		}, {
			desc: "out amount is under min amount",
			msg: types.MsgSwapExactTokensForTokens{
				Creator:      alice,
				PoolId:       0,
				AmountIn:     sdk.NewDec(1000),
				AmountOutMin: sdk.NewDec(100),
				Path:         []string{"foocoin", "token"},
				To:           alice,
				Deadline:     ctx.BlockTime(),
			},
			err: types.ErrUnderMinAmount,
		}, {
			desc: "valid params",
			msg: types.MsgSwapExactTokensForTokens{
				Creator:      alice,
				PoolId:       0,
				AmountIn:     sdk.NewDec(10000),
				AmountOutMin: sdk.NewDec(10),
				Path:         []string{"foocoin", "token"},
				To:           alice,
				Deadline:     ctx.BlockTime(),
			},
			err:              nil,
			expectedResponse: types.MsgSwapExactTokensForTokensResponse{AmountOut: 70},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			response, err := ms.SwapExactTokensForTokens(ctx, &tt.msg)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				require.EqualValues(t, &tt.expectedResponse, response)
			}
		})
	}
}
