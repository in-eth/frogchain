package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/investibc module sentinel errors
var (
	ErrSample                = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrInvalidPacketTimeout  = sdkerrors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion        = sdkerrors.Register(ModuleName, 1501, "invalid version")
	ErrAdminPermission       = sdkerrors.Register(ModuleName, 1502, "no permission")
	ErrInvalidDenom          = sdkerrors.Register(ModuleName, 1503, "invalid denom")
	ErrConnectionNotFound    = sdkerrors.Register(ModuleName, 1504, "connection not found")
	ErrPortIdNotFound        = sdkerrors.Register(ModuleName, 1505, "port id not found")
	ErrSwapResponseUnmarshal = sdkerrors.Register(ModuleName, 1506, "cannot unmarshal swap response message")
)
