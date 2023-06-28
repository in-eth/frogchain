package cli

import (
	"strconv"

	"frogchain/x/amm/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"strings"
)

var _ = strconv.Itoa(0)

func CmdAddLiquidity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-liquidity [pool-id] [desired-amounts] [min-amounts]",
		Short: "Broadcast message add-liquidity",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPoolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argCastDesiredAmounts := strings.Split(args[1], listSeparator)
			argDesiredAmounts := make([]uint64, len(argCastDesiredAmounts))
			for i, arg := range argCastDesiredAmounts {
				value, err := cast.ToUint64E(arg)
				if err != nil {
					return err
				}
				argDesiredAmounts[i] = value
			}
			argCastMinAmounts := strings.Split(args[2], listSeparator)
			argMinAmounts := make([]uint64, len(argCastMinAmounts))
			for i, arg := range argCastMinAmounts {
				value, err := cast.ToUint64E(arg)
				if err != nil {
					return err
				}
				argMinAmounts[i] = value
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddLiquidity(
				clientCtx.GetFromAddress().String(),
				argPoolId,
				argDesiredAmounts,
				argMinAmounts,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
