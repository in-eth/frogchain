package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"

	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) AddLiquidity(goCtx context.Context, msg *types.MsgAddLiquidity) (*types.MsgAddLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message

	// get pool
	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", msg.PoolId)
	}

	liquidityProvider, _ := sdk.AccAddressFromBech32(msg.Creator)

	// get liquidity amount
	liquidityAmount := msg.DesiredAmounts[0] * pool.ShareToken.TokenReserve / pool.PoolAssets[0].TokenReserve
	for i, desiredAmount := range msg.DesiredAmounts {
		castAmount := desiredAmount * pool.ShareToken.TokenReserve / pool.PoolAssets[i].TokenReserve
		if liquidityAmount > castAmount {
			liquidityAmount = castAmount
		}
	}

	if liquidityAmount == 0 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAmount, "no liquidity with amounts you deposit")
	}

	// send assets from account to pool
	collateral := sdk.NewCoins()
	for i, minAmount := range msg.MinAmounts {
		castAmount := liquidityAmount * pool.PoolAssets[i].TokenReserve / pool.ShareToken.TokenReserve

		if castAmount < minAmount {
			return nil, sdkerrors.Wrapf(types.ErrInvalidAmount,
				"calculated amount is below minimum, %d, %d, %d",
				fmt.Sprint(i),
				fmt.Sprint(castAmount),
				fmt.Sprint(minAmount),
			)
		}

		collateral.Add(
			sdk.NewCoin(
				pool.PoolAssets[i].TokenDenom,
				math.NewInt(int64(castAmount)),
			),
		)

		pool.PoolAssets[i].TokenReserve += castAmount
	}

	err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		liquidityProvider,
		types.ModuleName,
		collateral,
	)
	if err != nil {
		return nil, err
	}

	// mint new share token, amount is liquidityAmount
	shareToken := sdk.NewCoin(
		types.ShareTokenIndex(msg.PoolId),
		sdk.NewInt(int64(liquidityAmount)),
	)

	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(shareToken))
	if err != nil {
		return nil, err
	}

	pool.ShareToken.TokenReserve += liquidityAmount

	// send share token to liquidity provider
	err = k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		liquidityProvider,
		sdk.NewCoins(shareToken),
	)
	if err != nil {
		return nil, err
	}

	// update pool
	k.SetPool(ctx, pool)

	return &types.MsgAddLiquidityResponse{
		ShareToken: &shareToken,
	}, nil
}
