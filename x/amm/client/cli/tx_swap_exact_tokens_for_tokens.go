package cli

import (
	"strconv"

	"frogchain/x/amm/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdSwapExactTokensForTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap-exact-tokens-for-tokens [pool-id] [amount-in] [amount-out-min] [path] [to] [deadline]",
		Short: "Broadcast message swap-exact-tokens-for-tokens",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPoolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argAmountIn, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}
			argAmountOutMin, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}
			argPath, err := sdk.ParseCoinsNormalized(args[3])
			if err != nil {
				return err
			}
			argTo := args[4]
			argDeadline := args[5]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSwapExactTokensForTokens(
				clientCtx.GetFromAddress().String(),
				argPoolId,
				argAmountIn,
				argAmountOutMin,
				argPath,
				argTo,
				argDeadline,
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
