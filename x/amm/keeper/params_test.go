package keeper_test

import (
	"testing"

	testkeeper "frogchain/testutil/keeper"
	"frogchain/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.AmmKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
