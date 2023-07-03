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

	// TODO: Handling the message

	// get share token from pool with pool id
	shareToken, err := k.GetPoolShareTokenForId(ctx, msg.PoolId)
	if err != nil {
		return nil, err
	}

	// get provider address from msg.Creater
	liquidityProvider, _ := sdk.AccAddressFromBech32(msg.Creator)

	// get liquidity amount
	// liquidity amount is the minimum value of the list with this formula
	// liquidityAmount = tokenAmountIn * shareTokenAmount / tokenAmount
	liquidityAmount := uint64(math.MaxUint64)
	for i, desiredAmount := range msg.DesiredAmounts {
		poolAsset, err := k.GetPoolTokenForId(ctx, msg.PoolId, uint64(i))
		if err != nil {
			return nil, err
		}

		castAmount := desiredAmount * shareToken.TokenReserve / poolAsset.TokenReserve
		if liquidityAmount > castAmount {
			liquidityAmount = castAmount
		}
	}

	if liquidityAmount == 0 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAmount, "no liquidity with amounts you deposit")
	}

	// calculate asset amounts from account to pool and update pool data
	collateral := sdk.NewCoins()
	for i, minAmount := range msg.MinAmounts {
		poolAsset, err := k.GetPoolTokenForId(ctx, msg.PoolId, uint64(i))
		if err != nil {
			return nil, err
		}

		// input token amount calculated by liquidity amount
		castAmount := liquidityAmount * poolAsset.TokenReserve / shareToken.TokenReserve

		// if input token amount is below min amount, then revert
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
				poolAsset.TokenDenom,
				sdk.NewInt(int64(castAmount)),
			),
		)

		poolAsset.TokenReserve += castAmount
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
		collateral,
	)
	if err != nil {
		return nil, err
	}

	// mint new share token, amount is liquidityAmount
	newShareToken := sdk.NewCoin(
		shareToken.TokenDenom,
		sdk.NewInt(int64(liquidityAmount)),
	)

	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(newShareToken))
	if err != nil {
		return nil, err
	}

	shareToken.TokenReserve += liquidityAmount
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

	return &types.MsgAddLiquidityResponse{
		ShareToken: newShareToken,
	}, nil
}
