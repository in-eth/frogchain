package keeper

import (
	"frogchain/x/investibc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.AdminAccount(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// AdminAccount returns the AdminAccount param
func (k Keeper) AdminAccount(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyAdminAccount, &res)
	return
}
