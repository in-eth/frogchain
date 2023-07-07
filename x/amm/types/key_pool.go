package types

import (
	"encoding/binary"
	"fmt"
)

var _ binary.ByteOrder

const (
	// PoolKeyPrefix is the prefix to retrieve all Pool
	PoolKeyPrefix = "Pool/value/"
)

// PoolKey returns the store key to retrieve a Pool from the index fields

func PoolKey(poolId uint64) []byte {
	var key []byte

	indexBytes := []byte(PoolKeyPrefix)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)
	key = append(key, fmt.Sprintf("%d", poolId)...)

	return key
}
