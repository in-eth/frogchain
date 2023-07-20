package keeper

import (
	"context"

	"frogchain/x/investibc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SetAdminAccount(goCtx context.Context, msg *types.MsgSetAdminAccount) (*types.MsgSetAdminAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	k.SetParams(ctx, types.NewParams(msg.AdminAccount))

	return &types.MsgSetAdminAccountResponse{}, nil
}
