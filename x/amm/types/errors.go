package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/amm module sentinel errors
var (
	ErrSample = RegisterError(ModuleName, 1100, "sample error")

	ErrFeeOverflow       = RegisterError(ModuleName, 1102, "fee amount overflow")
	ErrWeightZero        = RegisterError(ModuleName, 1103, "weight amount zero. must be greater than 0 for pool create")
	ErrInValidToken      = RegisterError(ModuleName, 1104, "token is not a valid Coins object")
	ErrInvalidAddress    = RegisterError(ModuleName, 1105, "invalid address")
	ErrInvalidAmount     = RegisterError(ModuleName, 1106, "invalid amount")
	ErrInvalidDeadline   = RegisterError(ModuleName, 1107, "deadline cannot be parsed: %s")
	ErrDeadlinePassed    = RegisterError(ModuleName, 1108, "deadline is passed")
	ErrInvalidPath       = RegisterError(ModuleName, 1109, "invalid path for swap")
	ErrInvalidSwapDenom  = RegisterError(ModuleName, 1110, "invalid swap, in and out denom is same")
	ErrInvalidAssets     = RegisterError(ModuleName, 1111, "invalid assets")
	ErrDuplicateAssets   = RegisterError(ModuleName, 1112, "duplicate assets")
	ErrInvalidLength     = RegisterError(ModuleName, 1113, "invalid assets length")
	ErrUnderMinAmount    = RegisterError(ModuleName, 1114, "swaped value is under min amount")
	ErrInvalidSwapAmount = RegisterError(ModuleName, 1115, "swap amount is exceed pool token balance")
)

func RegisterError(moduleName string, errorId uint32, description string) *errors.Error {
	return sdkerrors.Register(moduleName, errorId, description)
}
