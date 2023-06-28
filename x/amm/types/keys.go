package types

import "fmt"

const (
	// ChainName defines the chain name
	ChainName = "frogchain"

	// ModuleName defines the module name
	ModuleName = "amm"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_amm"

	// MINIMUM_LIQUIDITY defines the minimum share amount to maintain pool
	MINIMUM_LIQUIDITY = 1000
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	PoolKey      = "Pool/"
	PoolCountKey = "Pool/count/"
)

func ShareTokenIndex(portID uint64) string {
	return fmt.Sprintf("%s-%s-pool-%s-shareToken", ChainName, ModuleName, fmt.Sprint(portID))
}
