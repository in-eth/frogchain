package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSwapExactTokensForTokens = "swap_exact_tokens_for_tokens"

var _ sdk.Msg = &MsgSwapExactTokensForTokens{}

func NewMsgSwapExactTokensForTokens(creator string, poolId uint64, amountIn uint64, amountOutMin uint64, path []string, to string, deadline string) *MsgSwapExactTokensForTokens {
	return &MsgSwapExactTokensForTokens{
		Creator:      creator,
		PoolId:       poolId,
		AmountIn:     amountIn,
		AmountOutMin: amountOutMin,
		Path:         path,
		To:           to,
		Deadline:     deadline,
	}
}

func (msg *MsgSwapExactTokensForTokens) Route() string {
	return RouterKey
}

func (msg *MsgSwapExactTokensForTokens) Type() string {
	return TypeMsgSwapExactTokensForTokens
}

func (msg *MsgSwapExactTokensForTokens) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSwapExactTokensForTokens) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSwapExactTokensForTokens) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return ErrInvalidAddress
	}
	_, err = sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		return ErrInvalidAddress
	}

	// deadline, err := time.Parse(DeadlineLayout, msg.Deadline)
	// if deadline < .BlockTime() {

	// }

	if len(msg.Path) == 1 {
		return ErrInvalidPath
	}

	return nil
}
