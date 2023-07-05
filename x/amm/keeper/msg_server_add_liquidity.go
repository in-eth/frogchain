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

	if assetLength, _ := k.GetPoolAssetsLength(ctx, msg.PoolId); assetLength < len(msg.DesiredAmounts) {
		return nil, types.ErrInvalidLength
	}

	for i, desiredAmount := range msg.DesiredAmounts {
		poolAsset, err := k.GetPoolTokenForId(ctx, msg.PoolId, uint64(i))
		if err != nil {
			return nil, err
		}

		castAmount := desiredAmount * shareToken.Amount.Uint64() / poolAsset.Amount.Uint64()
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

		// if input token amount is below min amount, then revert
		if msg.DesiredAmounts[i] < minAmount {
			return nil, sdkerrors.Wrapf(types.ErrInvalidAmount,
				"calculated amount is below minimum, %s, %s, %s",
				fmt.Sprint(i),
				fmt.Sprint(msg.DesiredAmounts[i]),
				fmt.Sprint(minAmount),
			)
		}

		castAmount := liquidityAmount * poolAsset.Amount.Uint64() / shareToken.Amount.Uint64()

		collateral = collateral.Add(
			sdk.NewCoin(
				poolAsset.Denom,
				sdk.NewInt(int64(msg.DesiredAmounts[i])),
			),
		)

		poolAsset = poolAsset.Add(
			sdk.NewCoin(
				poolAsset.Denom,
				sdk.NewInt(int64(castAmount)),
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
		collateral,
	)
	if err != nil {
		return nil, err
	}

	// mint new share token, amount is liquidityAmount
	newShareToken := sdk.NewCoin(
		shareToken.Denom,
		sdk.NewInt(int64(liquidityAmount)),
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

	return &types.MsgAddLiquidityResponse{
		ShareToken: newShareToken,
	}, nil
}
