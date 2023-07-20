package keeper

import (
	"context"

	"frogchain/x/investibc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DepositDenom(goCtx context.Context, req *types.QueryGetDepositDenomRequest) (*types.QueryGetDepositDenomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetDepositDenom(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetDepositDenomResponse{DepositDenom: val}, nil
}
