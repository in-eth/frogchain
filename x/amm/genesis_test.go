package amm_test

import (
	"testing"

	keepertest "frogchain/testutil/keeper"
	"frogchain/testutil/nullify"
	"frogchain/x/amm"
	"frogchain/x/amm/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		PoolList: []types.Pool{
			{
				Id: 0,
				PoolParam: types.PoolParam{
					SwapFee: sdk.NewDec(1),
					ExitFee: sdk.NewDec(1),
				},
				PoolAssets: []sdk.Coin{
					sdk.NewCoin("coin", sdk.NewInt(100)),
					sdk.NewCoin("val", sdk.NewInt(100)),
				},
				ShareToken:       sdk.NewCoin(types.ShareTokenIndex(0), sdk.NewInt(100)),
				MinimumLiquidity: sdk.NewDec(1000),
				IsActivated:      true,
			},
			{
				Id: 1,
				PoolParam: types.PoolParam{
					SwapFee: sdk.NewDec(1),
					ExitFee: sdk.NewDec(1),
				},
				PoolAssets: []sdk.Coin{
					sdk.NewCoin("coin", sdk.NewInt(100)),
					sdk.NewCoin("val", sdk.NewInt(100)),
				},
				ShareToken:       sdk.NewCoin(types.ShareTokenIndex(0), sdk.NewInt(100)),
				MinimumLiquidity: sdk.NewDec(1000),
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
