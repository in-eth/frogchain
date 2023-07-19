package keeper

import (
	"frogchain/x/investibc/types"
)

var _ types.QueryServer = Keeper{}
