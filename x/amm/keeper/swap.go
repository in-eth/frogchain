package keeper

import (
	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SwapExactAmountIn(
	ctx sdk.Context,
	poolId uint64,
	tokenInAmount uint64,
	path []string,
) (uint64, uint64, error) {
	// get pool params
	poolParam, err := k.GetPoolParamForId(ctx, poolId)
	if err != nil {
		return 0, 0, err
	}

	// calc fee and send it to feeCollector
	fee := tokenInAmount * poolParam.SwapFee / types.TOTALPERCENT
	tokenInAmount -= fee

	tokenOutAmount := tokenInAmount
	for i, tokenDenomIn := range path {
		if len(path)-1 == i {
			break
		}

		tokenDenomOut := path[i+1]

		tokenOutAmount, err = k.SwapToken(ctx, poolId, tokenOutAmount, tokenDenomIn, tokenDenomOut, types.SWAP_EXACT_TOKEN_IN)
		if err != nil {
			return 0, 0, err
		}
	}

	return tokenOutAmount, fee, err
}

func (k Keeper) SwapExactAmountOut(
	ctx sdk.Context,
	poolId uint64,
	tokenOutAmount uint64,
	path []string,
) (uint64, uint64, error) {
	poolParam, err := k.GetPoolParamForId(ctx, poolId)
	if err != nil {
		return 0, 0, err
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	tokenInAmount := tokenOutAmount
	for i, tokenDenomOut := range path {
		if i >= len(path)-1 {
			break
		}

		tokenDenomIn := path[i+1]

		tokenInAmount, err = k.SwapToken(ctx, poolId, tokenInAmount, tokenDenomIn, tokenDenomOut, types.SWAP_EXACT_TOKEN_OUT)
		if err != nil {
			return 0, 0, err
		}
	}

	// calc fee and send it to feeCollector
	fee := tokenInAmount * poolParam.SwapFee / (types.TOTALPERCENT - poolParam.SwapFee)
	tokenInAmount += fee

	return tokenInAmount, fee, err
}

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
	poolAssets, err := k.GetAllPoolAssets(
		ctx,
		poolId,
	)
	if err != nil {
		return 0, err
	}

	if tokenDenomIn == tokenDenomOut {
		return 0, types.ErrInvalidSwapDenom
	}

	inputId, outputId := int(-1), int(-1)
	for i, poolAsset := range poolAssets {
		if poolAsset.Denom == tokenDenomIn {
			inputId = i
		} else if poolAsset.Denom == tokenDenomOut {
			outputId = i
		}
		if inputId != -1 && outputId != -1 {
			break
		}
	}

	if inputId == -1 || outputId == -1 {
		return 0, types.ErrInvalidPath
	}

	reserve0 := poolAssets[inputId].Amount.Uint64()
	reserve1 := poolAssets[outputId].Amount.Uint64()

	tokenInAmount, tokenOutAmount := uint64(0), uint64(0)
	if swapType == types.SWAP_EXACT_TOKEN_IN {
		tokenInAmount = tokenAmount
		tokenOutAmount = (reserve1 * tokenInAmount) / (reserve0 + tokenInAmount)
	} else if swapType == types.SWAP_EXACT_TOKEN_OUT {
		tokenOutAmount = tokenAmount
		if tokenOutAmount >= reserve1 {
			return 0, types.ErrInvalidSwapAmount
		}
		tokenInAmount = (reserve0 * tokenOutAmount) / (reserve1 - tokenOutAmount)
	}

	inputAsset := poolAssets[inputId].AddAmount(sdk.NewInt(int64(tokenInAmount)))
	outputAsset := poolAssets[outputId].SubAmount(sdk.NewInt(int64(tokenOutAmount)))

	k.SetPoolToken(ctx, poolId, uint64(inputId), inputAsset)
	k.SetPoolToken(ctx, poolId, uint64(outputId), outputAsset)

	if swapType == types.SWAP_EXACT_TOKEN_IN {
		return tokenOutAmount, nil
	}
	return tokenInAmount, nil
}
