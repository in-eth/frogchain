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
		k.LiquidityDenom(ctx),
		k.CurrentLiquidityAmount(ctx),
		k.DepositLastTime(ctx),
		k.IcaConnectionId(ctx),
		k.JoinSwapExactAmountInPacketSent(ctx),
		k.LockTokensPacketSent(ctx),
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

// SetParams set the params
func (k Keeper) SetLiquidityDenomParam(ctx sdk.Context, liquidityDenom string) {
	params := k.GetParams(ctx)
	params.LiquidityDenom = liquidityDenom
	if params.CurrentLiquidityAmount.Amount.GT(math.ZeroInt()) {
		panic("could not change liquidity denom : current liquidity amount is not zero")
	}
	params.CurrentLiquidityAmount.Denom = liquidityDenom
	k.paramstore.SetParamSet(ctx, &params)
}

// SetParams set the params
func (k Keeper) SetCurrentLiquidityAmountParam(ctx sdk.Context, currentLiquidityAmount sdk.Coin) {
	params := k.GetParams(ctx)
	params.CurrentLiquidityAmount = currentLiquidityAmount
	k.paramstore.SetParamSet(ctx, &params)
}

// SetParams set the params
func (k Keeper) SetDepositLastTimeParam(ctx sdk.Context, depositLastTime uint64) {
	params := k.GetParams(ctx)
	params.DepositLastTime = depositLastTime
	k.paramstore.SetParamSet(ctx, &params)
}

// SetParams set the params
func (k Keeper) SetIcaConnectionIdParam(ctx sdk.Context, icaConnectionId string) {
	params := k.GetParams(ctx)
	params.IcaConnectionId = icaConnectionId
	k.paramstore.SetParamSet(ctx, &params)
}

// SetParams set the params
func (k Keeper) SetJoinSwapExactAmountInPacketSentParam(ctx sdk.Context, joinSwapExactAmountInPacketSent bool) {
	params := k.GetParams(ctx)
	params.JoinSwapExactAmountInPacketSent = joinSwapExactAmountInPacketSent
	k.paramstore.SetParamSet(ctx, &params)
}

// SetParams set the params
func (k Keeper) SetLockTokensPacketSentParam(ctx sdk.Context, lockTokensPacketSent bool) {
	params := k.GetParams(ctx)
	params.LockTokensPacketSent = lockTokensPacketSent
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

// DepositDenom returns the DepositDenom param
func (k Keeper) LiquidityDenom(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyLiquidityDenom, &res)
	return
}

// CurrentDepositAmount returns the CurrentDepositAmount param
func (k Keeper) CurrentLiquidityAmount(ctx sdk.Context) (res sdk.Coin) {
	k.paramstore.Get(ctx, types.KeyCurrentLiquidityAmount, &res)
	return
}

// CurrentDepositAmount returns the CurrentDepositAmount param
func (k Keeper) DepositLastTime(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyDepositLastTime, &res)
	return
}

// CurrentDepositAmount returns the CurrentDepositAmount param
func (k Keeper) IcaConnectionId(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyIcaConnectionId, &res)
	return
}

// CurrentDepositAmount returns the CurrentDepositAmount param
func (k Keeper) JoinSwapExactAmountInPacketSent(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyJoinSwapExactAmountInPacketSent, &res)
	return
}

// CurrentDepositAmount returns the CurrentDepositAmount param
func (k Keeper) LockTokensPacketSent(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyLockTokensPacketSent, &res)
	return
}
