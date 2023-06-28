package keeper

import (
	"context"

	"frogchain/x/amm/types"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	count := k.GetPoolCount(ctx)

	// get message creator
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	// castShareAmount is the cast amount for shares
	// castShareAmount = mulAll(tokenAmount)
	castShareAmount := math.NewInt(0)
	for i, assetAmount := range msg.AssetAmounts {
		collateral := sdk.NewCoins(
			sdk.NewCoin(
				msg.PoolAssets[i].TokenDenom,
				sdk.NewInt(int64(assetAmount)),
			),
		)
		sdkError := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, collateral)
		if sdkError != nil {
			return nil, sdkError
		}

		castShareAmount = castShareAmount.Mul(math.NewInt(int64(assetAmount)))
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
	sdkError := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, creatorShareToken)
	if sdkError != nil {
		return nil, sdkError
	}

	// pool data
	var pool = types.Pool{
		Id:          count,
		PoolParam:   msg.PoolParam,
		PoolAssets:  msg.PoolAssets,
		ShareToken:  shareToken,
		IsActivated: true,
	}

	// store pool in pool id
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKey))
	appendedValue := k.cdc.MustMarshal(&pool)
	store.Set(GetPoolIDBytes(pool.Id), appendedValue)
	k.SetPoolCount(ctx, count+1)

	return &types.MsgCreatePoolResponse{
		Id: int32(pool.Id),
	}, nil
}
