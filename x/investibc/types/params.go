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
	DefaultAdminAccount string = "frog1m9l358xunhhwds0568za49mzhvuxx9uxeu2c3r"
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
	KeyLockTokenTimestamp            = []byte("LockTokenTimestamp")
	DefaultLockTokenTimestamp uint64 = 10
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
	KeyDepositTokenToICAPacketSend          = []byte("DepositTokenToICAPacketSend")
	DefaultDepositTokenToICAPacketSend bool = true
)

var (
	KeyJoinSwapExactAmountInPacketSend          = []byte("JoinSwapExactAmountInPacketSend")
	DefaultJoinSwapExactAmountInPacketSend bool = false
)

var (
	KeyLockTokensPacketSend          = []byte("LockTokensPacketSend")
	DefaultLockTokensPacketSend bool = false
)

var (
	KeyUnLockLiquidityPacketSend          = []byte("UnLockLiquidityPacketSend")
	DefaultUnLockLiquidityPacketSend bool = false
)

var (
	KeyClaimRewardPacketSend          = []byte("ClaimRewardPacketSend")
	DefaultClaimRewardPacketSend bool = false
)

var (
	KeyDepositTokenToICAPacketSent          = []byte("DepositTokenToICAPacketSent")
	DefaultDepositTokenToICAPacketSent bool = false
)

var (
	KeyJoinSwapExactAmountInPacketSent          = []byte("JoinSwapExactAmountInPacketSent")
	DefaultJoinSwapExactAmountInPacketSent bool = false
)

var (
	KeyLockTokensPacketSent          = []byte("LockTokensPacketSent")
	DefaultLockTokensPacketSent bool = false
)

var (
	KeyUnLockLiquidityPacketSent          = []byte("UnLockLiquidityPacketSent")
	DefaultUnLockLiquidityPacketSent bool = false
)

var (
	KeyClaimRewardPacketSent          = []byte("ClaimRewardPacketSent")
	DefaultClaimRewardPacketSent bool = false
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
	lockTokenTimestamp uint64,
	depositLastTime uint64,
	icaConnectionId string,
	depositTokenToICAPacketSend bool,
	joinSwapExactAmountInPacketSend bool,
	lockTokensPacketSend bool,
	unLockLiquidityPacketSend bool,
	claimRewardPacketSend bool,
	depositTokenToICAPacketSent bool,
	joinSwapExactAmountInPacketSent bool,
	lockTokensPacketSent bool,
	unLockLiquidityPacketSent bool,
	claimRewardPacketSent bool,
) Params {
	return Params{
		AdminAccount:                    adminAccount,
		DepositDenom:                    depositDenom,
		CurrentDepositAmount:            currentDepositAmount,
		LiquidityDenom:                  liquidityDenom,
		LockTokenTimestamp:              lockTokenTimestamp,
		DepositLastTime:                 depositLastTime,
		IcaConnectionId:                 icaConnectionId,
		DepositTokenToICAPacketSend:     depositTokenToICAPacketSend,
		JoinSwapExactAmountInPacketSend: joinSwapExactAmountInPacketSend,
		LockTokensPacketSend:            lockTokensPacketSend,
		UnLockLiquidityPacketSend:       unLockLiquidityPacketSend,
		ClaimRewardPacketSend:           claimRewardPacketSend,
		DepositTokenToICAPacketSent:     depositTokenToICAPacketSent,
		JoinSwapExactAmountInPacketSent: joinSwapExactAmountInPacketSent,
		LockTokensPacketSent:            lockTokensPacketSent,
		UnLockLiquidityPacketSent:       unLockLiquidityPacketSent,
		ClaimRewardPacketSent:           claimRewardPacketSent,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultAdminAccount,
		DefaultDepositDenom,
		DefaultCurrentDepositAmount,
		DefaultLiquidityDenom,
		DefaultLockTokenTimestamp,
		DefaultDepositLastTime,
		DefaultIcaConnectionId,
		DefaultDepositTokenToICAPacketSend,
		DefaultJoinSwapExactAmountInPacketSend,
		DefaultLockTokensPacketSend,
		DefaultUnLockLiquidityPacketSend,
		DefaultClaimRewardPacketSend,
		DefaultDepositTokenToICAPacketSent,
		DefaultJoinSwapExactAmountInPacketSent,
		DefaultLockTokensPacketSent,
		DefaultUnLockLiquidityPacketSent,
		DefaultClaimRewardPacketSent,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAdminAccount, &p.AdminAccount, validateAdminAccount),
		paramtypes.NewParamSetPair(KeyDepositDenom, &p.DepositDenom, validateDenom),
		paramtypes.NewParamSetPair(KeyCurrentDepositAmount, &p.CurrentDepositAmount, validateCoin),
		paramtypes.NewParamSetPair(KeyLiquidityDenom, &p.LiquidityDenom, validateDenom),
		paramtypes.NewParamSetPair(KeyLockTokenTimestamp, &p.LockTokenTimestamp, validateTime),
		paramtypes.NewParamSetPair(KeyDepositLastTime, &p.DepositLastTime, validateTime),
		paramtypes.NewParamSetPair(KeyIcaConnectionId, &p.IcaConnectionId, validateIcaConnectionId),
		paramtypes.NewParamSetPair(KeyDepositTokenToICAPacketSend, &p.DepositTokenToICAPacketSend, validateBool),
		paramtypes.NewParamSetPair(KeyJoinSwapExactAmountInPacketSend, &p.JoinSwapExactAmountInPacketSend, validateBool),
		paramtypes.NewParamSetPair(KeyLockTokensPacketSend, &p.LockTokensPacketSend, validateBool),
		paramtypes.NewParamSetPair(KeyUnLockLiquidityPacketSend, &p.UnLockLiquidityPacketSend, validateBool),
		paramtypes.NewParamSetPair(KeyClaimRewardPacketSend, &p.ClaimRewardPacketSend, validateBool),
		paramtypes.NewParamSetPair(KeyDepositTokenToICAPacketSent, &p.DepositTokenToICAPacketSent, validateBool),
		paramtypes.NewParamSetPair(KeyJoinSwapExactAmountInPacketSent, &p.JoinSwapExactAmountInPacketSent, validateBool),
		paramtypes.NewParamSetPair(KeyLockTokensPacketSent, &p.LockTokensPacketSent, validateBool),
		paramtypes.NewParamSetPair(KeyUnLockLiquidityPacketSent, &p.UnLockLiquidityPacketSent, validateBool),
		paramtypes.NewParamSetPair(KeyClaimRewardPacketSent, &p.ClaimRewardPacketSent, validateBool),
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

func validateTime(i interface{}) error {
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
