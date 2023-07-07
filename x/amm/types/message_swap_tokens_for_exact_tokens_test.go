package types

import (
	"testing"

	"frogchain/testutil/sample"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgSwapTokensForExactTokens_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSwapTokensForExactTokens
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSwapTokensForExactTokens{
				Creator:   "invalid_address",
				PoolId:    1,
				AmountOut: sdk.NewDec(10),
				Path:      []string{"123", "123"},
				To:        "invalid_address",
			},
			err: ErrInvalidAddress,
		}, {
			name: "invalid to address",
			msg: MsgSwapTokensForExactTokens{
				Creator:   sample.AccAddress(),
				PoolId:    1,
				AmountOut: sdk.NewDec(10),
				Path:      []string{"123", "123"},
				To:        "invalid_address",
			},

			err: ErrInvalidAddress,
		}, {
			name: "invalid path",
			msg: MsgSwapTokensForExactTokens{
				Creator:   sample.AccAddress(),
				PoolId:    1,
				AmountOut: sdk.NewDec(10),
				Path:      []string{"123"},
				To:        sample.AccAddress(),
			},
			err: ErrInvalidPath,
		}, {
			name: "valid address",
			msg: MsgSwapTokensForExactTokens{
				Creator:   sample.AccAddress(),
				PoolId:    1,
				AmountOut: sdk.NewDec(10),
				Path:      []string{"123", "123"},
				To:        sample.AccAddress(),
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
