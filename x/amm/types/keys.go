package types

import (
	"fmt"
	"time"
)

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

	// TOTALPERCENT defines the 100% amount(this is for fee calc)
	TOTALPERCENT = 100000000

	// SWAP_EXACT_TOKEN_IN defines the type of token.
	// get output token amount for exact input token
	SWAP_EXACT_TOKEN_IN = 1

	// SWAP_EXACT_TOKEN_OUT defines the type of token.
	// get input token amount for exact output token
	SWAP_EXACT_TOKEN_OUT = 2
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

const (
	MaxTurnDuration = time.Duration(24 * 3_600 * 1000_000_000) // 1 day
	DeadlineLayout  = "2006-01-02 15:04:05.999999999 +0000 UTC"
)
