package keeper

import (
	"encoding/binary"
	"fmt"

	"frogchain/x/amm/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetPoolCount get the total number of pool
func (k Keeper) GetPoolCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.PoolCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetPoolCount set the total number of pool
func (k Keeper) SetPoolCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.PoolCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendPool appends a pool in the store with a new id and update the count
func (k Keeper) AppendPool(
	ctx sdk.Context,
	pool types.Pool,
) uint64 {
	// Create the pool
	count := k.GetPoolCount(ctx)

	// Set the ID of the appended value
	pool.Id = count

	k.SetPool(ctx, pool)

	// Update pool count
	k.SetPoolCount(ctx, count+1)

	return count
}

// SetPool set a specific pool in the store
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	b := k.cdc.MustMarshal(&pool)
	store.Set(types.PoolKey(
		fmt.Sprint(pool.Id),
	), b)
}

// GetPool returns a pool from its id
func (k Keeper) GetPool(ctx sdk.Context, id uint64) (val types.Pool, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))

	b := store.Get(types.PoolKey(
		fmt.Sprint(id),
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePool removes a pool from the store
func (k Keeper) RemovePool(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	store.Delete(types.PoolKey(
		fmt.Sprint(id),
	))
}

// GetAllPool returns all pool
func (k Keeper) GetAllPool(ctx sdk.Context) (list []types.Pool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Pool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetPoolIDBytes returns the byte representation of the ID
func GetPoolIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetPoolIDFromBytes returns ID in uint64 format from a byte array
func GetPoolIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}

// GetPoolAssetsLength returns length of assets for pool id
func (k Keeper) GetPoolAssetsLength(ctx sdk.Context, poolId uint64) (int, error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", poolId)
	}

	result := len(pool.PoolAssets)

	return result, nil
}

// GetAllPoolAssets returns all assets for pool id
func (k Keeper) GetAllPoolAssets(ctx sdk.Context, poolId uint64) ([]sdk.Coin, error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", poolId)
	}

	return pool.PoolAssets, nil
}

func (k Keeper) GetPoolTokenForId(ctx sdk.Context, poolId uint64, assetId uint64) (sdk.Coin, error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return sdk.NewCoin("null", sdk.NewInt(0)), sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", poolId)
	}

	if len(pool.PoolAssets) <= int(assetId) {
		return sdk.NewCoin("null", sdk.NewInt(0)), types.ErrInvalidLength
	}

	return pool.PoolAssets[assetId], nil
}

func (k Keeper) SetPoolToken(ctx sdk.Context, poolId uint64, assetId uint64, poolAsset sdk.Coin) error {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", poolId)
	}

	pool.PoolAssets[assetId] = poolAsset
	k.SetPool(ctx, pool)

	return nil
}

func (k Keeper) GetPoolShareTokenForId(ctx sdk.Context, poolId uint64) (sdk.Coin, error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return pool.ShareToken, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", poolId)
	}

	return pool.ShareToken, nil
}

func (k Keeper) SetPoolShareToken(ctx sdk.Context, poolId uint64, shareToken sdk.Coin) error {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", poolId)
	}

	pool.ShareToken = shareToken

	k.SetPool(ctx, pool)

	return nil
}

func (k Keeper) GetPoolParamForId(ctx sdk.Context, poolId uint64) (types.PoolParam, error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return pool.PoolParam, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", poolId)
	}

	return pool.PoolParam, nil
}

func (k Keeper) SetPoolParam(ctx sdk.Context, poolId uint64, poolParam types.PoolParam) error {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", poolId)
	}

	pool.PoolParam = poolParam
	k.SetPool(ctx, pool)

	return nil
}
