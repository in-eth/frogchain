package keeper

import (
	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SwapToken swaps input token to output token
// type = 1, 2
// if type == 1, input exact amount
// else if type == 2, output exact amount
func (k Keeper) SwapToken(
	ctx sdk.Context,
	poolId uint64,
	tokenAmount uint64,
	tokenDenomIn string,
	tokenDenomOut string,
	swapType uint,
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

	tokenInAmount, tokenOutAmount := uint64(0), uint64(0)
	if swapType == types.SWAP_EXACT_TOKEN_IN {
		tokenOutAmount = (reserve1 * tokenInAmount) / (reserve0 + tokenInAmount)
		tokenInAmount = tokenAmount
	} else if swapType == types.SWAP_EXACT_TOKEN_OUT {
		tokenInAmount = (reserve0 * tokenOutAmount) / (reserve1 - tokenOutAmount)
		tokenOutAmount = tokenAmount
	}

	pool.PoolAssets[inputId].TokenReserve += tokenInAmount
	pool.PoolAssets[outputId].TokenReserve -= tokenOutAmount

	k.SetPool(ctx, pool)

	if swapType == types.SWAP_EXACT_TOKEN_IN {
		return tokenOutAmount, nil
	}
	return tokenInAmount, nil
}
