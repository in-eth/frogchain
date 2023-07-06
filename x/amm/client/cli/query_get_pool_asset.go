package cli

import (
	"strconv"

	"frogchain/x/amm/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetPoolAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-pool-asset [pool-id] [asset-id]",
		Short: "Query get-pool-asset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqPoolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			reqAssetId, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetPoolAssetRequest{

				PoolId:  reqPoolId,
				AssetId: reqAssetId,
			}

			res, err := queryClient.GetPoolAsset(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
