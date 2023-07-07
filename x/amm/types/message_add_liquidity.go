package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgAddLiquidity = "add_liquidity"

var _ sdk.Msg = &MsgAddLiquidity{}

func NewMsgAddLiquidity(creator string, poolId uint64, desiredAmounts []sdk.Dec, minAmounts []sdk.Dec) *MsgAddLiquidity {
	return &MsgAddLiquidity{
		Creator:        creator,
		PoolId:         poolId,
		DesiredAmounts: desiredAmounts,
		MinAmounts:     minAmounts,
	}
}

func (msg *MsgAddLiquidity) Route() string {
	return RouterKey
}

func (msg *MsgAddLiquidity) Type() string {
	return TypeMsgAddLiquidity
}

func (msg *MsgAddLiquidity) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddLiquidity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddLiquidity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return ErrInvalidAddress
	}

	for i, desiredAmount := range msg.DesiredAmounts {
		if desiredAmount.LT(msg.MinAmounts[i]) || desiredAmount.Equal(sdk.NewDec(0)) || msg.MinAmounts[i].Equal(sdk.NewDec(0)) {
			return ErrInvalidAmount
		}
	}
	return nil
}
