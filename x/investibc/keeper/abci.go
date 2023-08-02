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
func (k *Keeper) EndBlocker(goCtx context.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	portId, icaAddr, err := k.GetPortICAAddr(goCtx)
	if err != nil {
		return
	}

	k.EndBlockerTransferDepositTokenToICA(goCtx, icaAddr)
	k.EndBlockerJoinSwapExactAmountIn(goCtx, portId, icaAddr)
	k.EndBlockerLockUpLiquidity(goCtx, portId, icaAddr)
	k.EndBlockerUnLockLiquidity(goCtx, portId, icaAddr)
	k.EndBlockerTransferRewards(goCtx, icaAddr)
}

func (k *Keeper) GetPortICAAddr(goCtx context.Context) (string, string, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	portID, err := icatypes.NewControllerPortID(types.ModuleName)
	if err != nil {
		return "", "", err
	}

	addr, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, k.IcaConnectionId(ctx), portID)
	if !found {
		return "", "", types.ErrPortIdNotFound
	}

	return portID, addr, nil
}

func (k *Keeper) EndBlockerIBCTransfer(goCtx context.Context, from string, to string, token sdk.Coin) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	channelID := types.IBCTransferChannelID

	// transfer tokens to Osmosis network
	timeoutTimestamp := ctx.BlockTime().Add(time.Minute).UnixNano()
	_, err := k.TransferKeeper.Transfer(goCtx, ibctransfertypes.NewMsgTransfer(
		ibctransfertypes.PortID,
		channelID,
		token,
		from,
		to,
		clienttypes.Height{},
		uint64(timeoutTimestamp),
		""))
	return err
}

func (k *Keeper) EndBlockerSendPacket(ctx sdk.Context, portID string, data []byte) error {
	channelID, found := k.icaControllerKeeper.GetOpenActiveChannel(ctx, k.IcaConnectionId(ctx), portID)
	if !found {
		k.Logger(ctx).Debug("endblocker - send packet", "channel not found for port id", portID)
		return types.ErrChannelNotFound
	}

	chanCap, found := k.IBCScopperKeeper.GetCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if !found {
		return channeltypes.ErrChannelCapabilityNotFound.Wrap("module does not own channel capability")
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	timeoutTimestamp := ctx.BlockTime().Add(time.Minute).UnixNano()
	_, err := k.icaControllerKeeper.SendTx(ctx, chanCap, k.IcaConnectionId(ctx), portID, packetData, uint64(timeoutTimestamp))
	if err != nil {
		return err
	}
	return nil
}

func (k *Keeper) EndBlockerTransferDepositTokenToICA(goCtx context.Context, icaAddr string) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.DepositTokenToICAPacketSend(ctx) {
		return nil
	}

	if k.DepositTokenToICAPacketSent(ctx) {
		return nil
	}

	if k.DepositLastTime(ctx)+uint64(time.Minute) > uint64(ctx.BlockTime().UnixNano()) {
		k.Logger(ctx).Debug("endblocker - deposit to ica", "depositlasttime", k.DepositLastTime(ctx)+uint64(time.Minute), "currenttime", uint64(ctx.BlockTime().UnixNano()))
		return nil
	}

	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	sendToken := k.bankKeeper.GetBalance(ctx, moduleAddr, k.DepositDenom(ctx))
	k.Logger(ctx).Debug("endblocker - deposit to ica", "sendtoken", sendToken)

	if sendToken.Amount.LTE(math.ZeroInt()) || sendToken.Denom == "deposit_denom" {
		return nil
	}

	err := k.EndBlockerIBCTransfer(goCtx, moduleAddr.String(), icaAddr, sendToken)
	if err != nil {
		k.Logger(ctx).Debug("endblocker - deposit to ica", "transfer error", err)
		return err
	}

	k.SetCurrentDepositAmountParam(ctx, sendToken)
	k.SetDepositTokenToICAPacketSentParam(ctx, true)

	return nil
}

func (k *Keeper) EndBlockerJoinSwapExactAmountIn(goCtx context.Context, portID string, addr string) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.JoinSwapExactAmountInPacketSend(ctx) {
		return nil
	}

	if k.JoinSwapExactAmountInPacketSent(ctx) {
		return nil
	}

	depositToken := k.CurrentDepositAmount(ctx)

	if depositToken.Amount.LTE(math.ZeroInt()) {
		return nil
	}

	if depositToken.Denom == types.DefaultDepositDenom {
		return nil
	}

	msg := osmosispool.MsgJoinSwapExternAmountIn{
		Sender:            addr,
		PoolId:            1,
		TokenIn:           depositToken,
		ShareOutMinAmount: math.ZeroInt(),
	}

	data, err := icatypes.SerializeCosmosTx(k.cdc, []proto.Message{&msg})
	if err != nil {
		k.Logger(ctx).Debug("endblocker - joinswapexternamountin", "serialize cosmos tx failed", err)
		return err
	}

	err = k.EndBlockerSendPacket(ctx, portID, data)
	if err != nil {
		k.Logger(ctx).Debug("endblocker - joinswapexternamountin", "send tx failed", err)
		return err
	}

	k.Logger(ctx).Debug("endblocker - joinswapexternamountin", "send tx success")

	k.SetJoinSwapExactAmountInPacketSentParam(ctx, true)

	return nil
}

func (k *Keeper) EndBlockerLockUpLiquidity(goCtx context.Context, portID string, addr string) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.LockTokensPacketSend(ctx) {
		return nil
	}

	if k.LockTokensPacketSent(ctx) {
		return nil
	}

	liquidity := k.bankKeeper.GetBalance(ctx, sdk.AccAddress(addr), k.LiquidityDenom(ctx))

	if liquidity.Amount.LTE(math.ZeroInt()) {
		return nil
	}

	if liquidity.Denom == types.DefaultLiquidityDenom {
		return nil
	}

	msg := osmosislockup.MsgLockTokens{
		Owner:    addr,
		Duration: types.LockDuration,
		Coins:    sdk.NewCoins(liquidity),
	}

	data, err := icatypes.SerializeCosmosTx(k.cdc, []proto.Message{&msg})
	if err != nil {
		k.Logger(ctx).Debug("endblocker - lockupliquidity", "serialize cosmos tx failed", err)
		return err
	}

	err = k.EndBlockerSendPacket(ctx, portID, data)
	if err != nil {
		k.Logger(ctx).Debug("endblocker - lockupliquidity", "send tx failed", err)
		return err
	}

	k.Logger(ctx).Debug("endblocker - lockupliquidity", "send tx success")

	k.SetLockTokensPacketSentParam(ctx, true)

	return nil
}

func (k *Keeper) EndBlockerUnLockLiquidity(goCtx context.Context, portID string, addr string) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.UnLockLiquidityPacketSend(ctx) {
		return nil
	}

	if k.UnLockLiquidityPacketSent(ctx) {
		return nil
	}

	if uint64(ctx.BlockTime().UnixNano())-k.LockTokenTimestamp(ctx) < uint64(types.LockDuration.Nanoseconds()) {
		return nil
	}

	msg := osmosislockup.MsgBeginUnlockingAll{
		Owner: addr,
	}

	data, err := icatypes.SerializeCosmosTx(k.cdc, []proto.Message{&msg})
	if err != nil {
		k.Logger(ctx).Debug("endblocker - unlock", "serialize cosmos tx failed", err)
		return err
	}

	err = k.EndBlockerSendPacket(ctx, portID, data)

	if err != nil {
		k.Logger(ctx).Debug("endblocker - unlock", "send tx failed", err)
		return err
	}

	k.Logger(ctx).Debug("endblocker - unlock", "unlock liquidity send tx success")

	k.SetUnLockLiquidityPacketSentParam(ctx, true)

	return nil
}

func (k *Keeper) EndBlockerTransferRewards(goCtx context.Context, icaAddr string) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.ClaimRewardPacketSend(ctx) {
		return
	}

	if k.ClaimRewardPacketSent(ctx) {
		return
	}

	balance := k.bankKeeper.GetBalance(ctx, sdk.AccAddress(icaAddr), k.DepositDenom(ctx))

	k.Logger(ctx).Debug("endblocker - transfer reward", "balance", balance)

	feeCollectorAddr := authtypes.NewModuleAddress(authtypes.FeeCollectorName)
	err := k.EndBlockerIBCTransfer(goCtx, icaAddr, feeCollectorAddr.String(), balance)
	if err != nil {
		k.Logger(ctx).Debug("endblocker - transfer reward", "transfer failed", err)
	}

	k.SetClaimRewardPacketSentParam(ctx, true)

	k.Logger(ctx).Debug("endblocker - transfer reward", "success")
}
