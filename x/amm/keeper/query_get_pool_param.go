package keeper

import (
	"context"

	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetPoolParam(goCtx context.Context, req *types.QueryGetPoolParamRequest) (*types.QueryGetPoolParamResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Process the query
	poolParam, err := k.GetPoolParamForId(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetPoolParamResponse{
		PoolParam: poolParam,
	}, nil
}
