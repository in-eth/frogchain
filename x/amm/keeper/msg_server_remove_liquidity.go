package keeper

import (
	"context"
	"fmt"

	"frogchain/x/amm/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) RemoveLiquidity(goCtx context.Context, msg *types.MsgRemoveLiquidity) (*types.MsgRemoveLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get pool param
	poolParam, err := k.GetPoolParamForId(ctx, msg.PoolId)
	if err != nil {
		return nil, err
	}

	// get pool share token
	shareToken, err := k.GetPoolShareTokenForId(ctx, msg.PoolId)
	if err != nil {
		return nil, err
	}

	liquidityProvider := sdk.MustAccAddressFromBech32(msg.Creator)

	// get pool fee collector
	feeCollector := sdk.MustAccAddressFromBech32(poolParam.FeeCollector)

	// burn share token, amount is desiredAmount
	fee := msg.Liquidity.Mul(poolParam.ExitFee).QuoInt(math.NewInt(types.TOTALPERCENT))
	liquidity := msg.Liquidity.Sub(fee)

	burnShareToken := sdk.NewCoin(
		shareToken.Denom,
		liquidity.TruncateInt(),
	)

	if poolAssetsLen, _ := k.GetPoolAssetsLength(ctx, msg.PoolId); len(msg.MinAmounts) != poolAssetsLen {
		return nil, types.ErrInvalidLength
	}

	// send assets from pool to account
	sendTokens := sdk.NewCoins()
	for i, minAmount := range msg.MinAmounts {
		// get pool asset
		poolAsset, err := k.GetPoolTokenForId(ctx, msg.PoolId, uint64(i))
		if err != nil {
			return nil, err
		}

		// calculate token amount for liquidity
		castAmount := liquidity.MulInt(poolAsset.Amount).QuoInt(shareToken.Amount)

		if castAmount.LT(minAmount) {
			return nil, sdkerrors.Wrapf(types.ErrInvalidAmount,
				"calculated amount is below minimum, %s, %s, %s",
				fmt.Sprint(i),
				fmt.Sprint(castAmount),
				fmt.Sprint(minAmount),
			)
		}

		sendTokens = sendTokens.Add(
			sdk.NewCoin(
				poolAsset.Denom,
				castAmount.TruncateInt(),
			),
		)

		// update pool asset data
		poolAsset = poolAsset.SubAmount(castAmount.TruncateInt())
		k.SetPoolToken(ctx, msg.PoolId, uint64(i), poolAsset)
	}

	// send tokens from module to account
	err = k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		liquidityProvider,
		sendTokens,
	)
	if err != nil {
		return nil, err
	}

	//move share tokens from account to module
	err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		liquidityProvider,
		types.ModuleName,
		sdk.NewCoins(burnShareToken),
	)
	if err != nil {
		return nil, err
	}

	// burn share token except fee
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(burnShareToken))
	if err != nil {
		return nil, err
	}

	// send fee share token to fee collector
	err = k.bankKeeper.SendCoins(ctx, liquidityProvider, feeCollector, sdk.NewCoins(
		sdk.NewCoin(
			shareToken.Denom,
			fee.RoundInt(),
		),
	))
	if err != nil {
		return nil, err
	}

	// update pool share token data
	shareToken = shareToken.SubAmount(liquidity.TruncateInt())
	err = k.SetPoolShareToken(ctx, msg.PoolId, shareToken)
	if err != nil {
		return nil, err
	}

	// emit mint event
	ctx.EventManager().EmitEvent(
		types.NewRemoveLiquidityEvent(liquidityProvider, msg.PoolId, sendTokens),
	)

	return &types.MsgRemoveLiquidityResponse{
		ReceivedTokens: sendTokens,
	}, nil
}
