package keeper

import (
	"context"
	"fmt"

	"frogchain/x/amm/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	liquidityProvider, _ := sdk.AccAddressFromBech32(msg.Creator)

	// burn share token, amount is desiredAmount
	fee := msg.Liquidity * poolParam.ExitFee / (types.TOTALPERCENT)
	liquidity := msg.Liquidity - fee

	burnShareToken := sdk.NewCoin(
		shareToken.Denom,
		sdk.NewInt(int64(liquidity)),
	)

	if poolAssetsLen, _ := k.GetPoolAssetsLength(ctx, msg.PoolId); len(msg.MinAmounts) != poolAssetsLen {
		return nil, types.ErrInvalidLength
	}

	// send assets from pool to account
	receiveTokens := sdk.NewCoins()
	for i, minAmount := range msg.MinAmounts {
		// get pool asset
		poolAsset, err := k.GetPoolTokenForId(ctx, msg.PoolId, uint64(i))
		if err != nil {
			return nil, err
		}

		// calculate token amount for liquidity
		castAmount := liquidity * poolAsset.Amount.Uint64() / shareToken.Amount.Uint64()

		if castAmount < minAmount {
			return nil, ErrorWrap(types.ErrInvalidAmount,
				"calculated amount is below minimum, %s, %s, %s",
				fmt.Sprint(i),
				fmt.Sprint(castAmount),
				fmt.Sprint(minAmount),
			)
		}

		receiveTokens = receiveTokens.Add(
			sdk.NewCoin(
				poolAsset.Denom,
				math.NewInt(int64(castAmount)),
			),
		)

		// update pool asset data
		poolAsset = poolAsset.SubAmount(math.NewInt(int64(castAmount)))
		k.SetPoolToken(ctx, msg.PoolId, uint64(i), poolAsset)
	}

	// send tokens from module to account
	err = k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		liquidityProvider,
		receiveTokens,
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

	// get pool fee collector
	feeCollector, err := sdk.AccAddressFromBech32(poolParam.FeeCollector)
	if err != nil {
		return nil, err
	}

	// send fee share token to fee collector
	err = k.bankKeeper.SendCoins(ctx, liquidityProvider, feeCollector, sdk.NewCoins(
		sdk.NewCoin(
			shareToken.Denom,
			sdk.NewInt(int64(fee)),
		),
	))
	if err != nil {
		return nil, err
	}

	// update pool share token data
	shareToken = shareToken.SubAmount(sdk.NewInt(int64(liquidity)))
	err = k.SetPoolShareToken(ctx, msg.PoolId, shareToken)
	if err != nil {
		return nil, err
	}

	// emit mint event
	ctx.EventManager().EmitEvent(
		types.NewRemoveLiquidityEvent(liquidityProvider, msg.PoolId, receiveTokens),
	)

	return &types.MsgRemoveLiquidityResponse{
		ReceivedTokens: receiveTokens,
	}, nil
}
