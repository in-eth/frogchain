package amm_test

import (
	"testing"

	keepertest "frogchain/testutil/keeper"
	"frogchain/testutil/nullify"
	"frogchain/x/amm"
	"frogchain/x/amm/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		PoolList: []types.Pool{
			{
				Id: 0,
				PoolParam: types.PoolParam{
					SwapFee: 1,
					ExitFee: 1,
				},
				PoolAssets: []types.PoolToken{
					types.PoolToken{
						TokenDenom:   "coin",
						TokenWeight:  1,
						TokenReserve: 100,
					},
					types.PoolToken{
						TokenDenom:   "val",
						TokenWeight:  1,
						TokenReserve: 100,
					},
				},
				ShareToken: &types.PoolToken{
					TokenDenom:   types.ShareTokenIndex(0),
					TokenWeight:  1,
					TokenReserve: 100,
				},
				Mininumliquidity: 1000,
				IsActivated:      true,
			},
			{
				Id: 1,
				PoolParam: types.PoolParam{
					SwapFee: 1,
					ExitFee: 1,
				},
				PoolAssets: []types.PoolToken{
					types.PoolToken{
						TokenDenom:   "coin",
						TokenWeight:  1,
						TokenReserve: 100,
					},
					types.PoolToken{
						TokenDenom:   "val",
						TokenWeight:  1,
						TokenReserve: 100,
					},
				},
				ShareToken: &types.PoolToken{
					TokenDenom:   types.ShareTokenIndex(1),
					TokenWeight:  1,
					TokenReserve: 100,
				},
				Mininumliquidity: 1000,
				IsActivated:      true,
			},
		},
		PoolCount: 2,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.AmmKeeper(t)
	amm.InitGenesis(ctx, *k, genesisState)
	got := amm.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.PoolList, got.PoolList)
	require.Equal(t, genesisState.PoolCount, got.PoolCount)
	// this line is used by starport scaffolding # genesis/test/assert
}
