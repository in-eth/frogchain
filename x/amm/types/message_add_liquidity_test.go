package types

import (
	"testing"

	"frogchain/testutil/sample"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgAddLiquidity_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAddLiquidity
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgAddLiquidity{
				Creator:        "invalid_address",
				PoolId:         1,
				DesiredAmounts: []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10)},
				MinAmounts:     []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10)},
			},
			err: ErrInvalidAddress,
		}, {
			name: "invalid desired amounts",
			msg: MsgAddLiquidity{
				Creator:        sample.AccAddress(),
				PoolId:         1,
				DesiredAmounts: []sdk.Dec{sdk.NewDec(9), sdk.NewDec(10)},
				MinAmounts:     []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10)},
			},
			err: ErrInvalidAmount,
		}, {
			name: "zero desired amounts",
			msg: MsgAddLiquidity{
				Creator:        sample.AccAddress(),
				PoolId:         1,
				DesiredAmounts: []sdk.Dec{sdk.NewDec(0), sdk.NewDec(10)},
				MinAmounts:     []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10)},
			},
			err: ErrInvalidAmount,
		}, {
			name: "valid address",
			msg: MsgAddLiquidity{
				Creator:        sample.AccAddress(),
				PoolId:         1,
				DesiredAmounts: []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10)},
				MinAmounts:     []sdk.Dec{sdk.NewDec(10), sdk.NewDec(10)},
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
