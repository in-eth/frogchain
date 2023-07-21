package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyAdminAccount = []byte("AdminAccount")
	// TODO: Determine the default value
	DefaultAdminAccount string = "admin_account"
)

var (
	KeyDepositDenom = []byte("DepositDenom")
	// TODO: Determine the default value
	DefaultDepositDenom string = "deposit_denom"
)

var (
	KeyCurrentDepositAmount = []byte("CurrentDepositAmount")
	// TODO: Determine the default value
	DefaultCurrentDepositAmount sdk.Coin = sdk.NewCoin(DefaultDepositDenom, sdk.ZeroInt())
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
) Params {
	return Params{
		AdminAccount:         adminAccount,
		DepositDenom:         depositDenom,
		CurrentDepositAmount: currentDepositAmount,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultAdminAccount,
		DefaultDepositDenom,
		DefaultCurrentDepositAmount,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAdminAccount, &p.AdminAccount, validateAdminAccount),
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
	adminAccount, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// TODO implement validation
	_ = adminAccount

	return nil
}
