package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/amm module sentinel errors
var (
	ErrSample = RegisterError(ModuleName, 1100, "sample error")

	ErrFeeOverflow         = RegisterError(ModuleName, 1102, "Fee amount overflow")
	ErrWeightZero          = RegisterError(ModuleName, 1103, "Weight amount is 0")
	ErrInvalidWeightlength = RegisterError(ModuleName, 1104, "Weights and assets length are not equal")
	ErrInvalidAddress      = RegisterError(ModuleName, 1105, "Invalid address")
	ErrInvalidAmount       = RegisterError(ModuleName, 1106, "Invalid amount")
	ErrInvalidAssetsLength = RegisterError(ModuleName, 1107, "Invalid assets length")
	ErrDeadlinePassed      = RegisterError(ModuleName, 1108, "Deadline passed")
	ErrInvalidPath         = RegisterError(ModuleName, 1109, "Invalid path for swap")
	ErrInvalidSwapDenom    = RegisterError(ModuleName, 1110, "In and out denom should not be equal")
	ErrInvalidAssets       = RegisterError(ModuleName, 1111, "Invalid assets")
	ErrDuplicateAssets     = RegisterError(ModuleName, 1112, "Duplicate assets")
	ErrInvalidLength       = RegisterError(ModuleName, 1113, "Invalid assets length")
	ErrUnderMinAmount      = RegisterError(ModuleName, 1114, "Swaped value is under min amount")
	ErrInvalidSwapAmount   = RegisterError(ModuleName, 1115, "Swap amount exceeds pool token balance")
)

func RegisterError(moduleName string, errorId uint32, description string) *errors.Error {
	return sdkerrors.Register(moduleName, errorId, description)
}
