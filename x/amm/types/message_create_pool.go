package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgCreatePool = "create_pool"

var _ sdk.Msg = &MsgCreatePool{}

func NewMsgCreatePool(creator string, poolParam PoolParam, poolAssets []sdk.Coin, assetWeights []sdk.Dec) *MsgCreatePool {
	return &MsgCreatePool{
		Creator:      creator,
		PoolParam:    &poolParam,
		PoolAssets:   poolAssets,
		AssetWeights: assetWeights,
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
	if swapFeeAmount.GTE(sdk.NewDec(TOTALPERCENT)) {
		return ErrFeeOverflow
	}

	exitFeeAmount := msg.PoolParam.ExitFee
	if exitFeeAmount.GTE(sdk.NewDec(TOTALPERCENT)) {
		return ErrFeeOverflow
	}

	if len(msg.PoolAssets) == 1 {
		return ErrInvalidAssets
	}
	if len(msg.PoolAssets) != len(msg.AssetWeights) {
		return ErrInvalidLength
	}

	for _, weight := range msg.AssetWeights {
		if weight.Equal(sdk.NewDec(0)) {
			return ErrWeightZero
		}
	}

	return nil
}
