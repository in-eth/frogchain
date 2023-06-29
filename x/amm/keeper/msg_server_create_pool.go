package keeper

import (
	"context"

	"frogchain/x/amm/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	count := k.GetPoolCount(ctx)

	// pool data
	var pool = types.Pool{
		PoolParam:   msg.PoolParam,
		PoolAssets:  msg.PoolAssets,
		IsActivated: true,
	}

	// get message creator
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	// castShareAmount is the cast amount for shares
	// castShareAmount = mulAll(tokenAmount)
	castShareAmount := math.NewInt(0)
	collateral := sdk.NewCoins()
	for i, assetAmount := range msg.AssetAmounts {
		collateral.Add(
			sdk.NewCoin(
				msg.PoolAssets[i].TokenDenom,
				sdk.NewInt(int64(assetAmount)),
			),
		)

		pool.PoolAssets[i].TokenReserve = assetAmount

		castShareAmount = castShareAmount.Mul(math.NewInt(int64(assetAmount)))
	}

	// send asset tokens from creator to module
	sdkError := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, collateral)
	if sdkError != nil {
		return nil, sdkError
	}

	// shareAmount is share token amount for liquidity
	// shareAmount ^ 2 = castShareAmount = mulAll(tokenAmount)
	shareAmount := math.Int.BigInt(castShareAmount).Sqrt(castShareAmount.BigInt())

	// mint share token
	shareToken := sdk.NewCoin(types.ShareTokenIndex(count), sdk.NewInt(shareAmount.Int64()))

	if err := k.bankKeeper.MintCoins(
		ctx, types.ModuleName, sdk.NewCoins(shareToken),
	); err != nil {
		return nil, err
	}

	// share token data
	pool.ShareToken.TokenDenom = types.ShareTokenIndex(count)
	pool.ShareToken.TokenReserve = shareAmount.Uint64()

	// minimun liquidity is truncated to maintain pool
	minLiquidity := sdk.NewInt(types.MINIMUM_LIQUIDITY)

	// creator receive share token except minimum liquidity which is for maintaining pool
	creatorShareAmount := shareAmount.Sub(shareAmount, math.Int.BigInt(minLiquidity))
	creatorShareToken := sdk.NewCoins(
		sdk.NewCoin(
			types.ShareTokenIndex(count),
			sdk.NewInt(creatorShareAmount.Int64()),
		),
	)

	// send share_token to the pool creator
	sdkError = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, creatorShareToken)
	if sdkError != nil {
		return nil, sdkError
	}

	// append pool
	poolId := k.AppendPool(ctx, pool)

	return &types.MsgCreatePoolResponse{
		Id: poolId,
	}, nil
}
