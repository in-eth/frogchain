package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"frogchain/testutil/network"
	"frogchain/testutil/nullify"
	"frogchain/x/investibc/client/cli"
	"frogchain/x/investibc/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithDepositBalanceObjects(t *testing.T, n int) (*network.Network, []types.DepositBalance) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{PortId: types.PortID}
	for i := 0; i < n; i++ {
		depositBalance := types.DepositBalance{
			Index: strconv.Itoa(i),
		}
		nullify.Fill(&depositBalance)
		state.DepositBalanceList = append(state.DepositBalanceList, depositBalance)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.DepositBalanceList
}

func TestShowDepositBalance(t *testing.T) {
	net, objs := networkWithDepositBalanceObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc    string
		idIndex string

		args []string
		err  error
		obj  types.DepositBalance
	}{
		{
			desc:    "found",
			idIndex: objs[0].Index,

			args: common,
			obj:  objs[0],
		},
		{
			desc:    "not found",
			idIndex: strconv.Itoa(100000),

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idIndex,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowDepositBalance(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetDepositBalanceResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.DepositBalance)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.DepositBalance),
				)
			}
		})
	}
}

func TestListDepositBalance(t *testing.T) {
	net, objs := networkWithDepositBalanceObjects(t, 5)

	ctx := net.Validators[0].ClientCtx
	request := func(next []byte, offset, limit uint64, total bool) []string {
		args := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		if next == nil {
			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
		} else {
			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
		}
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
		if total {
			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
		}
		return args
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListDepositBalance(), args)
			require.NoError(t, err)
			var resp types.QueryAllDepositBalanceResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.DepositBalance), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.DepositBalance),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListDepositBalance(), args)
			require.NoError(t, err)
			var resp types.QueryAllDepositBalanceResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.DepositBalance), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.DepositBalance),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListDepositBalance(), args)
		require.NoError(t, err)
		var resp types.QueryAllDepositBalanceResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.DepositBalance),
		)
	})
}
