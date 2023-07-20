package keeper

import (
	"context"

	"frogchain/x/investibc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SetDepositDenom(goCtx context.Context, msg *types.MsgSetDepositDenom) (*types.MsgSetDepositDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get message creator
	sdk.MustAccAddressFromBech32(msg.Creator)

	if msg.Creator != k.AdminAccount(ctx) {
		return nil, types.ErrAdminPermission
	}

	k.SetDepositDenomStore(ctx, types.DepositDenom{Denom: msg.Denom})

	return &types.MsgSetDepositDenomResponse{}, nil
}
