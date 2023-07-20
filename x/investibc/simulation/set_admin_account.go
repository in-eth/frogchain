package simulation

import (
	"math/rand"

	"frogchain/x/investibc/keeper"
	"frogchain/x/investibc/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgSetAdminAccount(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgSetAdminAccount{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the SetAdminAccount simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "SetAdminAccount simulation not implemented"), nil, nil
	}
}
