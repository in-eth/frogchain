package keeper

import (
	"context"
	"time"

	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) SwapTokensForExactTokens(goCtx context.Context, msg *types.MsgSwapTokensForExactTokens) (*types.MsgSwapTokensForExactTokensResponse, error) {
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
	poolParam, err := k.GetPoolParamForId(ctx, msg.PoolId)
	if err != nil {
		return nil, err
	}

	// get token sender
	sender, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	// get token receiver
	tokenReceiver, err := sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		panic(err)
	}

	// get fee collector
	feeCollector, err := sdk.AccAddressFromBech32(poolParam.FeeCollector)
	if err != nil {
		panic(err)
	}

	// send output token from module to receiver
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, tokenReceiver, sdk.NewCoins(
		sdk.NewCoin(
			msg.Path[len(msg.Path)-1],
			sdk.NewInt(int64(msg.AmountOut)),
		),
	))
	if err != nil {
		return nil, err
	}

	tokenInAmount, fee, err := k.SwapExactAmountOut(ctx, msg.PoolId, msg.AmountOut, msg.Path)

	// send fee token to fee collector
	err = k.bankKeeper.SendCoins(ctx, sender, feeCollector, sdk.NewCoins(
		sdk.NewCoin(
			msg.Path[len(msg.Path)-1],
			sdk.NewInt(int64(fee)),
		),
	))
	if err != nil {
		return nil, err
	}

	// send input token from sender to module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(
		sdk.NewCoin(
			msg.Path[len(msg.Path)-1],
			sdk.NewInt(int64(tokenInAmount)),
		),
	))
	if err != nil {
		return nil, err
	}

	return &types.MsgSwapTokensForExactTokensResponse{
		AmountIn: tokenInAmount,
	}, nil
}
