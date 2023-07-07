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

func CmdSwapExactTokensForTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap-exact-tokens-for-tokens [pool-id] [amount-in] [amount-out-min] [path] [to] [deadline]",
		Short: "Broadcast message swap-exact-tokens-for-tokens",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			// get pool id
			argPoolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			// get input token amount
			argAmountIn := sdk.MustNewDecFromStr(args[1])

			// get minimum token output amount
			argAmountOutMin := sdk.MustNewDecFromStr(args[2])

			// unmarshal token path
			argPath := strings.Split(args[3], listSeparator)

			// get receiver address
			argTo := args[4]

			// get swap deadline
			argDeadline, err := cast.ToTimeE(args[5])
			if err != nil {
				return err
			}

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
