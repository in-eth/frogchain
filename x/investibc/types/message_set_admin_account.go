package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetAdminAccount = "set_admin_account"

var _ sdk.Msg = &MsgSetAdminAccount{}

func NewMsgSetAdminAccount(creator string, adminAccount string) *MsgSetAdminAccount {
	return &MsgSetAdminAccount{
		Creator:      creator,
		AdminAccount: adminAccount,
	}
}

func (msg *MsgSetAdminAccount) Route() string {
	return RouterKey
}

func (msg *MsgSetAdminAccount) Type() string {
	return TypeMsgSetAdminAccount
}

func (msg *MsgSetAdminAccount) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetAdminAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetAdminAccount) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
