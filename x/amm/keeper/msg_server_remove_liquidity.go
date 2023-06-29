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

	// get pool param
	poolParam, err := k.GetPoolParam(ctx, msg.PoolId)
	if err != nil {
		return nil, err
	}

	// get pool share token
	shareToken, err := k.GetPoolShareToken(ctx, msg.PoolId)
	if err != nil {
		return nil, err
	}

	liquidityProvider, _ := sdk.AccAddressFromBech32(msg.Creator)

	// burn share token, amount is desiredAmount
	fee := msg.Liquidity * poolParam.ExitFee / (types.TOTALPERCENT)
	liquidity := msg.Liquidity - fee

	burnShareToken := sdk.NewCoin(
		types.ShareTokenIndex(msg.PoolId),
		sdk.NewInt(int64(liquidity)),
	)

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
			shareToken.TokenDenom,
			sdk.NewInt(int64(fee)),
		),
	))
	if err != nil {
		return nil, err
	}

	// update pool share token data
	shareToken.TokenReserve -= liquidity
	err = k.SetPoolShareToken(ctx, msg.PoolId, shareToken)
	if err != nil {
		return nil, err
	}

	// send assets from pool to account
	receiveTokens := sdk.NewCoins()
	for i, minAmount := range msg.MinAmounts {
		// get pool asset
		poolAsset, err := k.GetPoolToken(ctx, msg.PoolId, uint64(i))
		if err != nil {
			return nil, err
		}

		// calculate token amount for liquidity
		castAmount := liquidity * poolAsset.TokenReserve / shareToken.TokenReserve

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
				poolAsset.TokenDenom,
				math.NewInt(int64(castAmount)),
			),
		)

		// update pool asset data
		poolAsset.TokenReserve -= castAmount
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

	return &types.MsgRemoveLiquidityResponse{
		ReceivedTokens: receiveTokens,
	}, nil
}
