package investibc_test

import (
	"testing"

	keepertest "frogchain/testutil/keeper"
	"frogchain/testutil/nullify"
	"frogchain/x/investibc"
	"frogchain/x/investibc/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		DepositDenom: &types.DepositDenom{
			Denom: "6",
		},
		DepositBalanceList: []types.DepositBalance{
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.InvestibcKeeper(t)
	investibc.InitGenesis(ctx, *k, genesisState)
	got := investibc.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	require.Equal(t, genesisState.DepositDenom, got.DepositDenom)
	require.ElementsMatch(t, genesisState.DepositBalanceList, got.DepositBalanceList)
	// this line is used by starport scaffolding # genesis/test/assert
}
