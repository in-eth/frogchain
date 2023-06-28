package cli

import (
	"strconv"

	"frogchain/x/amm/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdRemoveLiquidity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-liquidity [pool-id] [desired-amount] [min-amount]",
		Short: "Broadcast message remove-liquidity",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPoolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argDesiredAmount, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}
			argMinAmount, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveLiquidity(
				clientCtx.GetFromAddress().String(),
				argPoolId,
				argDesiredAmount,
				argMinAmount,
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
