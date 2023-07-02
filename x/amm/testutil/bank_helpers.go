package testutil

import (
	"context"
	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
)

func (escrow *MockBankKeeper) ExpectAny(context context.Context) {
	escrow.EXPECT().SendCoinsFromAccountToModule(sdk.UnwrapSDKContext(context), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	escrow.EXPECT().SendCoinsFromModuleToAccount(sdk.UnwrapSDKContext(context), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
}

func coinsOf(tokenDenom string, amount uint64) sdk.Coins {
	return sdk.Coins{
		sdk.Coin{
			Denom:  tokenDenom,
			Amount: sdk.NewInt(int64(amount)),
		},
	}
}

func (escrow *MockBankKeeper) ExpectPay(context context.Context, who string, tokenDenom string, amount uint64) *gomock.Call {
	whoAddr, err := sdk.AccAddressFromBech32(who)
	if err != nil {
		panic(err)
	}
	return escrow.EXPECT().SendCoinsFromAccountToModule(sdk.UnwrapSDKContext(context), whoAddr, types.ModuleName, coinsOf(tokenDenom, amount))
}

func (escrow *MockBankKeeper) ExpectRefund(context context.Context, who string, tokenDenom string, amount uint64) *gomock.Call {
	whoAddr, err := sdk.AccAddressFromBech32(who)
	if err != nil {
		panic(err)
	}
	return escrow.EXPECT().SendCoinsFromModuleToAccount(sdk.UnwrapSDKContext(context), types.ModuleName, whoAddr, coinsOf(tokenDenom, amount))
}
