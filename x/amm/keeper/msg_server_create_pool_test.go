package keeper_test

import (
	"testing"

	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgCreatePool(t *testing.T) {
	ms, _, context, ctrl, _ := setupMsgManagePools(t, false)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()

	tests := []struct {
		desc             string
		msg              types.MsgCreatePool
		expectedResponse types.MsgCreatePoolResponse
		err              error
	}{
		{
			desc: "invalid swap fee amount",
			msg: types.MsgCreatePool{
				Creator: alice,
				PoolParam: &types.PoolParam{
					SwapFee:      sdk.NewDec(1000000000),
					ExitFee:      sdk.NewDec(1),
					FeeCollector: alice,
				},
			},
			err: types.ErrFeeOverflow,
		}, {
			desc: "invalid exit fee amount",
			msg: types.MsgCreatePool{
				Creator: alice,
				PoolParam: &types.PoolParam{
					SwapFee:      sdk.NewDec(1),
					ExitFee:      sdk.NewDec(10000000000),
					FeeCollector: alice,
				},
			},
			err: types.ErrFeeOverflow,
		}, {
			desc: "invalid assets length",
			msg: types.MsgCreatePool{
				Creator: alice,
				PoolParam: &types.PoolParam{
					SwapFee:      sdk.NewDec(1),
					ExitFee:      sdk.NewDec(1),
					FeeCollector: alice,
				},
				PoolAssets: []sdk.Coin{
					sdk.NewCoin(
						"token",
						sdk.NewInt(100),
					),
				},
			},
			err: types.ErrInvalidAssetsLength,
		}, {
			desc: "invalid weights length",
			msg: types.MsgCreatePool{
				Creator: alice,
				PoolParam: &types.PoolParam{
					SwapFee:      sdk.NewDec(1),
					ExitFee:      sdk.NewDec(1),
					FeeCollector: alice,
				},
				PoolAssets: []sdk.Coin{
					sdk.NewCoin(
						"token",
						sdk.NewInt(100),
					),
					sdk.NewCoin(
						"foocoin",
						sdk.NewInt(100),
					),
				},
				AssetWeights: []sdk.Dec{sdk.NewDec(10)},
			},
			err: types.ErrInvalidWeightlength,
		}, {
			desc: "mint share token amount is under minimum liquidity",
			msg: types.MsgCreatePool{
				Creator: alice,
				PoolParam: &types.PoolParam{
					SwapFee:      sdk.NewDec(1),
					ExitFee:      sdk.NewDec(1),
					FeeCollector: alice,
				},
				PoolAssets: []sdk.Coin{
					sdk.NewCoin(
						"token",
						sdk.NewInt(100),
					),
					sdk.NewCoin(
						"foocoin",
						sdk.NewInt(100),
					),
				},
				AssetWeights: []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10)},
			},
			err: sdkerrors.Wrapf(types.ErrInvalidAmount,
				"not enough coins for minimum liquidity",
			),
		}, {
			desc: "valid params",
			msg: types.MsgCreatePool{
				Creator: alice,
				PoolParam: &types.PoolParam{
					SwapFee:      sdk.NewDec(1),
					ExitFee:      sdk.NewDec(1),
					FeeCollector: alice,
				},
				PoolAssets: []sdk.Coin{
					sdk.NewCoin(
						"token",
						sdk.NewInt(10000),
					),
					sdk.NewCoin(
						"foocoin",
						sdk.NewInt(1000),
					),
				},
				AssetWeights: []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10)},
			},
			expectedResponse: types.MsgCreatePoolResponse{
				Id: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			response, err := ms.CreatePool(ctx, &tt.msg)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				require.EqualValues(t, &tt.expectedResponse, response)
			}
		})
	}
}
