package cli

import (
	"strconv"

	"frogchain/x/amm/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"strings"
)

var _ = strconv.Itoa(0)

func CmdGetSwapExactTokensForTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-swap-exact-tokens-for-tokens [pool-id] [amount-in] [path]",
		Short: "Query get-swap-exact-tokens-for-tokens",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqPoolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			reqAmountIn, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}
			reqPath := strings.Split(args[2], listSeparator)

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetSwapExactTokensForTokensRequest{

				PoolId:   reqPoolId,
				AmountIn: reqAmountIn,
				Path:     reqPath,
			}

			res, err := queryClient.GetSwapExactTokensForTokens(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
