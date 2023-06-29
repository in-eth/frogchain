package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SwapToken swaps input token to output token
func (k Keeper) SwapToken(
	ctx sdk.Context,
	poolId uint64,
	tokenInAmount uint64,
	tokenDenomIn string,
	tokenDenomOut string,
) (uint64, error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", poolId)
	}

	inputId, outputId := int(0), int(0)
	for i, poolAsset := range pool.PoolAssets {
		if poolAsset.TokenDenom == tokenDenomIn {
			inputId = i
		} else if poolAsset.TokenDenom == tokenDenomOut {
			outputId = i
		}
	}

	reserve0 := pool.PoolAssets[inputId].TokenReserve
	reserve1 := pool.PoolAssets[outputId].TokenReserve

	tokenOutAmount := reserve1 - (reserve0*reserve1)/(reserve0+tokenInAmount)

	pool.PoolAssets[inputId].TokenReserve += tokenInAmount
	pool.PoolAssets[outputId].TokenReserve -= tokenOutAmount

	k.SetPool(ctx, pool)

	return tokenOutAmount, nil
}
