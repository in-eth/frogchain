package keeper

import (
	"context"

	"frogchain/x/investibc/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DepositBalanceAll(goCtx context.Context, req *types.QueryAllDepositBalanceRequest) (*types.QueryAllDepositBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var depositBalances []types.DepositBalance
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	depositBalanceStore := prefix.NewStore(store, types.KeyPrefix(types.DepositBalanceKeyPrefix))

	pageRes, err := query.Paginate(depositBalanceStore, req.Pagination, func(key []byte, value []byte) error {
		var depositBalance types.DepositBalance
		if err := k.cdc.Unmarshal(value, &depositBalance); err != nil {
			return err
		}

		depositBalances = append(depositBalances, depositBalance)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDepositBalanceResponse{DepositBalance: depositBalances, Pagination: pageRes}, nil
}

func (k Keeper) DepositBalance(goCtx context.Context, req *types.QueryGetDepositBalanceRequest) (*types.QueryGetDepositBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetDepositBalance(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetDepositBalanceResponse{DepositBalance: val}, nil
}
