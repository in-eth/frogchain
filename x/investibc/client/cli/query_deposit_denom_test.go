package cli_test

import (
	"fmt"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"

	"frogchain/testutil/network"
	"frogchain/testutil/nullify"
	"frogchain/x/investibc/client/cli"
	"frogchain/x/investibc/types"
)

func networkWithDepositDenomObjects(t *testing.T) (*network.Network, types.DepositDenom) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{PortId: types.PortID}
	depositDenom := &types.DepositDenom{}
	nullify.Fill(&depositDenom)
	state.DepositDenom = depositDenom
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), *state.DepositDenom
}

func TestShowDepositDenom(t *testing.T) {
	net, obj := networkWithDepositDenomObjects(t)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc string
		args []string
		err  error
		obj  types.DepositDenom
	}{
		{
			desc: "get",
			args: common,
			obj:  obj,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			var args []string
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowDepositDenom(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetDepositDenomResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.DepositDenom)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.DepositDenom),
				)
			}
		})
	}
}
