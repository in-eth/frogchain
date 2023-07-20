package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// DepositBalanceKeyPrefix is the prefix to retrieve all DepositBalance
	DepositBalanceKeyPrefix = "DepositBalance/value/"
)

// DepositBalanceKey returns the store key to retrieve a DepositBalance from the index fields
func DepositBalanceKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
