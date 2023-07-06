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
					SwapFee:      10,
					ExitFee:      10,
					FeeCollector: "invalid_address",
				},
			},
			err: ErrInvalidAddress,
		}, {
			name: "exit fee amount overflow",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      10,
					ExitFee:      100000001,
					FeeCollector: sample.AccAddress(),
				},
			},
			err: ErrFeeOverflow,
		}, {
			name: "swap fee amount overflow",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      100000001,
					ExitFee:      10,
					FeeCollector: sample.AccAddress(),
				},
			},
			err: ErrFeeOverflow,
		}, {
			name: "invalid weights length",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      10,
					ExitFee:      10,
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
				AssetWeights: []uint64{1},
			},
			err: ErrInvalidLength,
		}, {
			name: "invalid weight",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      10,
					ExitFee:      10,
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
				AssetWeights: []uint64{0, 1},
			},

			err: ErrWeightZero,
		}, {
			name: "valid address",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      10,
					ExitFee:      10,
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
				AssetWeights: []uint64{1, 1},
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
