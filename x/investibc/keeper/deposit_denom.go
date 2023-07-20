package keeper

import (
	"frogchain/x/investibc/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetDepositDenom set depositDenom in the store
func (k Keeper) SetDepositDenomStore(ctx sdk.Context, depositDenom types.DepositDenom) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositDenomKey))
	b := k.cdc.MustMarshal(&depositDenom)
	store.Set([]byte{0}, b)
}

// GetDepositDenom returns depositDenom
func (k Keeper) GetDepositDenom(ctx sdk.Context) (val types.DepositDenom, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositDenomKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveDepositDenom removes depositDenom from the store
func (k Keeper) RemoveDepositDenom(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositDenomKey))
	store.Delete([]byte{0})
}
