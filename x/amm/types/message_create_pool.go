package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgCreatePool = "create_pool"

var _ sdk.Msg = &MsgCreatePool{}

func NewMsgCreatePool(creator string, poolParam *PoolParam, poolAssets []*PoolToken, assetAmounts []uint64) *MsgCreatePool {
	return &MsgCreatePool{
		Creator:      creator,
		PoolParam:    poolParam,
		PoolAssets:   poolAssets,
		AssetAmounts: assetAmounts,
	}
}

func (msg *MsgCreatePool) Route() string {
	return RouterKey
}

func (msg *MsgCreatePool) Type() string {
	return TypeMsgCreatePool
}

func (msg *MsgCreatePool) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreatePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreatePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return ErrInvalidAddress
	}
	_, err = sdk.AccAddressFromBech32(msg.PoolParam.FeeCollector)
	if err != nil {
		return ErrInvalidAddress
	}

	swapFeeAmount := msg.PoolParam.SwapFee
	if swapFeeAmount >= TOTALPERCENT {
		return ErrFeeOverflow
	}

	exitFeeAmount := msg.PoolParam.ExitFee
	if exitFeeAmount >= TOTALPERCENT {
		return ErrFeeOverflow
	}

	for _, poolAsset := range msg.PoolAssets {
		weight := poolAsset.TokenWeight
		if weight == 0 {
			return ErrWeightZero
		}
	}
	return nil
}
