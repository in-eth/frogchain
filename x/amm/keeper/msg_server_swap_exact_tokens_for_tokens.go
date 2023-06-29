package keeper

import (
	"context"
	"time"

	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) SwapExactTokensForTokens(goCtx context.Context, msg *types.MsgSwapExactTokensForTokens) (*types.MsgSwapExactTokensForTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message

	// check deadline is passed
	deadline, err := time.Parse(types.DeadlineLayout, msg.Deadline)
	if err != nil {
		return nil, err
	}

	if ctx.BlockTime().After(deadline) {
		return nil, sdkerrors.Wrapf(
			types.ErrDeadlinePassed,
			types.ErrDeadlinePassed.Error(),
			deadline.UTC().Format(types.DeadlineLayout),
			ctx.BlockTime().UTC().Format(types.DeadlineLayout),
		)
	}

	// get pool param
	poolParam, err := k.GetPoolParam(ctx, msg.PoolId)
	if err != nil {
		return nil, err
	}

	// calc fee and send it to feeCollector
	fee := msg.AmountIn * poolParam.SwapFee / types.TOTALPERCENT
	tokenOutAmount := msg.AmountIn - fee

	// get token sender
	sender, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	// get fee collector
	feeCollector, err := sdk.AccAddressFromBech32(poolParam.FeeCollector)
	if err != nil {
		panic(err)
	}

	// send fee token to fee collector
	err = k.bankKeeper.SendCoins(ctx, sender, feeCollector, sdk.NewCoins(
		sdk.NewCoin(
			msg.Path[0],
			sdk.NewInt(int64(fee)),
		),
	))
	if err != nil {
		return nil, err
	}

	// send input token from sender to module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(
		sdk.NewCoin(
			msg.Path[0],
			sdk.NewInt(int64(tokenOutAmount)),
		),
	))
	if err != nil {
		return nil, err
	}

	for i, tokenDenomIn := range msg.Path {
		if len(msg.Path)-1 == i {
			break
		}

		tokenDenomOut := msg.Path[i+1]

		tokenOutAmount, err = k.SwapToken(ctx, msg.PoolId, tokenOutAmount, tokenDenomIn, tokenDenomOut)
		if err != nil {
			return nil, err
		}

		tokenDenomIn = tokenDenomOut
	}

	// send output token from module to sender
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(
		sdk.NewCoin(
			msg.Path[len(msg.Path)-1],
			sdk.NewInt(int64(tokenOutAmount)),
		),
	))
	if err != nil {
		return nil, err
	}

	return &types.MsgSwapExactTokensForTokensResponse{
		AmountOut: tokenOutAmount,
	}, nil
}
