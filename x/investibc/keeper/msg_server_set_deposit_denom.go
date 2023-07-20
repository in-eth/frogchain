package keeper

import (
	"context"

	"frogchain/x/investibc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SetDepositDenom(goCtx context.Context, msg *types.MsgSetDepositDenom) (*types.MsgSetDepositDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	k.SetDepositDenom(ctx, types.NewDepositDenom(msg.Denom))

	return &types.MsgSetDepositDenomResponse{}, nil
}
