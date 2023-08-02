package investibc

import (
	"fmt"

	"cosmossdk.io/math"
	proto "github.com/gogo/protobuf/proto"

	"frogchain/x/investibc/keeper"
	"frogchain/x/investibc/types"

	osmosislockup "frogchain/osmosis/lockup"
	osmosispool "frogchain/osmosis/pool"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
)

type IBCModule struct {
	keeper keeper.Keeper
}

func NewIBCModule(k keeper.Keeper) IBCModule {
	return IBCModule{
		keeper: k,
	}
}

// OnChanOpenInit implements the IBCModule interface
func (im IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {

	if chanCap == nil {
		path := host.ChannelCapabilityPath(portID, channelID)
		chanCap, _ := im.keeper.IBCScopperKeeper.GetCapability(ctx, path)

		chanCap = chanCap
	}
	return version, im.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID))

	// // Require portID is the portID module is bound to
	// boundPort := im.keeper.GetPort(ctx)
	// if boundPort != portID {
	// 	return "", sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, boundPort)
	// }

	// if version != types.Version {
	// 	return "", sdkerrors.Wrapf(types.ErrInvalidVersion, "got %s, expected %s", version, types.Version)
	// }

	// // Claim channel capability passed back by IBC module
	// if err := im.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
	// 	return "", err
	// }

	// return version, nil
}

// OnChanOpenTry implements the IBCModule interface
func (im IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {

	// Require portID is the portID module is bound to
	// boundPort := im.keeper.GetPort(ctx)
	// if boundPort != portID {
	// 	return "", sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, boundPort)
	// }

	// if counterpartyVersion != types.Version {
	// 	return "", sdkerrors.Wrapf(types.ErrInvalidVersion, "invalid counterparty version: got: %s, expected %s", counterpartyVersion, types.Version)
	// }

	// // Module may have already claimed capability in OnChanOpenInit in the case of crossing hellos
	// // (ie chainA and chainB both call ChanOpenInit before one of them calls ChanOpenTry)
	// // If module can already authenticate the capability then module already owns it so we don't need to claim
	// // Otherwise, module does not have channel capability and we must claim it from IBC
	// if !im.keeper.AuthenticateCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)) {
	// 	// Only claim channel capability passed back by IBC module if we do not already own it
	// 	if err := im.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
	// 		return "", err
	// 	}
	// }

	// return types.Version, nil

	return "", nil
}

// OnChanOpenAck implements the IBCModule interface
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	_,
	counterpartyVersion string,
) error {
	// if counterpartyVersion != types.Version {
	// 	return sdkerrors.Wrapf(types.ErrInvalidVersion, "invalid counterparty version: %s, expected %s", counterpartyVersion, types.Version)
	// }
	return nil
}

// OnChanOpenConfirm implements the IBCModule interface
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnChanCloseInit implements the IBCModule interface
func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// Disallow user-initiated channel closing for channels
	// return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
	return nil
}

// OnChanCloseConfirm implements the IBCModule interface
func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnRecvPacket implements the IBCModule interface
func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	// var ack channeltypes.Acknowledgement

	// // this line is used by starport scaffolding # oracle/packet/module/recv

	// var modulePacketData types.InvestibcPacketData
	// if err := modulePacketData.Unmarshal(modulePacket.GetData()); err != nil {
	// 	return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet data: %s", err.Error()))
	// }

	// // Dispatch packet
	// switch packet := modulePacketData.Packet.(type) {
	// // this line is used by starport scaffolding # ibc/packet/module/recv
	// default:
	// 	err := fmt.Errorf("unrecognized %s packet type: %T", types.ModuleName, packet)
	// 	return channeltypes.NewErrorAcknowledgement(err)
	// }

	// // NOTE: acknowledgement will be written synchronously during IBC handler execution.
	// return ack
	return channeltypes.NewErrorAcknowledgement(fmt.Errorf("cannot receive packet via interchain accounts authentication module"))
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	var ack channeltypes.Acknowledgement
	if err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet acknowledgement: %v", err)
	}

	txMsgData := &sdk.TxMsgData{}
	if err := proto.Unmarshal(ack.GetResult(), txMsgData); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 tx message data: %v", err)
	}

	switch len(txMsgData.Data) {
	case 0:
		// TODO: handle for sdk 0.46.x
		return nil
	default:
		for _, msgData := range txMsgData.Data {
			response, err := handleMsgData(ctx, im.keeper, msgData)
			if err != nil {
				return err
			}

			im.keeper.Logger(ctx).Debug("message response in ICS-27 packet", "response", response)
		}
		return nil
	}
}

// OnTimeoutPacket implements the IBCModule interface
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	if im.keeper.DepositTokenToICAPacketSent(ctx) {
		im.keeper.SetDepositTokenToICAPacketSentParam(ctx, false)
		im.keeper.Logger(ctx).Debug("message timeout", "deposit token to ica")
	} else if im.keeper.JoinSwapExactAmountInPacketSent(ctx) {
		im.keeper.SetJoinSwapExactAmountInPacketSentParam(ctx, false)
		im.keeper.Logger(ctx).Debug("message timeout", "join swap exact amount in")
	} else if im.keeper.LockTokensPacketSent(ctx) {
		im.keeper.SetLockTokensPacketSentParam(ctx, false)
		im.keeper.Logger(ctx).Debug("message timeout", "lock liquidity")
	} else if im.keeper.UnLockLiquidityPacketSent(ctx) {
		im.keeper.SetUnLockLiquidityPacketSentParam(ctx, false)
		im.keeper.Logger(ctx).Debug("message timeout", "unlock liquidity")
	} else if im.keeper.ClaimRewardPacketSent(ctx) {
		im.keeper.SetClaimRewardPacketSentParam(ctx, false)
		im.keeper.Logger(ctx).Debug("message timeout", "claim reward")
	}

	// // Dispatch packet
	// switch packet := modulePacketData.Packet.(type) {
	// // this line is used by starport scaffolding # ibc/packet/module/timeout

	// default:
	// 	errMsg := fmt.Sprintf("unrecognized %s packet type: %T", types.ModuleName, packet)
	// 	return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	// }

	return nil
}

func handleMsgData(ctx sdk.Context, k keeper.Keeper, msgData *sdk.MsgData) (string, error) {
	switch msgData.MsgType {
	case sdk.MsgTypeURL(&ibctransfertypes.MsgTransfer{}):
		msgResponse := &ibctransfertypes.MsgTransferResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal send response message: %s", err.Error())
		}

		if k.ClaimRewardPacketSent(ctx) {
			k.SetClaimRewardPacketSendParam(ctx, false)
			k.SetClaimRewardPacketSentParam(ctx, false)
			k.SetDepositTokenToICAPacketSendParam(ctx, true)
			return msgResponse.String(), nil
		}

		k.SetDepositLastTimeParam(ctx, uint64(ctx.BlockTime().UnixNano()))

		k.SetDepositTokenToICAPacketSendParam(ctx, false)
		k.SetDepositTokenToICAPacketSentParam(ctx, false)
		k.SetJoinSwapExactAmountInPacketSendParam(ctx, true)

		return msgResponse.String(), nil
	case sdk.MsgTypeURL(&osmosispool.MsgJoinSwapExternAmountIn{}):
		msgResponse := &osmosispool.MsgJoinSwapExternAmountInResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal send response message: %s", err.Error())
		}

		sendToken := k.CurrentDepositAmount(ctx)
		sendToken.Amount = math.ZeroInt()
		k.SetCurrentDepositAmountParam(ctx, sendToken)

		k.SetJoinSwapExactAmountInPacketSendParam(ctx, false)
		k.SetJoinSwapExactAmountInPacketSentParam(ctx, false)
		k.SetLockTokensPacketSendParam(ctx, true)

		return msgResponse.String(), nil

	case sdk.MsgTypeURL(&osmosislockup.MsgLockTokens{}):
		msgResponse := &osmosislockup.MsgLockTokensResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal send response message: %s", err.Error())
		}

		lockTokenTimestamp := uint64(ctx.BlockTime().UnixNano())
		k.SetLockTokenTimestampParam(ctx, lockTokenTimestamp)

		k.SetLockTokensPacketSendParam(ctx, false)
		k.SetLockTokensPacketSentParam(ctx, false)
		k.SetUnLockLiquidityPacketSendParam(ctx, true)

		return msgResponse.String(), nil

	case sdk.MsgTypeURL(&osmosislockup.MsgBeginUnlockingAll{}):
		msgResponse := &osmosislockup.MsgBeginUnlockingAllResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal send response message: %s", err.Error())
		}

		k.SetUnLockLiquidityPacketSendParam(ctx, false)
		k.SetUnLockLiquidityPacketSentParam(ctx, false)
		k.SetClaimRewardPacketSendParam(ctx, true)

		return msgResponse.String(), nil
	default:
		return "", nil
	}
}
