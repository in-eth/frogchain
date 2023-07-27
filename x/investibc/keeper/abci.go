package keeper

import (
	"context"
	"frogchain/x/investibc/types"
	"time"

	osmosislockup "frogchain/osmosis/lockup"
	osmosispool "frogchain/osmosis/pool"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/gogoproto/proto"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"

	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
)

func (k *Keeper) BeginBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
}

// Called every block, update validator set
func (k *Keeper) EndBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	err := k.EndBlockerTransferRewards(ctx)
	if err != nil {
		return err
	}
	err = k.EndBlockerSendPacketToOsmosis(ctx)
	return err
}

func (k *Keeper) EndBlockerTransferRewards(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	portID, err := icatypes.NewControllerPortID(types.ModuleName)
	if err != nil {
		return err
	}

	addr, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, k.IcaConnectionId(ctx), portID)
	if !found {
		return types.ErrPortIdNotFound
	}

	balance := k.bankKeeper.GetBalance(ctx, sdk.AccAddress(addr), k.DepositDenom(ctx)).Sub(k.CurrentDepositAmount(ctx))

	if balance.Amount.LTE(math.ZeroInt()) {
		return nil
	}

	channelID, found := k.icaControllerKeeper.GetActiveChannelID(ctx, k.IcaConnectionId(ctx), ibctransfertypes.PortID)
	if !found {
		return icatypes.ErrActiveChannelNotFound.Wrapf("failed to retrieve active channel for port %s", portID)
	}

	// transfer tokens to Osmosis network
	timeoutTimestamp := ctx.BlockTime().Add(time.Minute).UnixNano()
	_, err = k.TransferKeeper.Transfer(goCtx, ibctransfertypes.NewMsgTransfer(
		ibctransfertypes.PortID,
		channelID,
		balance,
		addr,
		sdk.AccAddress(authtypes.FeeCollectorName).String(),
		clienttypes.Height{},
		uint64(timeoutTimestamp),
		""))
	if err != nil {
		return err
	}
	return nil
}

func (k *Keeper) EndBlockerSendPacketToOsmosis(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	portID, err := icatypes.NewControllerPortID(types.ModuleName)
	if err != nil {
		return err
	}

	channelID, found := k.icaControllerKeeper.GetActiveChannelID(ctx, k.IcaConnectionId(ctx), portID)
	if !found {
		return icatypes.ErrActiveChannelNotFound.Wrapf("failed to retrieve active channel for port %s", portID)
	}

	addr, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, k.IcaConnectionId(ctx), portID)
	if !found {
		return types.ErrPortIdNotFound
	}

	k.UnlockLiquidity(ctx, portID, channelID, addr)

	sendToken := k.CurrentDepositAmount(ctx)
	if k.JoinSwapExactAmountInPacketSent(ctx) == true || k.LockTokensPacketSent(ctx) == true {
		return nil
	}
	if k.DepositLastTime(ctx)+uint64(time.Minute) < uint64(ctx.BlockTime().Unix()) &&
		sendToken.Amount.GT(math.ZeroInt()) &&
		sendToken.Denom != "deposit_denom" {

		sdkError := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(addr), sdk.NewCoins(sendToken))
		if sdkError != nil {
			return sdkError
		}

		k.JoinSwapExactAmountIn(ctx, portID, channelID, addr, sendToken)
		return nil
	}

	liquidityToken := k.CurrentLiquidityAmount(ctx)
	if liquidityToken.Amount.GT(math.ZeroInt()) && liquidityToken.Denom != "liquidity_denom" {
		k.LockUpLiquidity(ctx, portID, channelID, addr, liquidityToken)
	}

	return nil
}

func (k *Keeper) JoinSwapExactAmountIn(ctx sdk.Context, portID string, channelID string, addr string, sendToken sdk.Coin) error {
	chanCap, found := k.IBCScopperKeeper.GetCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if !found {
		return channeltypes.ErrChannelCapabilityNotFound.Wrap("module does not own channel capability")
	}

	msg := osmosispool.MsgJoinSwapExternAmountIn{
		Sender:            addr,
		PoolId:            1,
		TokenIn:           sendToken,
		ShareOutMinAmount: math.ZeroInt(),
	}

	data, err := icatypes.SerializeCosmosTx(k.cdc, []proto.Message{&msg})
	if err != nil {
		return err
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	timeoutTimestamp := ctx.BlockTime().Add(time.Minute).UnixNano()
	_, err = k.icaControllerKeeper.SendTx(ctx, chanCap, k.IcaConnectionId(ctx), portID, packetData, uint64(timeoutTimestamp))
	if err != nil {
		return err
	}

	k.Logger(ctx).Debug("message send packet endblocker", "joinswapexactamountin", packetData)

	k.SetJoinSwapExactAmountInPacketSentParam(ctx, true)

	return nil
}

func (k *Keeper) LockUpLiquidity(ctx sdk.Context, portID string, channelID string, addr string, sendToken sdk.Coin) error {
	chanCap, found := k.IBCScopperKeeper.GetCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if !found {
		return channeltypes.ErrChannelCapabilityNotFound.Wrap("module does not own channel capability")
	}

	msg := osmosislockup.MsgLockTokens{
		Owner:    addr,
		Duration: time.Duration(10 * time.Second),
		Coins:    sdk.NewCoins(sendToken),
	}

	data, err := icatypes.SerializeCosmosTx(k.cdc, []proto.Message{&msg})
	if err != nil {
		return err
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	timeoutTimestamp := ctx.BlockTime().Add(time.Minute).UnixNano()
	_, err = k.icaControllerKeeper.SendTx(ctx, chanCap, k.IcaConnectionId(ctx), portID, packetData, uint64(timeoutTimestamp))
	if err != nil {
		return err
	}

	k.Logger(ctx).Debug("message send packet endblocker", "lockupliquidity", packetData)

	k.SetLockTokensPacketSentParam(ctx, true)

	return nil
}

func (k *Keeper) UnlockLiquidity(ctx sdk.Context, portID string, channelID string, addr string) error {
	chanCap, found := k.IBCScopperKeeper.GetCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if !found {
		return channeltypes.ErrChannelCapabilityNotFound.Wrap("module does not own channel capability")
	}

	msg := osmosislockup.MsgBeginUnlockingAll{
		Owner: addr,
	}

	data, err := icatypes.SerializeCosmosTx(k.cdc, []proto.Message{&msg})
	if err != nil {
		return err
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	timeoutTimestamp := ctx.BlockTime().Add(time.Minute).UnixNano()
	_, err = k.icaControllerKeeper.SendTx(ctx, chanCap, k.IcaConnectionId(ctx), portID, packetData, uint64(timeoutTimestamp))
	if err != nil {
		return err
	}

	k.Logger(ctx).Debug("message send packet endblocker", "unlock liquidity", packetData)

	k.SetLockTokensPacketSentParam(ctx, true)

	return nil
}
