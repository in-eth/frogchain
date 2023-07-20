package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetDepositDenom = "set_deposit_denom"

var _ sdk.Msg = &MsgSetDepositDenom{}

func NewMsgSetDepositDenom(creator string, denom string) *MsgSetDepositDenom {
	return &MsgSetDepositDenom{
		Creator: creator,
		Denom:   denom,
	}
}

func (msg *MsgSetDepositDenom) Route() string {
	return RouterKey
}

func (msg *MsgSetDepositDenom) Type() string {
	return TypeMsgSetDepositDenom
}

func (msg *MsgSetDepositDenom) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetDepositDenom) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetDepositDenom) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
