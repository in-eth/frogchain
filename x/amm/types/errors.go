package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/amm module sentinel errors
var (
	ErrSample = sdkerrors.Register(ModuleName, 1100, "sample error")

	ErrFeeOverflow    = sdkerrors.Register(ModuleName, 1102, "fee amount overflow. must be less than 10^8")
	ErrWeightZero     = sdkerrors.Register(ModuleName, 1103, "weight amount zero. must be greater than 0 for pool create")
	ErrInValidToken   = sdkerrors.Register(ModuleName, 1104, "token is not a valid Coins object")
	ErrInvalidAddress = sdkerrors.Register(ModuleName, 1105, "invalid address")
	ErrInvalidAmount  = sdkerrors.Register(ModuleName, 1106, "invalid amount")
)
