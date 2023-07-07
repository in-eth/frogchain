package keeper_test

import (
	"testing"

	"frogchain/x/amm/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgAddLiquidity(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgManagePools(t, true)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()

	tests := []struct {
		desc             string
		msg              types.MsgAddLiquidity
		expectedResponse types.MsgAddLiquidityResponse
		err              error
	}{
		{
			desc: "invalid pool id",
			msg: types.MsgAddLiquidity{
				Creator:        alice,
				PoolId:         1,
				DesiredAmounts: []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10)},
				MinAmounts:     []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10)},
			},
			err: sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key 1 doesn't exist"),
		}, {
			desc: "invalid desired amounts length",
			msg: types.MsgAddLiquidity{
				Creator:        alice,
				PoolId:         0,
				DesiredAmounts: []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10), sdk.NewDec(10)},
				MinAmounts:     []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10)},
			},
			err: types.ErrInvalidLength,
		}, {
			desc: "invalid min amounts length",
			msg: types.MsgAddLiquidity{
				Creator:        alice,
				PoolId:         0,
				DesiredAmounts: []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10)},
				MinAmounts:     []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10), sdk.NewDec(10)},
			},
			err: types.ErrInvalidLength,
		}, {
			desc: "invalid liquidity amount",
			msg: types.MsgAddLiquidity{
				Creator:        alice,
				PoolId:         0,
				DesiredAmounts: []sdk.Dec{sdk.NewDec(0), sdk.NewDec(10)},
				MinAmounts:     []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10)},
			},
			err: sdkerrors.Wrapf(types.ErrInvalidAmount, "no liquidity with amounts you deposit"),
		}, {
			desc: "calculated amount is under min amount",
			msg: types.MsgAddLiquidity{
				Creator:        alice,
				PoolId:         0,
				DesiredAmounts: []sdk.Dec{sdk.NewDec(5), sdk.NewDec(10)},
				MinAmounts:     []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10)},
			},
			err: sdkerrors.Wrapf(types.ErrInvalidAmount, "calculated amount is below min amount, 0, 5, 10"),
		}, {
			desc: "valid params",
			msg: types.MsgAddLiquidity{
				Creator:        alice,
				PoolId:         0,
				DesiredAmounts: []sdk.Dec{sdk.NewDec(100), sdk.NewDec(10)},
				MinAmounts:     []sdk.Dec{sdk.NewDec(10), sdk.NewDec(1)},
			},
			err: nil,
			expectedResponse: types.MsgAddLiquidityResponse{
				ShareToken: sdk.NewCoin(types.ShareTokenIndex(0), math.NewInt(10)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			response, err := ms.AddLiquidity(ctx, &tt.msg)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				require.EqualValues(t, &tt.expectedResponse, response)
			}
		})
	}
}
