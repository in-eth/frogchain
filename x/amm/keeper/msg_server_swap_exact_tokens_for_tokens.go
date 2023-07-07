package keeper

import (
	"context"

	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) SwapExactTokensForTokens(goCtx context.Context, msg *types.MsgSwapExactTokensForTokens) (*types.MsgSwapExactTokensForTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get pool param
	poolParam, err := k.GetPoolParamForId(ctx, msg.PoolId)
	if err != nil {
		return nil, err
	}

	// check deadline passed
	if msg.Deadline.Before(ctx.BlockTime()) {
		return nil, sdkerrors.Wrapf(
			types.ErrDeadlinePassed,
			msg.Deadline.String(),
		)
	}

	// get token sender
	sender := sdk.MustAccAddressFromBech32(msg.Creator)

	// get token sender
	tokenReceiver := sdk.MustAccAddressFromBech32(msg.To)

	// get fee collector
	feeCollector := sdk.MustAccAddressFromBech32(poolParam.FeeCollector)

	tokenOutAmount, fee, err := k.SwapExactAmountIn(ctx, msg.PoolId, msg.AmountIn, msg.Path)
	if err != nil {
		return nil, err
	}

	// if result is below min value, then revert
	if tokenOutAmount.LT(msg.AmountOutMin) {
		return nil, types.ErrUnderMinAmount
	}

	// send fee token to fee collector
	err = k.bankKeeper.SendCoins(ctx, sender, feeCollector, sdk.NewCoins(
		sdk.NewCoin(
			msg.Path[0],
			fee.RoundInt(),
		),
	))
	if err != nil {
		return nil, err
	}

	// send input token from sender to module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(
		sdk.NewCoin(
			msg.Path[0],
			msg.AmountIn.Sub(fee).RoundInt(),
		),
	))
	if err != nil {
		return nil, err
	}

	// send output token from module to `to` receiver
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, tokenReceiver, sdk.NewCoins(
		sdk.NewCoin(
			msg.Path[len(msg.Path)-1],
			tokenOutAmount.TruncateInt(),
		),
	))

	if err != nil {
		return nil, err
	}

	// emit mint event
	ctx.EventManager().EmitEvent(
		types.NewSwapExactTokensForTokensEvent(sender, msg.PoolId, tokenOutAmount.TruncateInt().Uint64()),
	)

	return &types.MsgSwapExactTokensForTokensResponse{
		AmountOut: tokenOutAmount.TruncateInt().Uint64(),
	}, nil
}
