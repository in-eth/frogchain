package types

import "time"

const (
	// ModuleName defines the module name
	ModuleName = "investibc"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_investibc"

	// Version defines the current version the IBC module supports
	Version = "investibc-1"

	// PortID is the default port id that module binds to
	PortID = "investibc"

	// ModuleToken defines the native token of module
	ModuleToken = "frog"
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("investibc-port-")
)

var (
	LockDuration = time.Duration(10 * time.Second)

	IBCTransferChannelID = "channel-0"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
