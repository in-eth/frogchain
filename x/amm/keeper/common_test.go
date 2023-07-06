package keeper_test

import (
	"frogchain/x/amm/testutil"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	alice = testutil.Alice
	bob   = testutil.Bob
	// carol = testutil.Carol
)

func ErrorWrap(err error, format string, args ...interface{}) error {
	return sdkerrors.Wrapf(
		err,
		format,
		args,
	)
}
