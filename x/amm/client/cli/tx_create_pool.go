package cli

import (
	"strconv"
	"strings"

	"encoding/json"
	"frogchain/x/amm/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdCreatePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [pool-param] [pool-assets]",
		Short: "Broadcast message create-pool",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPoolParam := new(types.PoolParam)

			// unmarshal PoolParam
			err = json.Unmarshal([]byte(args[0]), argPoolParam)
			if err != nil {
				return err
			}

			// unmarshal PoolAssets
			argCastPoolAssets := strings.Split(args[1], listSeparator)
			argPoolAssets := make([]*types.PoolToken, len(argCastPoolAssets))
			for i, arg := range argCastPoolAssets {
				argPoolAsset := new(types.PoolToken)
				err = json.Unmarshal([]byte(arg), argPoolAssets[i])
				if err != nil {
					return err
				}
				argPoolAssets[i] = argPoolAsset
			}

			// get asset Amounts
			argCastAssetAmounts := strings.Split(args[2], listSeparator)
			argAssetAmounts := make([]uint64, len(argCastAssetAmounts))
			for i, arg := range argCastAssetAmounts {
				value, err := cast.ToUint64E(arg)
				if err != nil {
					return err
				}
				argAssetAmounts[i] = value
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreatePool(
				clientCtx.GetFromAddress().String(),
				argPoolParam,
				argPoolAssets,
				argAssetAmounts,
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
