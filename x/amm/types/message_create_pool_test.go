package types

import (
	"frogchain/testutil/sample"
	"testing"

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
			name: "assets not enough",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      10,
					ExitFee:      10,
					FeeCollector: sample.AccAddress(),
				},
				PoolAssets: []PoolToken{
					PoolToken{
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
				PoolAssets: []PoolToken{
					PoolToken{
						TokenDenom:   "123",
						TokenWeight:  1,
						TokenReserve: 0,
					},
					PoolToken{
						TokenDenom:   "123",
						TokenWeight:  1,
						TokenReserve: 0,
					},
				},
			},
			err: ErrDuplicateAssets,
		}, {
			name: "invalid amounts length",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      10,
					ExitFee:      10,
					FeeCollector: sample.AccAddress(),
				},
				PoolAssets: []PoolToken{
					PoolToken{
						TokenDenom:   "123",
						TokenWeight:  1,
						TokenReserve: 0,
					},
					PoolToken{
						TokenDenom:   "124",
						TokenWeight:  1,
						TokenReserve: 0,
					},
				},
				AssetAmounts: []uint64{10},
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
				PoolAssets: []PoolToken{
					PoolToken{
						TokenDenom:   "123",
						TokenWeight:  0,
						TokenReserve: 0,
					},
					PoolToken{
						TokenDenom:   "124",
						TokenWeight:  1,
						TokenReserve: 0,
					},
				},
				AssetAmounts: []uint64{10, 10},
			},

			err: ErrWeightZero,
		}, {
			name: "invalid asset amounts",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      10,
					ExitFee:      10,
					FeeCollector: sample.AccAddress(),
				},
				PoolAssets: []PoolToken{
					PoolToken{
						TokenDenom:   "123",
						TokenWeight:  1,
						TokenReserve: 0,
					},
					PoolToken{
						TokenDenom:   "124",
						TokenWeight:  1,
						TokenReserve: 0,
					},
				},
				AssetAmounts: []uint64{0, 10},
			},
			err: ErrInvalidAmount,
		}, {
			name: "valid address",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				PoolParam: &PoolParam{
					SwapFee:      10,
					ExitFee:      10,
					FeeCollector: sample.AccAddress(),
				},
				PoolAssets: []PoolToken{
					PoolToken{
						TokenDenom:   "123",
						TokenWeight:  1,
						TokenReserve: 0,
					},
					PoolToken{
						TokenDenom:   "124",
						TokenWeight:  1,
						TokenReserve: 0,
					},
				},
				AssetAmounts: []uint64{10, 10},
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
