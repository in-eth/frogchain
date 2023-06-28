package keeper

import (
	"context"

	"frogchain/x/amm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgCreatePoolResponse{}, nil
}
