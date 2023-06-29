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

	// TODO: Handling the message

	// get pool
	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", msg.PoolId)
	}

	liquidityProvider, _ := sdk.AccAddressFromBech32(msg.Creator)

	// burn share token, amount is desiredAmount
	liquidity := msg.Liquidity * pool.PoolParam.ExitFee / (10 ^ 8)

	shareToken := sdk.NewCoin(
		types.ShareTokenIndex(msg.PoolId),
		sdk.NewInt(int64(liquidity)),
	)

	//move share tokens from account to module
	err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		liquidityProvider,
		types.ModuleName,
		sdk.NewCoins(shareToken),
	)
	if err != nil {
		return nil, err
	}

	// burn share token except fee
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(shareToken))
	if err != nil {
		return nil, err
	}

	// send fee share token to feeCollector
	fee := msg.Liquidity - liquidity
	feeCollector, err := sdk.AccAddressFromBech32(pool.PoolParam.FeeCollector)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.SendCoins(ctx, liquidityProvider, feeCollector, sdk.NewCoins(
		sdk.NewCoin(
			types.ShareTokenIndex(msg.PoolId),
			sdk.NewInt(int64(fee)),
		),
	))
	if err != nil {
		return nil, err
	}

	pool.ShareToken.TokenReserve -= liquidity

	// send assets from pool to account
	receiveTokens := sdk.NewCoins()
	for i, minAmount := range msg.MinAmounts {
		castAmount := liquidity * pool.PoolAssets[i].TokenReserve / pool.ShareToken.TokenReserve

		if castAmount < minAmount {
			return nil, sdkerrors.Wrapf(types.ErrInvalidAmount,
				"calculated amount is below minimum, %d, %d, %d",
				fmt.Sprint(i),
				fmt.Sprint(castAmount),
				fmt.Sprint(minAmount),
			)
		}

		receiveTokens.Add(
			sdk.NewCoin(
				pool.PoolAssets[i].TokenDenom,
				math.NewInt(int64(castAmount)),
			),
		)

		pool.PoolAssets[i].TokenReserve -= castAmount
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		liquidityProvider,
		receiveTokens,
	)
	if err != nil {
		return nil, err
	}

	// update pool
	k.SetPool(ctx, pool)

	return &types.MsgRemoveLiquidityResponse{
		ReceivedTokens: receiveTokens,
	}, nil
}
