package keeper

import (
	"frogchain/x/investibc/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.AdminAccount(ctx),
		k.DepositDenom(ctx),
		k.CurrentDepositAmount(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// SetAdminAccountParam set the adminAccount param
func (k Keeper) SetAdminAccountParam(ctx sdk.Context, adminAccount string) {
	params := k.GetParams(ctx)
	params.AdminAccount = adminAccount
	k.paramstore.SetParamSet(ctx, &params)
}

// SetParams set the params
func (k Keeper) SetDepositDenomParam(ctx sdk.Context, depositDenom string) {
	params := k.GetParams(ctx)
	params.DepositDenom = depositDenom
	if params.CurrentDepositAmount.Amount.GT(math.ZeroInt()) {
		panic("could not change deposit denom : current deposit amount is not zero")
	}
	params.CurrentDepositAmount.Denom = depositDenom
	k.paramstore.SetParamSet(ctx, &params)
}

// SetParams set the params
func (k Keeper) SetCurrentDepositAmountParam(ctx sdk.Context, currentDepositAmount sdk.Coin) {
	params := k.GetParams(ctx)
	params.CurrentDepositAmount = currentDepositAmount
	k.paramstore.SetParamSet(ctx, &params)
}

// AdminAccount returns the AdminAccount param
func (k Keeper) AdminAccount(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyAdminAccount, &res)
	return
}

// DepositDenom returns the DepositDenom param
func (k Keeper) DepositDenom(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyDepositDenom, &res)
	return
}

// CurrentDepositAmount returns the CurrentDepositAmount param
func (k Keeper) CurrentDepositAmount(ctx sdk.Context) (res sdk.Coin) {
	k.paramstore.Get(ctx, types.KeyCurrentDepositAmount, &res)
	return
}
