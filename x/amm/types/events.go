package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// amm module event types
const (
	// supply and balance tracking events name and attributes
	EventTypeCreatePool               = "pool_created"
	EventTypeAddLiquidity             = "liquidity_added"
	EventTypeRemoveLiquidity          = "liquidity_removed"
	EventTypeSwapExactTokensForTokens = "tokens_swapped"
	EventTypeSwapTokensForExactTokens = "exact_tokens_swapped"

	AttributeKeySpender  = "spender"
	AttributeKeyReceiver = "receiver"
	AttributePoolId      = "poolid"
	AttirbuteShareToken  = "sharetoken"
	AttirbuteAssets      = "assets"
	AttributeAmountIn    = "swapamountin"
	AttributeAmountOut   = "swapamountout"
)

// NewCreatePoolEvent constructs a new pool created sdk.Event
func NewCreatePoolEvent(creator sdk.AccAddress, poolId uint64, shareToken sdk.Coin) sdk.Event {
	return sdk.NewEvent(
		EventTypeCreatePool,
		sdk.NewAttribute(AttributeKeySpender, creator.String()),
		sdk.NewAttribute(AttributePoolId, fmt.Sprint(poolId)),
		sdk.NewAttribute(AttirbuteShareToken, shareToken.String()),
	)
}

// NewAddLiquidityEvent constructs a new liquidity provided sdk.Event
func NewAddLiquidityEvent(creator sdk.AccAddress, poolId uint64, shareToken sdk.Coin) sdk.Event {
	return sdk.NewEvent(
		EventTypeAddLiquidity,
		sdk.NewAttribute(AttributeKeySpender, creator.String()),
		sdk.NewAttribute(AttributePoolId, fmt.Sprint(poolId)),
		sdk.NewAttribute(AttirbuteShareToken, shareToken.String()),
	)
}

// NewRemoveLiquidityEvent constructs a new liquidity burnt sdk.Event
func NewRemoveLiquidityEvent(creator sdk.AccAddress, poolId uint64, amounts sdk.Coins) sdk.Event {
	return sdk.NewEvent(
		EventTypeRemoveLiquidity,
		sdk.NewAttribute(AttributeKeySpender, creator.String()),
		sdk.NewAttribute(AttributePoolId, fmt.Sprint(poolId)),
		sdk.NewAttribute(AttirbuteAssets, amounts.String()),
	)
}

// NewSwapExactTokensForTokensEvent constructs a new swap sdk.Event
func NewSwapExactTokensForTokensEvent(creator sdk.AccAddress, poolId uint64, amountOut uint64) sdk.Event {
	return sdk.NewEvent(
		EventTypeSwapExactTokensForTokens,
		sdk.NewAttribute(AttributeKeySpender, creator.String()),
		sdk.NewAttribute(AttributePoolId, fmt.Sprint(poolId)),
		sdk.NewAttribute(AttributeAmountOut, fmt.Sprint(amountOut)),
	)
}

// NewSwapTokensForExactTokensEvent constructs a new liquidity burnt sdk.Event
func NewSwapTokensForExactTokensEvent(creator sdk.AccAddress, poolId uint64, amountIn uint64) sdk.Event {
	return sdk.NewEvent(
		EventTypeSwapTokensForExactTokens,
		sdk.NewAttribute(AttributeKeySpender, creator.String()),
		sdk.NewAttribute(AttributePoolId, fmt.Sprint(poolId)),
		sdk.NewAttribute(AttributeAmountOut, fmt.Sprint(amountIn)),
	)
}
