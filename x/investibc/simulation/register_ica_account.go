package simulation

import (
	"math/rand"

	"frogchain/x/investibc/keeper"
	"frogchain/x/investibc/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgRegisterIcaAccount(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgRegisterIcaAccount{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the RegisterIcaAccount simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "RegisterIcaAccount simulation not implemented"), nil, nil
	}
}
