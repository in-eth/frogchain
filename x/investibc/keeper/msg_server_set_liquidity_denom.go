package keeper

import (
	"context"

	"frogchain/x/investibc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SetLiquidityDenom(goCtx context.Context, msg *types.MsgSetLiquidityDenom) (*types.MsgSetLiquidityDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	k.SetLiquidityDenomParam(ctx, msg.Denom)

	return &types.MsgSetLiquidityDenomResponse{}, nil
}
