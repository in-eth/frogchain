package cli

import (
	"strconv"

	"frogchain/x/amm/types"

	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdRemoveLiquidity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-liquidity [pool-id] [desired-amount] [min-amounts]",
		Short: "Broadcast message remove-liquidity",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			// get pool id
			argPoolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			// get liquidity amount to burn
			argLiquidity, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			// get min token amounts to get
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

			msg := types.NewMsgRemoveLiquidity(
				clientCtx.GetFromAddress().String(),
				argPoolId,
				argLiquidity,
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
