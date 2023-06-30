package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/amm module sentinel errors
var (
	ErrSample = sdkerrors.Register(ModuleName, 1100, "sample error")

	ErrFeeOverflow      = sdkerrors.Register(ModuleName, 1102, "fee amount overflow")
	ErrWeightZero       = sdkerrors.Register(ModuleName, 1103, "weight amount zero. must be greater than 0 for pool create")
	ErrInValidToken     = sdkerrors.Register(ModuleName, 1104, "token is not a valid Coins object")
	ErrInvalidAddress   = sdkerrors.Register(ModuleName, 1105, "invalid address")
	ErrInvalidAmount    = sdkerrors.Register(ModuleName, 1106, "invalid amount")
	ErrInvalidDeadline  = sdkerrors.Register(ModuleName, 1107, "deadline cannot be parsed: %s")
	ErrDeadlinePassed   = sdkerrors.Register(ModuleName, 1108, "deadline is passed: %s %s")
	ErrInvalidPath      = sdkerrors.Register(ModuleName, 1109, "invalid path for swap")
	ErrInvalidSwapDenom = sdkerrors.Register(ModuleName, 1110, "invalid swap, in and out denom is same")
)
