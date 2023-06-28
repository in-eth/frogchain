package keeper

import (
	"context"

	"frogchain/x/amm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddLiquidity(goCtx context.Context, msg *types.MsgAddLiquidity) (*types.MsgAddLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgAddLiquidityResponse{}, nil
}
