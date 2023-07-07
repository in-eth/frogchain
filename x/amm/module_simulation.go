package amm

import (
	"math/rand"

	"frogchain/testutil/sample"
	ammsimulation "frogchain/x/amm/simulation"
	"frogchain/x/amm/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = ammsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgCreatePool          = "op_weight_msg_create_pool"
	defaultWeightMsgCreatePool int = 100

	opWeightMsgAddLiquidity          = "op_weight_msg_add_liquidity"
	defaultWeightMsgAddLiquidity int = 100

	opWeightMsgRemoveLiquidity          = "op_weight_msg_remove_liquidity"
	defaultWeightMsgRemoveLiquidity int = 100

	opWeightMsgSwapExactTokensForTokens          = "op_weight_msg_swap_exact_tokens_for_tokens"
	defaultWeightMsgSwapExactTokensForTokens int = 100

	opWeightMsgSwapTokensForExactTokens          = "op_weight_msg_swap_tokens_for_exact_tokens"
	defaultWeightMsgSwapTokensForExactTokens int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	ammGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&ammGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreatePool int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreatePool, &weightMsgCreatePool, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePool = defaultWeightMsgCreatePool
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePool,
		ammsimulation.SimulateMsgCreatePool(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAddLiquidity int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgAddLiquidity, &weightMsgAddLiquidity, nil,
		func(_ *rand.Rand) {
			weightMsgAddLiquidity = defaultWeightMsgAddLiquidity
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddLiquidity,
		ammsimulation.SimulateMsgAddLiquidity(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgRemoveLiquidity int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgRemoveLiquidity, &weightMsgRemoveLiquidity, nil,
		func(_ *rand.Rand) {
			weightMsgRemoveLiquidity = defaultWeightMsgRemoveLiquidity
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRemoveLiquidity,
		ammsimulation.SimulateMsgRemoveLiquidity(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSwapExactTokensForTokens int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSwapExactTokensForTokens, &weightMsgSwapExactTokensForTokens, nil,
		func(_ *rand.Rand) {
			weightMsgSwapExactTokensForTokens = defaultWeightMsgSwapExactTokensForTokens
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSwapExactTokensForTokens,
		ammsimulation.SimulateMsgSwapExactTokensForTokens(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSwapTokensForExactTokens int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSwapTokensForExactTokens, &weightMsgSwapTokensForExactTokens, nil,
		func(_ *rand.Rand) {
			weightMsgSwapTokensForExactTokens = defaultWeightMsgSwapTokensForExactTokens
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSwapTokensForExactTokens,
		ammsimulation.SimulateMsgSwapTokensForExactTokens(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreatePool,
			defaultWeightMsgCreatePool,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				ammsimulation.SimulateMsgCreatePool(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddLiquidity,
			defaultWeightMsgAddLiquidity,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				ammsimulation.SimulateMsgAddLiquidity(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgRemoveLiquidity,
			defaultWeightMsgRemoveLiquidity,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				ammsimulation.SimulateMsgRemoveLiquidity(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgSwapExactTokensForTokens,
			defaultWeightMsgSwapExactTokensForTokens,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				ammsimulation.SimulateMsgSwapExactTokensForTokens(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgSwapTokensForExactTokens,
			defaultWeightMsgSwapTokensForExactTokens,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				ammsimulation.SimulateMsgSwapTokensForExactTokens(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
