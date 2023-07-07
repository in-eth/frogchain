package cli

import (
	"strconv"

	"frogchain/x/amm/types"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetSwapTokensForExactTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-swap-tokens-for-exact-tokens [pool-id] [amount-out] [path]",
		Short: "Query get-swap-tokens-for-exact-tokens",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqPoolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			reqAmountOut, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}
			reqPath := strings.Split(args[2], listSeparator)

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetSwapTokensForExactTokensRequest{

				PoolId:    reqPoolId,
				AmountOut: reqAmountOut,
				Path:      reqPath,
			}

			res, err := queryClient.GetSwapTokensForExactTokens(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
