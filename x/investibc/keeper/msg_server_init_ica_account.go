package keeper

import (
	"context"

	"frogchain/x/investibc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
)

func (k msgServer) InitIcaAccount(goCtx context.Context, msg *types.MsgInitIcaAccount) (*types.MsgInitIcaAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Creator != k.AdminAccount(ctx) {
		return nil, types.ErrAdminPermission
	}

	k.SetIcaConnectionIdParam(ctx, msg.ConnectionId)

	// Get ConnectionEnd (for counterparty connection)
	connectionEnd, found := k.Keeper.IBCKeeper.ConnectionKeeper.GetConnection(ctx, msg.ConnectionId)
	if !found {
		return nil, types.ErrConnectionNotFound
	}
	counterpartyConnection := connectionEnd.Counterparty

	appVersion := string(icatypes.ModuleCdc.MustMarshalJSON(&icatypes.Metadata{
		Version:                icatypes.Version,
		ControllerConnectionId: msg.ConnectionId,
		HostConnectionId:       counterpartyConnection.ConnectionId,
		Encoding:               icatypes.EncodingProtobuf,
		TxType:                 icatypes.TxTypeSDKMultiMsg,
	}))

	msgServer := icacontrollerkeeper.NewMsgServerImpl(&k.icaControllerKeeper)
	msgRegisterInterchainAccount := icacontrollertypes.NewMsgRegisterInterchainAccount(msg.ConnectionId, types.ModuleName, appVersion)

	_, err := msgServer.RegisterInterchainAccount(sdk.WrapSDKContext(ctx), msgRegisterInterchainAccount)
	if err != nil {
		return nil, err
	}

	portID, err := icatypes.NewControllerPortID(types.ModuleName)
	if err != nil {
		return nil, err
	}

	k.icaControllerKeeper.SetMiddlewareEnabled(ctx, portID, msg.ConnectionId)

	return &types.MsgInitIcaAccountResponse{}, nil
}
