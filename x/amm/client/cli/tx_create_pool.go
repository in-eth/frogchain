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

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ = strconv.Itoa(0)

func CmdCreatePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [pool-param] [pool-assets] [asset-weights]",
		Short: "Broadcast message create-pool",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			// unmarshal PoolParam
			argPoolParam := new(types.PoolParam)

			err = json.Unmarshal([]byte(args[0]), argPoolParam)
			if err != nil {
				return err
			}

			// unmarshal PoolAssets
			argPoolAssets, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			// get asset Amounts
			argCastAssetWeights := strings.Split(args[2], listSeparator)
			argAssetWeights := make([]uint64, len(argCastAssetWeights))
			for i, arg := range argCastAssetWeights {
				value, err := cast.ToUint64E(arg)
				if err != nil {
					return err
				}
				argAssetWeights[i] = value
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreatePool(
				clientCtx.GetFromAddress().String(),
				*argPoolParam,
				argPoolAssets,
				argAssetWeights,
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
