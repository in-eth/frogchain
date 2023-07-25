package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSetAdminAccount{}, "investibc/SetAdminAccount", nil)
	cdc.RegisterConcrete(&MsgDeposit{}, "investibc/Deposit", nil)
	cdc.RegisterConcrete(&MsgInitIcaAccount{}, "investibc/RegisterIcaAccount", nil)
	cdc.RegisterConcrete(&MsgSetDepositDenom{}, "investibc/SetDepositDenom", nil)
	cdc.RegisterConcrete(&MsgWithdraw{}, "investibc/Withdraw", nil)
	cdc.RegisterConcrete(&MsgSetLiquidityDenom{}, "investibc/SetLiquidityDenom", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetAdminAccount{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeposit{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgInitIcaAccount{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetDepositDenom{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgWithdraw{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetLiquidityDenom{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
