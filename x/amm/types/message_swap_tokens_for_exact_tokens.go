package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSwapTokensForExactTokens = "swap_tokens_for_exact_tokens"

var _ sdk.Msg = &MsgSwapTokensForExactTokens{}

func NewMsgSwapTokensForExactTokens(creator string, poolId uint64, amountOut sdk.Dec, path []string, to string, deadline time.Time) *MsgSwapTokensForExactTokens {
	return &MsgSwapTokensForExactTokens{
		Creator:   creator,
		PoolId:    poolId,
		AmountOut: amountOut,
		Path:      path,
		To:        to,
		Deadline:  deadline,
	}
}

func (msg *MsgSwapTokensForExactTokens) Route() string {
	return RouterKey
}

func (msg *MsgSwapTokensForExactTokens) Type() string {
	return TypeMsgSwapTokensForExactTokens
}

func (msg *MsgSwapTokensForExactTokens) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSwapTokensForExactTokens) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSwapTokensForExactTokens) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return ErrInvalidAddress
	}

	_, err = sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		return ErrInvalidAddress
	}

	if len(msg.Path) == 1 {
		return ErrInvalidPath
	}
	return nil
}
