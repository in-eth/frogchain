package keeper

import (
	"context"
	"frogchain/x/investibc/types"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) BeginBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
}

// Called every block, update validator set
func (k *Keeper) EndBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	err := k.BlockSendTokensToICA(ctx)
	return err
}

func (k *Keeper) BlockSendTokensToICA(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// sendToken := k.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(types.ModuleName), denom)
	sendToken := k.CurrentDepositAmount(ctx)
	sdkError := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(k.AdminAccount(ctx)), sdk.NewCoins(sendToken))
	if sdkError != nil {
		sendToken.Amount = math.ZeroInt()
		k.SetCurrentDepositAmountParam(ctx, sendToken)
	}

	return sdkError
}
