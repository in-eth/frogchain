package keeper

import (
	"context"

	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetPoolAsset(goCtx context.Context, req *types.QueryGetPoolAssetRequest) (*types.QueryGetPoolAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	poolAsset, err := k.GetPoolTokenForId(ctx, req.PoolId, req.AssetId)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetPoolAssetResponse{
		PoolAsset: poolAsset,
	}, nil
}
