package keeper

import (
	"context"

	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetPoolTotalLiquidity(goCtx context.Context, req *types.QueryGetPoolTotalLiquidityRequest) (*types.QueryGetPoolTotalLiquidityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Process the query
	shareToken, err := k.GetPoolShareTokenForId(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetPoolTotalLiquidityResponse{
		TotalLiquidity: shareToken.TokenReserve,
	}, nil
}
