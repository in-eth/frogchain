package keeper

import (
	"context"

	"frogchain/x/amm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SwapTokensForExactTokens(goCtx context.Context, msg *types.MsgSwapTokensForExactTokens) (*types.MsgSwapTokensForExactTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgSwapTokensForExactTokensResponse{}, nil
}
