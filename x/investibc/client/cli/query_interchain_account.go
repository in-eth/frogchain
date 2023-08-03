package cli

import (
	"strconv"

	"frogchain/x/investibc/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdInterchainAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "interchain-account [owner] [connection-id]",
		Short: "Query interchain-account",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqOwner := args[0]
			reqConnectionId := args[1]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryInterchainAccountRequest{
				Owner:        reqOwner,
				ConnectionId: reqConnectionId,
			}

			res, err := queryClient.InterchainAccount(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
