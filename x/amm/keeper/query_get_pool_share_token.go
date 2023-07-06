package keeper

import (
	"context"

	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetPoolShareToken(goCtx context.Context, req *types.QueryGetPoolShareTokenRequest) (*types.QueryGetPoolShareTokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	shareToken, err := k.GetPoolShareTokenForId(ctx, req.PoolId)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetPoolShareTokenResponse{
		ShareToken: shareToken,
	}, nil
}
