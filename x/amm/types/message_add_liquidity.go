package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddLiquidity = "add_liquidity"

var _ sdk.Msg = &MsgAddLiquidity{}

func NewMsgAddLiquidity(creator string, poolId uint64, desiredAmounts []uint64, minAmounts []uint64) *MsgAddLiquidity {
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
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	for i, desiredAmount := range msg.DesiredAmounts {
		if desiredAmount < msg.MinAmounts[i] {
			return sdkerrors.Wrapf(
				ErrInvalidAmount,
				"invalid desired and min amounts (%s, %s)",
				fmt.Sprint(desiredAmount),
				fmt.Sprint(msg.MinAmounts[i]))
		}
	}
	return nil
}
