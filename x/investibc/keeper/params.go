package keeper

import (
	"frogchain/x/investibc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.AdminAccount(ctx),
		k.DepositDenom(ctx),
		k.CurrentDepositAmount(ctx),
		k.LiquidityDenom(ctx),
		k.LockTokenTimestamp(ctx),
		k.DepositLastTime(ctx),
		k.IcaConnectionId(ctx),
		k.DepositTokenToICAPacketSend(ctx),
		k.JoinSwapExactAmountInPacketSend(ctx),
		k.LockTokensPacketSend(ctx),
		k.UnLockLiquidityPacketSend(ctx),
		k.ClaimRewardPacketSend(ctx),
		k.DepositTokenToICAPacketSent(ctx),
		k.JoinSwapExactAmountInPacketSent(ctx),
		k.LockTokensPacketSent(ctx),
		k.UnLockLiquidityPacketSent(ctx),
		k.ClaimRewardPacketSent(ctx),
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
	k.paramstore.SetParamSet(ctx, &params)
}

// SetParams set the params
func (k Keeper) SetLockTokenTimestampParam(ctx sdk.Context, lockTokenTimestamp uint64) {
	params := k.GetParams(ctx)
	params.LockTokenTimestamp = lockTokenTimestamp
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
func (k Keeper) SetDepositTokenToICAPacketSendParam(ctx sdk.Context, depositTokenToICAPacketSend bool) {
	params := k.GetParams(ctx)
	params.DepositTokenToICAPacketSend = depositTokenToICAPacketSend
	k.paramstore.SetParamSet(ctx, &params)
}

// SetParams set the params
func (k Keeper) SetJoinSwapExactAmountInPacketSendParam(ctx sdk.Context, joinSwapExactAmountInPacketSend bool) {
	params := k.GetParams(ctx)
	params.JoinSwapExactAmountInPacketSend = joinSwapExactAmountInPacketSend
	k.paramstore.SetParamSet(ctx, &params)
}

// SetParams set the params
func (k Keeper) SetLockTokensPacketSendParam(ctx sdk.Context, lockTokensPacketSend bool) {
	params := k.GetParams(ctx)
	params.LockTokensPacketSend = lockTokensPacketSend
	k.paramstore.SetParamSet(ctx, &params)
}

// SetParams set the params
func (k Keeper) SetUnLockLiquidityPacketSendParam(ctx sdk.Context, unLockLiquidityPacketSend bool) {
	params := k.GetParams(ctx)
	params.UnLockLiquidityPacketSend = unLockLiquidityPacketSend
	k.paramstore.SetParamSet(ctx, &params)
}

// SetParams set the params
func (k Keeper) SetClaimRewardPacketSendParam(ctx sdk.Context, claimRewardPacketSend bool) {
	params := k.GetParams(ctx)
	params.ClaimRewardPacketSend = claimRewardPacketSend
	k.paramstore.SetParamSet(ctx, &params)
}

// SetParams set the params
func (k Keeper) SetDepositTokenToICAPacketSentParam(ctx sdk.Context, depositTokenToICAPacketSent bool) {
	params := k.GetParams(ctx)
	params.DepositTokenToICAPacketSent = depositTokenToICAPacketSent
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

// SetParams set the params
func (k Keeper) SetUnLockLiquidityPacketSentParam(ctx sdk.Context, unLockLiquidityPacketSent bool) {
	params := k.GetParams(ctx)
	params.UnLockLiquidityPacketSent = unLockLiquidityPacketSent
	k.paramstore.SetParamSet(ctx, &params)
}

// SetParams set the params
func (k Keeper) SetClaimRewardPacketSentParam(ctx sdk.Context, claimRewardPacketSent bool) {
	params := k.GetParams(ctx)
	params.ClaimRewardPacketSent = claimRewardPacketSent
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

// LiquidityDenom returns the LiquidityDenom param
func (k Keeper) LiquidityDenom(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyLiquidityDenom, &res)
	return
}

// CurrentLiquidityAmount returns the CurrentLiquidityAmount param
func (k Keeper) LockTokenTimestamp(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyLockTokenTimestamp, &res)
	return
}

// DepositLastTime returns the DepositLastTime param
func (k Keeper) DepositLastTime(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyDepositLastTime, &res)
	return
}

// IcaConnectionId returns the IcaConnectionId param
func (k Keeper) IcaConnectionId(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyIcaConnectionId, &res)
	return
}

// DepositTokenToICAPacketSend returns the DepositTokenToICAPacketSend param
func (k Keeper) DepositTokenToICAPacketSend(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyDepositTokenToICAPacketSend, &res)
	return
}

// JoinSwapExactAmountInPacketSend returns the JoinSwapExactAmountInPacketSend param
func (k Keeper) JoinSwapExactAmountInPacketSend(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyJoinSwapExactAmountInPacketSend, &res)
	return
}

// LockTokensPacketSend returns the LockTokensPacketSend param
func (k Keeper) LockTokensPacketSend(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyLockTokensPacketSend, &res)
	return
}

// UnLockLiquidityPacketSend returns the UnLockLiquidityPacketSend param
func (k Keeper) UnLockLiquidityPacketSend(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyUnLockLiquidityPacketSend, &res)
	return
}

// ClaimRewardPacketSend returns the ClaimRewardPacketSend param
func (k Keeper) ClaimRewardPacketSend(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyClaimRewardPacketSend, &res)
	return
}

// DepositTokenToICAPacketSent returns the DepositTokenToICAPacketSent param
func (k Keeper) DepositTokenToICAPacketSent(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyDepositTokenToICAPacketSent, &res)
	return
}

// JoinSwapExactAmountInPacketSent returns the JoinSwapExactAmountInPacketSent param
func (k Keeper) JoinSwapExactAmountInPacketSent(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyJoinSwapExactAmountInPacketSent, &res)
	return
}

// LockTokensPacketSent returns the LockTokensPacketSent param
func (k Keeper) LockTokensPacketSent(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyLockTokensPacketSent, &res)
	return
}

// UnLockLiquidityPacketSent returns the UnLockLiquidityPacketSent param
func (k Keeper) UnLockLiquidityPacketSent(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyUnLockLiquidityPacketSent, &res)
	return
}

// ClaimRewardPacketSent returns the ClaimRewardPacketSent param
func (k Keeper) ClaimRewardPacketSent(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyClaimRewardPacketSent, &res)
	return
}
