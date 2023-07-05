package keeper

import (
	"context"

	"frogchain/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetPoolShareTokenDenom(goCtx context.Context, req *types.QueryGetPoolShareTokenDenomRequest) (*types.QueryGetPoolShareTokenDenomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Process the query
	shareToken, err := k.GetPoolShareTokenForId(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetPoolShareTokenDenomResponse{
		ShareDenom: shareToken.Denom,
	}, nil
}
