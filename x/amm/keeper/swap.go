package keeper

import (
	"frogchain/x/amm/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SwapExactAmountIn(
	ctx sdk.Context,
	poolId uint64,
	tokenInAmount sdk.Dec,
	path []string,
) (sdk.Dec, sdk.Dec, error) {
	// get pool params
	poolParam, err := k.GetPoolParamForId(ctx, poolId)
	if err != nil {
		return sdk.NewDec(0), sdk.NewDec(0), err
	}

	// calc fee and send it to feeCollector
	fee := tokenInAmount.Mul(poolParam.SwapFee).QuoRoundUp(sdk.NewDec(types.TOTALPERCENT))
	tokenInAmount = tokenInAmount.Sub(fee)

	tokenOutAmount := tokenInAmount
	for i, tokenDenomIn := range path {
		if len(path)-1 == i {
			break
		}

		tokenDenomOut := path[i+1]

		tokenOutAmount, err = k.SwapToken(ctx, poolId, tokenOutAmount, tokenDenomIn, tokenDenomOut, types.SWAP_EXACT_TOKEN_IN)
		if err != nil {
			return sdk.NewDec(0), sdk.NewDec(0), err
		}
	}

	return tokenOutAmount, fee, err
}

func (k Keeper) SwapExactAmountOut(
	ctx sdk.Context,
	poolId uint64,
	tokenOutAmount sdk.Dec,
	path []string,
) (sdk.Dec, sdk.Dec, error) {
	poolParam, err := k.GetPoolParamForId(ctx, poolId)
	if err != nil {
		return sdk.NewDec(0), sdk.NewDec(0), err
	}

	for i := len(path) - 1; i > 0; i-- {
		tokenDenomOut := path[i]
		tokenDenomIn := path[i-1]

		tokenInAmount, err := k.SwapToken(ctx, poolId, tokenOutAmount, tokenDenomIn, tokenDenomOut, types.SWAP_EXACT_TOKEN_OUT)
		if err != nil {
			return sdk.NewDec(0), sdk.NewDec(0), err
		}

		tokenOutAmount = tokenInAmount
	}

	// calc fee and send it to feeCollector
	fee := tokenOutAmount.Mul(poolParam.SwapFee).QuoRoundUp(sdk.NewDec(types.TOTALPERCENT).Sub(poolParam.SwapFee))

	tokenInAmount := tokenOutAmount.Add(math.LegacyDec(fee.RoundInt()))

	return tokenInAmount, fee, err
}

// SwapToken swaps input token to output token
// type = 1, 2
// if type == 1, input exact amount
// else if type == 2, output exact amount
func (k Keeper) SwapToken(
	ctx sdk.Context,
	poolId uint64,
	tokenAmount sdk.Dec,
	tokenDenomIn string,
	tokenDenomOut string,
	swapType uint,
) (sdk.Dec, error) {
	poolAssets, err := k.GetAllPoolAssets(
		ctx,
		poolId,
	)
	if err != nil {
		return sdk.NewDec(0), err
	}

	if tokenDenomIn == tokenDenomOut {
		return sdk.NewDec(0), types.ErrInvalidSwapDenom
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
		return sdk.NewDec(0), types.ErrInvalidPath
	}

	reserve0 := sdk.NewDecFromInt(poolAssets[inputId].Amount)
	reserve1 := sdk.NewDecFromInt(poolAssets[outputId].Amount)

	tokenInAmount, tokenOutAmount := sdk.NewDec(0), sdk.NewDec(0)
	if swapType == types.SWAP_EXACT_TOKEN_IN {
		tokenInAmount = tokenAmount
		tokenOutAmount = reserve1.Mul(tokenInAmount).Quo(reserve0.Add(tokenInAmount))
	} else if swapType == types.SWAP_EXACT_TOKEN_OUT {
		tokenOutAmount = tokenAmount
		if tokenOutAmount.GTE(reserve1) {
			return sdk.NewDec(0), types.ErrInvalidSwapAmount
		}
		tokenInAmount = reserve0.Mul(tokenOutAmount).Quo(reserve1.Sub(tokenOutAmount))
	}

	inputAsset := poolAssets[inputId].AddAmount(tokenInAmount.RoundInt())
	outputAsset := poolAssets[outputId].SubAmount(tokenOutAmount.TruncateInt())

	k.SetPoolToken(ctx, poolId, uint64(inputId), inputAsset)
	k.SetPoolToken(ctx, poolId, uint64(outputId), outputAsset)

	if swapType == types.SWAP_EXACT_TOKEN_IN {
		return tokenOutAmount, nil
	}
	return tokenInAmount, nil
}
