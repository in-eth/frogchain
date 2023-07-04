package keeper

import (
	"context"

	"frogchain/x/amm/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// pool data
	var pool = types.Pool{
		PoolParam:   *msg.PoolParam,
		PoolAssets:  msg.PoolAssets,
		IsActivated: true,
	}

	// append pool
	poolId := k.AppendPool(ctx, pool)

	// get message creator
	creator, _ := sdk.AccAddressFromBech32(msg.Creator)

	// castShareAmount is the cast amount for shares
	// castShareAmount = mulAll(tokenAmount)
	castShareAmount := math.NewInt(1)
	collateral := sdk.NewCoins()
	for i, assetAmount := range msg.AssetAmounts {
		collateral.Add(
			sdk.NewCoin(
				msg.PoolAssets[i].TokenDenom,
				sdk.NewInt(int64(assetAmount)),
			),
		)

		poolAsset := msg.PoolAssets[i]
		poolAsset.TokenReserve = assetAmount
		err := k.SetPoolToken(ctx, poolId, uint64(i), poolAsset)
		if err != nil {
			return nil, err
		}

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

	// set share token data
	shareToken := types.PoolToken{
		TokenDenom:   types.ShareTokenIndex(poolId),
		TokenWeight:  1,
		TokenReserve: math.NewIntFromBigInt(shareAmount).Uint64(),
	}
	err := k.SetPoolShareToken(ctx, poolId, &shareToken)
	if err != nil {
		return nil, err
	}

	// mint share token
	newShareToken := sdk.NewCoin(shareToken.TokenDenom, sdk.NewInt(shareAmount.Int64()))

	if err := k.bankKeeper.MintCoins(
		ctx, types.ModuleName, sdk.NewCoins(newShareToken),
	); err != nil {
		return nil, err
	}

	// minimun liquidity is truncated to maintain pool
	minLiquidity := sdk.NewInt(types.MINIMUM_LIQUIDITY)

	if shareAmount.Cmp(math.Int.BigInt(math.NewInt(types.MINIMUM_LIQUIDITY))) != 1 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAmount,
			"not enough coins for minimum liquidity",
		)
	}

	// creator receive share token except minimum liquidity which is for maintaining pool
	creatorShareAmount := shareAmount.Sub(shareAmount, math.Int.BigInt(minLiquidity))
	creatorShareToken := sdk.NewCoins(
		sdk.NewCoin(
			shareToken.TokenDenom,
			sdk.NewInt(creatorShareAmount.Int64()),
		),
	)

	// send share_token to the pool creator
	sdkError = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, creatorShareToken)
	if sdkError != nil {
		return nil, sdkError
	}

	return &types.MsgCreatePoolResponse{
		Id: poolId,
	}, nil
}
