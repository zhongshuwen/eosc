// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"context"

	"github.com/zhongshuwen/zswchain-go/system"
	"github.com/spf13/cobra"
)

var systemClaimRewardsCmd = &cobra.Command{
	Use:   "claimrewards [owner]",
	Short: "Claim block production rewards. Once per day, don't forget it!",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		owner := toAccount(args[0], "owner")

		pushEOSCActions(context.Background(), api,
			system.NewClaimRewards(owner),
		)
	},
}

func init() {
	systemCmd.AddCommand(systemClaimRewardsCmd)
}
