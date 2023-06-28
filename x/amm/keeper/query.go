package keeper

import (
	"frogchain/x/amm/types"
)

var _ types.QueryServer = Keeper{}
