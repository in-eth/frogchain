package keeper

import (
	"context"

	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetPoolAssets(goCtx context.Context, req *types.QueryGetPoolAssetsRequest) (*types.QueryGetPoolAssetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	poolAssets, err := k.GetAllPoolAssets(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetPoolAssetsResponse{
		Assets: poolAssets,
	}, nil
}
