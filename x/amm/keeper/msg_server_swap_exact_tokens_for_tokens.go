package keeper

import (
	"context"
	"time"

	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
		return nil, ErrorWrap(
			types.ErrDeadlinePassed,
			deadline.UTC().Format(types.DeadlineLayout),
		)
	}

	// get pool param
	poolParam, err := k.GetPoolParamForId(ctx, msg.PoolId)
	if err != nil {
		return nil, err
	}

	// get token sender
	sender, _ := sdk.AccAddressFromBech32(msg.Creator)

	// get token sender
	tokenReceiver, _ := sdk.AccAddressFromBech32(msg.To)

	// get fee collector
	feeCollector, _ := sdk.AccAddressFromBech32(poolParam.FeeCollector)

	tokenOutAmount, fee, err := k.SwapExactAmountIn(ctx, msg.PoolId, msg.AmountIn, msg.Path)
	if err != nil {
		return nil, err
	}

	// send fee token to fee collector
	if fee > 0 {
		err = k.bankKeeper.SendCoins(ctx, sender, feeCollector, sdk.NewCoins(
			sdk.NewCoin(
				msg.Path[0],
				sdk.NewInt(int64(fee)),
			),
		))
		if err != nil {
			return nil, err
		}
	}

	// send input token from sender to module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(
		sdk.NewCoin(
			msg.Path[0],
			sdk.NewInt(int64(msg.AmountIn-fee)),
		),
	))
	if err != nil {
		return nil, err
	}

	// if result is below min value, then revert
	if tokenOutAmount < msg.AmountOutMin {
		return nil, types.ErrUnderMinAmount
	}

	// send output token from module to `to` receiver
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, tokenReceiver, sdk.NewCoins(
		sdk.NewCoin(
			msg.Path[len(msg.Path)-1],
			sdk.NewInt(int64(tokenOutAmount)),
		),
	))

	if err != nil {
		return nil, err
	}

	// emit mint event
	ctx.EventManager().EmitEvent(
		types.NewSwapExactTokensForTokensEvent(sender, msg.PoolId, tokenOutAmount),
	)

	return &types.MsgSwapExactTokensForTokensResponse{
		AmountOut: tokenOutAmount,
	}, nil
}
