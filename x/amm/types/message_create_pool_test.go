package types

import (
	"frogchain/testutil/sample"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMsgCreatePool_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreatePool
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreatePool{
				Creator: "invalid_address",
			},
			err: ErrInvalidAddress,
		}, {
			name: "fee collector invalid address",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      sdk.NewDec(10),
					ExitFee:      sdk.NewDec(10),
					FeeCollector: "invalid_address",
				},
			},
			err: ErrInvalidAddress,
		}, {
			name: "exit fee amount overflow",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      sdk.NewDec(10),
					ExitFee:      sdk.NewDec(100000001),
					FeeCollector: sample.AccAddress(),
				},
			},
			err: ErrFeeOverflow,
		}, {
			name: "swap fee amount overflow",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      sdk.NewDec(100000001),
					ExitFee:      sdk.NewDec(10),
					FeeCollector: sample.AccAddress(),
				},
			},
			err: ErrFeeOverflow,
		}, {
			name: "invalid weights length",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      sdk.NewDec(10),
					ExitFee:      sdk.NewDec(10),
					FeeCollector: sample.AccAddress(),
				},
				PoolAssets: []sdk.Coin{
					sdk.NewCoin(
						"token",
						sdk.NewInt(10),
					),
					sdk.NewCoin(
						"foocoin",
						sdk.NewInt(10),
					),
				},
				AssetWeights: []sdk.Dec{sdk.NewDec(1)},
			},
			err: ErrInvalidLength,
		}, {
			name: "invalid weight",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      sdk.NewDec(10),
					ExitFee:      sdk.NewDec(10),
					FeeCollector: sample.AccAddress(),
				},
				PoolAssets: []sdk.Coin{
					sdk.NewCoin(
						"token",
						sdk.NewInt(10),
					),
					sdk.NewCoin(
						"foocoin",
						sdk.NewInt(10),
					),
				},
				AssetWeights: []sdk.Dec{sdk.NewDec(0), sdk.NewDec(1)},
			},

			err: ErrWeightZero,
		}, {
			name: "valid address",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      sdk.NewDec(10),
					ExitFee:      sdk.NewDec(10),
					FeeCollector: sample.AccAddress(),
				},
				PoolAssets: []sdk.Coin{
					sdk.NewCoin(
						"token",
						sdk.NewInt(10),
					),
					sdk.NewCoin(
						"foocoin",
						sdk.NewInt(10),
					),
				},
				AssetWeights: []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
