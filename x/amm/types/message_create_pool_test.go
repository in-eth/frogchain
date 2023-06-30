package types

import (
	"testing"

	"frogchain/testutil/sample"

	"github.com/stretchr/testify/require"
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
				PoolParam: &PoolParam{
					SwapFee:      10,
					ExitFee:      10,
					FeeCollector: "invalid_address",
				},
				PoolAssets: []*PoolToken{},
			},
			err: ErrInvalidAddress,
		}, {
			name: "invalid address",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      10,
					ExitFee:      10,
					FeeCollector: "invalid_address",
				},
				PoolAssets: []*PoolToken{},
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
				PoolAssets: []*PoolToken{},
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
				PoolAssets: []*PoolToken{},
			},
			err: ErrFeeOverflow,
		}, {
			name: "assets not enough",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      10,
					ExitFee:      10,
					FeeCollector: sample.AccAddress(),
				},
				PoolAssets: []*PoolToken{
					&PoolToken{
						TokenDenom:   "123",
						TokenWeight:  1,
						TokenReserve: 0,
					},
				},
			},
			err: ErrInvalidAssets,
		}, {
			name: "same assets exist in assets",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      10,
					ExitFee:      10,
					FeeCollector: sample.AccAddress(),
				},
				PoolAssets: []*PoolToken{
					&PoolToken{
						TokenDenom:   "123",
						TokenWeight:  1,
						TokenReserve: 0,
					},
					&PoolToken{
						TokenDenom:   "123",
						TokenWeight:  1,
						TokenReserve: 0,
					},
				},
			},
			err: ErrInvalidAssets,
		}, {
			name: "valid address",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      10,
					ExitFee:      10,
					FeeCollector: sample.AccAddress(),
				},
				PoolAssets: []*PoolToken{
					&PoolToken{
						TokenDenom:   "123",
						TokenWeight:  1,
						TokenReserve: 0,
					},
					&PoolToken{
						TokenDenom:   "124",
						TokenWeight:  1,
						TokenReserve: 0,
					},
				},
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
