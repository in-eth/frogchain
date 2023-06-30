package types

import (
	"testing"

	"frogchain/testutil/sample"

	"github.com/stretchr/testify/require"
)

func TestMsgRemoveLiquidity_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgRemoveLiquidity
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgRemoveLiquidity{
				Creator: "invalid_address",
			},
			err: ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgRemoveLiquidity{
				Creator: sample.AccAddress(),
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
