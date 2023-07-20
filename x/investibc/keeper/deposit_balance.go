package keeper

import (
	"frogchain/x/investibc/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetDepositBalance set a specific depositBalance in the store from its index
func (k Keeper) SetDepositBalance(ctx sdk.Context, depositBalance types.DepositBalance) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositBalanceKeyPrefix))
	b := k.cdc.MustMarshal(&depositBalance)
	store.Set(types.DepositBalanceKey(
		depositBalance.Index,
	), b)
}

// GetDepositBalance returns a depositBalance from its index
func (k Keeper) GetDepositBalance(
	ctx sdk.Context,
	index string,

) (val types.DepositBalance, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositBalanceKeyPrefix))

	b := store.Get(types.DepositBalanceKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveDepositBalance removes a depositBalance from the store
func (k Keeper) RemoveDepositBalance(
	ctx sdk.Context,
	index string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositBalanceKeyPrefix))
	store.Delete(types.DepositBalanceKey(
		index,
	))
}

// GetAllDepositBalance returns all depositBalance
func (k Keeper) GetAllDepositBalance(ctx sdk.Context) (list []types.DepositBalance) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositBalanceKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DepositBalance
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
