package keeper

import (
	"context"

	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get message creator
	creator := sdk.MustAccAddressFromBech32(msg.Creator)

	sdk.MustAccAddressFromBech32(msg.PoolParam.FeeCollector)

	if msg.PoolParam.SwapFee.GTE(sdk.NewDec(types.TOTALPERCENT)) || msg.PoolParam.ExitFee.GTE(sdk.NewDec(types.TOTALPERCENT)) {
		return nil, types.ErrFeeOverflow
	}

	if len(msg.PoolAssets) < 2 {
		return nil, types.ErrInvalidAssetsLength
	}

	if len(msg.PoolAssets) != len(msg.AssetWeights) {
		return nil, types.ErrInvalidWeightlength
	}

	// shareAmount is share token amount for liquidity
	// Pow(shareAmount, 2) = castShareAmount = mulAll(tokenAmount)
	castShareAmount := sdk.OneDec()
	assetCoins := sdk.NewCoins()
	for _, asset := range msg.PoolAssets {
		assetCoins = assetCoins.Add(asset)
		castShareAmount = castShareAmount.MulInt(asset.Amount)
	}
	shareAmount, err := castShareAmount.ApproxSqrt()
	if err != nil {
		return nil, err
	}

	// pool data
	var pool = types.Pool{
		PoolParam:        *msg.PoolParam,
		PoolAssets:       msg.PoolAssets,
		AssetWeights:     msg.AssetWeights,
		MinimumLiquidity: sdk.NewDec(types.MINIMUM_LIQUIDITY),
		IsActivated:      true,
	}

	poolId := k.GetPoolCount(ctx)

	// send asset tokens from creator to module
	sdkError := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, assetCoins)
	if sdkError != nil {
		return nil, sdkError
	}

	// set share token data
	shareToken := sdk.Coin{
		Denom:  types.ShareTokenIndex(poolId),
		Amount: shareAmount.TruncateInt(),
	}

	pool.ShareToken = shareToken

	// mint share token

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(shareToken)); err != nil {
		return nil, err
	}

	// minimun liquidity is truncated from creator's share token to maintain pool
	if shareAmount.LTE(pool.MinimumLiquidity) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAmount,
			"not enough coins for minimum liquidity",
		)
	}

	shareAmount = shareAmount.Sub(pool.MinimumLiquidity)
	creatorShareToken := sdk.NewCoins(
		sdk.NewCoin(
			shareToken.Denom,
			shareAmount.TruncateInt(),
		),
	)

	// send share_token to the pool creator
	sdkError = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, creatorShareToken)
	if sdkError != nil {
		return nil, sdkError
	}

	// append pool
	k.SetPool(ctx, pool)
	k.SetPoolCount(ctx, poolId+1)

	// emit mint event
	ctx.EventManager().EmitEvent(
		types.NewCreatePoolEvent(creator, poolId, creatorShareToken[0]),
	)

	return &types.MsgCreatePoolResponse{
		Id: poolId,
	}, nil
}
