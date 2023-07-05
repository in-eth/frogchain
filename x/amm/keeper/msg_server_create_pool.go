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

	// castShareAmount is the cast amount for shares
	// castShareAmount = mulAll(tokenAmount)
	castShareAmount := math.NewInt(1)
	collateral := sdk.NewCoins()
	for _, asset := range msg.PoolAssets {
		collateral = collateral.Add(asset)

		castShareAmount = castShareAmount.Mul(asset.Amount)
	}

	// pool data
	var pool = types.Pool{
		PoolParam:        *msg.PoolParam,
		PoolAssets:       msg.PoolAssets,
		AssetWeights:     msg.AssetWeights,
		MinimumLiquidity: types.MINIMUM_LIQUIDITY,
		IsActivated:      true,
	}

	// append pool
	poolId := k.AppendPool(ctx, pool)

	// get message creator
	creator, _ := sdk.AccAddressFromBech32(msg.Creator)

	// send asset tokens from creator to module
	sdkError := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, collateral)
	if sdkError != nil {
		return nil, sdkError
	}

	// shareAmount is share token amount for liquidity
	// shareAmount ^ 2 = castShareAmount = mulAll(tokenAmount)
	shareAmount := math.Int.BigInt(castShareAmount).Sqrt(castShareAmount.BigInt())

	// set share token data
	shareToken := sdk.Coin{
		Denom:  types.ShareTokenIndex(poolId),
		Amount: sdk.NewInt(shareAmount.Int64()),
	}
	err := k.SetPoolShareToken(ctx, poolId, shareToken)
	if err != nil {
		return nil, err
	}

	// mint share token

	if err := k.bankKeeper.MintCoins(
		ctx, types.ModuleName, sdk.NewCoins(shareToken),
	); err != nil {
		return nil, err
	}

	// minimun liquidity is truncated to maintain pool
	minLiquidity := pool.MinimumLiquidity

	if shareAmount.Uint64() <= minLiquidity {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAmount,
			"not enough coins for minimum liquidity",
		)
	}

	// creator receive share token except minimum liquidity which is for maintaining pool
	creatorShareAmount := shareAmount.Uint64() - minLiquidity
	creatorShareToken := sdk.NewCoins(
		sdk.NewCoin(
			shareToken.Denom,
			sdk.NewInt(int64(creatorShareAmount)),
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
