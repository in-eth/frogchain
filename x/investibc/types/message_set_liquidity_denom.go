package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetLiquidityDenom = "set_liquidity_denom"

var _ sdk.Msg = &MsgSetLiquidityDenom{}

func NewMsgSetLiquidityDenom(creator string, denom string) *MsgSetLiquidityDenom {
	return &MsgSetLiquidityDenom{
		Creator: creator,
		Denom:   denom,
	}
}

func (msg *MsgSetLiquidityDenom) Route() string {
	return RouterKey
}

func (msg *MsgSetLiquidityDenom) Type() string {
	return TypeMsgSetLiquidityDenom
}

func (msg *MsgSetLiquidityDenom) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetLiquidityDenom) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetLiquidityDenom) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
