package keeper

import (
	"context"

	"frogchain/x/investibc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RegisterIcaAccount(goCtx context.Context, msg *types.MsgRegisterIcaAccount) (*types.MsgRegisterIcaAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get message creator
	sdk.MustAccAddressFromBech32(msg.Creator)

	if msg.Creator != k.AdminAccount(ctx) {
		return nil, types.ErrAdminPermission
	}

	if err := k.icaControllerKeeper.RegisterInterchainAccount(ctx, msg.ConnectionId, msg.Creator, msg.Version); err != nil {
		return nil, err
	}

	return &types.MsgRegisterIcaAccountResponse{}, nil
}
