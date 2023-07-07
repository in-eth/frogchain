package keeper_test

import (
	"testing"

	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgRemoveLiquidity(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgManagePools(t, true)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()

	tests := []struct {
		desc             string
		msg              types.MsgRemoveLiquidity
		expectedResponse types.MsgRemoveLiquidityResponse
		err              error
	}{
		{
			desc: "invalid pool id",
			msg: types.MsgRemoveLiquidity{
				Creator:    alice,
				PoolId:     1,
				Liquidity:  sdk.NewDec(10),
				MinAmounts: []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10)},
			},
			err: sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key 1 doesn't exist"),
		}, {
			desc: "invalid min amounts length",
			msg: types.MsgRemoveLiquidity{
				Creator:    alice,
				PoolId:     0,
				Liquidity:  sdk.NewDec(10),
				MinAmounts: []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10), sdk.NewDec(10)},
			},
			err: types.ErrInvalidLength,
		}, {
			desc: "calculated amount is under min amount",
			msg: types.MsgRemoveLiquidity{
				Creator:    alice,
				PoolId:     0,
				Liquidity:  sdk.NewDec(100),
				MinAmounts: []sdk.Dec{sdk.NewDec(1001), sdk.NewDec(3)},
			},
			err: sdkerrors.Wrapf(types.ErrInvalidAmount, "calculated amount is below minimum, 0, 1000, 1001"),
		}, {
			desc: "valid params",
			msg: types.MsgRemoveLiquidity{
				Creator:    alice,
				PoolId:     0,
				Liquidity:  sdk.NewDec(100),
				MinAmounts: []sdk.Dec{sdk.NewDec(100), sdk.NewDec(1)},
			},
			err: nil,
			expectedResponse: types.MsgRemoveLiquidityResponse{
				ReceivedTokens: []sdk.Coin{
					sdk.NewCoin("foocoin", sdk.NewInt(999)),
					sdk.NewCoin("token", sdk.NewInt(9)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			response, err := ms.RemoveLiquidity(ctx, &tt.msg)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				require.EqualValues(t, &tt.expectedResponse, response)
			}
		})
	}
}
