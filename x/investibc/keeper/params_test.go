package keeper_test

import (
	"testing"

	testkeeper "frogchain/testutil/keeper"
	"frogchain/x/investibc/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.InvestibcKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
	require.EqualValues(t, params.AdminAccount, k.AdminAccount(ctx))
}
