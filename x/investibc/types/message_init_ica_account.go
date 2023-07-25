package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgInitIcaAccount = "Init_ica_account"

var _ sdk.Msg = &MsgInitIcaAccount{}

func NewMsgInitIcaAccount(creator string, connectionId string) *MsgInitIcaAccount {
	return &MsgInitIcaAccount{
		Creator:      creator,
		ConnectionId: connectionId,
	}
}

func (msg *MsgInitIcaAccount) Route() string {
	return RouterKey
}

func (msg *MsgInitIcaAccount) Type() string {
	return TypeMsgInitIcaAccount
}

func (msg *MsgInitIcaAccount) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgInitIcaAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgInitIcaAccount) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
