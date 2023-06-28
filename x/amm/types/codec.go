package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreatePool{}, "amm/CreatePool", nil)
	cdc.RegisterConcrete(&MsgAddLiquidity{}, "amm/AddLiquidity", nil)
	cdc.RegisterConcrete(&MsgRemoveLiquidity{}, "amm/RemoveLiquidity", nil)
	cdc.RegisterConcrete(&MsgSwapExactTokensForTokens{}, "amm/SwapExactTokensForTokens", nil)
	cdc.RegisterConcrete(&MsgSwapTokensForExactTokens{}, "amm/SwapTokensForExactTokens", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreatePool{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddLiquidity{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRemoveLiquidity{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSwapExactTokensForTokens{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSwapTokensForExactTokens{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
