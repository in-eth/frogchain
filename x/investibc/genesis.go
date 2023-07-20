package investibc

import (
	"frogchain/x/investibc/keeper"
	"frogchain/x/investibc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set if defined
	if genState.DepositDenom != nil {
		k.SetDepositDenomStore(ctx, *genState.DepositDenom)
	}
	// Set all the depositBalance
	for _, elem := range genState.DepositBalanceList {
		k.SetDepositBalance(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetPort(ctx, genState.PortId)
	// Only try to bind to port if it is not already bound, since we may already own
	// port capability from capability InitGenesis
	if !k.IsBound(ctx, genState.PortId) {
		// module binds to the port on InitChain
		// and claims the returned capability
		err := k.BindPort(ctx, genState.PortId)
		if err != nil {
			panic("could not claim port capability: " + err.Error())
		}
	}
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.PortId = k.GetPort(ctx)
	// Get all depositDenom
	depositDenom, found := k.GetDepositDenom(ctx)
	if found {
		genesis.DepositDenom = &depositDenom
	}
	genesis.DepositBalanceList = k.GetAllDepositBalance(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
