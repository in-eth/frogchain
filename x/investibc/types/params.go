package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyAdminAccount            = []byte("AdminAccount")
	DefaultAdminAccount string = "admin_account"
)

var (
	KeyDepositDenom            = []byte("DepositDenom")
	DefaultDepositDenom string = "deposit_denom"
)

var (
	KeyCurrentDepositAmount              = []byte("CurrentDepositAmount")
	DefaultCurrentDepositAmount sdk.Coin = sdk.NewCoin(DefaultDepositDenom, sdk.ZeroInt())
)

var (
	KeyLiquidityDenom            = []byte("LiquidityDenom")
	DefaultLiquidityDenom string = "liquidity_denom"
)

var (
	KeyCurrentLiquidityAmount              = []byte("CurrentLiquidityAmount")
	DefaultCurrentLiquidityAmount sdk.Coin = sdk.NewCoin(DefaultLiquidityDenom, sdk.ZeroInt())
)

var (
	KeyDepositLastTime            = []byte("DepositLastTime")
	DefaultDepositLastTime uint64 = 10
)

var (
	KeyIcaConnectionId            = []byte("IcaConnectionId")
	DefaultIcaConnectionId string = ""
)

var (
	KeyJoinSwapExactAmountInPacketSent          = []byte("JoinSwapExactAmountInPacketSent")
	DefaultJoinSwapExactAmountInPacketSent bool = false
)

var (
	KeyLockTokensPacketSent          = []byte("LockTokensPacketSent")
	DefaultLockTokensPacketSent bool = false
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	adminAccount string,
	depositDenom string,
	currentDepositAmount sdk.Coin,
	liquidityDenom string,
	currentLiquidityAmount sdk.Coin,
	depositLastTime uint64,
	icaConnectionId string,
	liquidityBootstrapping bool,
	liquidityBootstrapped bool,
) Params {
	return Params{
		AdminAccount:                    adminAccount,
		DepositDenom:                    depositDenom,
		CurrentDepositAmount:            currentDepositAmount,
		LiquidityDenom:                  liquidityDenom,
		CurrentLiquidityAmount:          currentLiquidityAmount,
		DepositLastTime:                 depositLastTime,
		IcaConnectionId:                 icaConnectionId,
		JoinSwapExactAmountInPacketSent: liquidityBootstrapping,
		LockTokensPacketSent:            liquidityBootstrapped,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultAdminAccount,
		DefaultDepositDenom,
		DefaultCurrentDepositAmount,
		DefaultLiquidityDenom,
		DefaultCurrentLiquidityAmount,
		DefaultDepositLastTime,
		DefaultIcaConnectionId,
		DefaultJoinSwapExactAmountInPacketSent,
		DefaultLockTokensPacketSent,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAdminAccount, &p.AdminAccount, validateAdminAccount),
		paramtypes.NewParamSetPair(KeyDepositDenom, &p.DepositDenom, validateDenom),
		paramtypes.NewParamSetPair(KeyCurrentDepositAmount, &p.CurrentDepositAmount, validateCoin),
		paramtypes.NewParamSetPair(KeyLiquidityDenom, &p.LiquidityDenom, validateDenom),
		paramtypes.NewParamSetPair(KeyCurrentLiquidityAmount, &p.CurrentLiquidityAmount, validateCoin),
		paramtypes.NewParamSetPair(KeyDepositLastTime, &p.DepositLastTime, validateDepositEndTime),
		paramtypes.NewParamSetPair(KeyIcaConnectionId, &p.IcaConnectionId, validateIcaConnectionId),
		paramtypes.NewParamSetPair(KeyJoinSwapExactAmountInPacketSent, &p.JoinSwapExactAmountInPacketSent, validateBool),
		paramtypes.NewParamSetPair(KeyLockTokensPacketSent, &p.LockTokensPacketSent, validateBool),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateAdminAccount(p.AdminAccount); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// validateAdminAccount validates the AdminAccount param
func validateAdminAccount(v interface{}) error {
	_, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	return nil
}

func validateDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateVestingDuration(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("vesting duration should be positive")
	}
	return nil
}

func validateDepositEndTime(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// if v == 0 {
	// 	return fmt.Errorf("deposit end time should be positive")
	// }
	return nil
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateIcaConnectionId(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateCoin(i interface{}) error {
	_, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
