package keeper

import (
	"context"

	"frogchain/x/amm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RemoveLiquidity(goCtx context.Context, msg *types.MsgRemoveLiquidity) (*types.MsgRemoveLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgRemoveLiquidityResponse{}, nil
}
