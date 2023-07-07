package keeper

import (
	"context"
	"fmt"
	"math"

	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) AddLiquidity(goCtx context.Context, msg *types.MsgAddLiquidity) (*types.MsgAddLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get provider address from msg.Creater
	liquidityProvider := sdk.MustAccAddressFromBech32(msg.Creator)

	// get share token from pool with pool id
	shareToken, err := k.GetPoolShareTokenForId(ctx, msg.PoolId)
	if err != nil {
		return nil, err
	}

	if assetLength, _ := k.GetPoolAssetsLength(ctx, msg.PoolId); assetLength != len(msg.DesiredAmounts) || len(msg.DesiredAmounts) != len(msg.MinAmounts) {
		return nil, types.ErrInvalidLength
	}

	// get liquidity amount
	// liquidity amount is the minimum value of the list with this formula
	// liquidityAmount = tokenAmountIn * shareTokenAmount / tokenAmount
	liquidityAmount := sdk.MustNewDecFromStr(fmt.Sprintf("%d", math.MaxInt))

	for i, desiredAmount := range msg.DesiredAmounts {
		poolAsset, err := k.GetPoolTokenForId(ctx, msg.PoolId, uint64(i))
		if err != nil {
			return nil, err
		}

		castAmount := desiredAmount.MulInt(shareToken.Amount).QuoInt(poolAsset.Amount)
		if liquidityAmount.GT(castAmount) {
			liquidityAmount = castAmount
		}
	}

	if liquidityAmount.Equal(sdk.ZeroDec()) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAmount, "no liquidity with amounts you deposit")
	}

	// calculate asset amounts from account to pool and update pool data
	receiveCoins := sdk.NewCoins()
	for i, minAmount := range msg.MinAmounts {
		poolAsset, err := k.GetPoolTokenForId(ctx, msg.PoolId, uint64(i))
		if err != nil {
			return nil, err
		}

		castAmount := liquidityAmount.MulInt(poolAsset.Amount).QuoInt(shareToken.Amount)

		// if input token amount is below min amount, then revert
		if castAmount.LT(minAmount) {
			return nil, sdkerrors.Wrapf(types.ErrInvalidAmount,
				"calculated amount is below min amount, %s, %s, %s",
				fmt.Sprint(i),
				fmt.Sprint(castAmount),
				fmt.Sprint(minAmount),
			)
		}

		receiveCoins = receiveCoins.Add(
			sdk.NewCoin(
				poolAsset.Denom,
				msg.DesiredAmounts[i].RoundInt(),
			),
		)

		poolAsset = poolAsset.Add(
			sdk.NewCoin(
				poolAsset.Denom,
				castAmount.TruncateInt(),
			),
		)
		err = k.SetPoolToken(ctx, msg.PoolId, uint64(i), poolAsset)
		if err != nil {
			return nil, err
		}
	}

	// send assets from account to pool
	err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		liquidityProvider,
		types.ModuleName,
		receiveCoins,
	)
	if err != nil {
		return nil, err
	}

	// mint new share token, amount is liquidityAmount
	newShareToken := sdk.NewCoin(
		shareToken.Denom,
		liquidityAmount.TruncateInt(),
	)

	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(newShareToken))
	if err != nil {
		return nil, err
	}

	shareToken = shareToken.Add(newShareToken)
	err = k.SetPoolShareToken(ctx, msg.PoolId, shareToken)
	if err != nil {
		return nil, err
	}

	// send share token to liquidity provider
	err = k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		liquidityProvider,
		sdk.NewCoins(newShareToken),
	)
	if err != nil {
		return nil, err
	}

	// emit mint event
	ctx.EventManager().EmitEvent(
		types.NewAddLiquidityEvent(liquidityProvider, msg.PoolId, newShareToken),
	)

	return &types.MsgAddLiquidityResponse{
		ShareToken: newShareToken,
	}, nil
}
