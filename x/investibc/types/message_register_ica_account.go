package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRegisterIcaAccount = "register_ica_account"

var _ sdk.Msg = &MsgRegisterIcaAccount{}

func NewMsgRegisterIcaAccount(creator string, connectionId string, version string) *MsgRegisterIcaAccount {
	return &MsgRegisterIcaAccount{
		Creator:      creator,
		ConnectionId: connectionId,
		Version:      version,
	}
}

func (msg *MsgRegisterIcaAccount) Route() string {
	return RouterKey
}

func (msg *MsgRegisterIcaAccount) Type() string {
	return TypeMsgRegisterIcaAccount
}

func (msg *MsgRegisterIcaAccount) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRegisterIcaAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterIcaAccount) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
