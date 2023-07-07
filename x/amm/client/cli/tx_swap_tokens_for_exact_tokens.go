package cli

import (
	"strconv"
	"strings"

	"frogchain/x/amm/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdSwapTokensForExactTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap-tokens-for-exact-tokens [pool-id] [amount-out] [path] [to] [deadline]",
		Short: "Broadcast message swap-tokens-for-exact-tokens",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			// get pool id
			argPoolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			// get output token amount out
			argAmountOut := sdk.MustNewDecFromStr(args[1])
			if err != nil {
				return err
			}

			// unmarshal token path
			argPath := strings.Split(args[2], listSeparator)

			// get receiver address
			argTo := args[3]

			// get swap deadline
			argDeadline, err := cast.ToTimeE(args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSwapTokensForExactTokens(
				clientCtx.GetFromAddress().String(),
				argPoolId,
				argAmountOut,
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
