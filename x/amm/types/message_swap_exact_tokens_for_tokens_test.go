package types

import (
	"testing"

	"frogchain/testutil/sample"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgSwapExactTokensForTokens_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSwapExactTokensForTokens
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSwapExactTokensForTokens{
				Creator:      "invalid_address",
				PoolId:       1,
				AmountIn:     sdk.NewDec(10),
				AmountOutMin: sdk.NewDec(100),
				Path:         []string{"1"},
				To:           "invalid_address",
			},
			err: ErrInvalidAddress,
		}, {
			name: "invalid to address",
			msg: MsgSwapExactTokensForTokens{
				Creator:      sample.AccAddress(),
				PoolId:       1,
				AmountIn:     sdk.NewDec(10),
				AmountOutMin: sdk.NewDec(100),
				Path:         []string{"1"},
				To:           "invalid_address",
			},
			err: ErrInvalidAddress,
		}, {
			name: "invalid path",
			msg: MsgSwapExactTokensForTokens{
				Creator:      sample.AccAddress(),
				PoolId:       1,
				AmountIn:     sdk.NewDec(10),
				AmountOutMin: sdk.NewDec(100),
				Path:         []string{"1"},
				To:           sample.AccAddress(),
			},
			err: ErrInvalidPath,
		}, {
			name: "valid address",
			msg: MsgSwapExactTokensForTokens{
				Creator:      sample.AccAddress(),
				PoolId:       1,
				AmountIn:     sdk.NewDec(10),
				AmountOutMin: sdk.NewDec(100),
				Path:         []string{"1", "2"},
				To:           sample.AccAddress(),
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
