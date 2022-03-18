package cmd

import (
	"context"
	"fmt"

	"github.com/zhongshuwen/zswchain-go"
	"github.com/zhongshuwen/zswchain-go/system"
	"github.com/spf13/cobra"
)

var voteCancelAllCmd = &cobra.Command{
	Use:   "cancel-all [voter name]",
	Short: "Cancel all votes currently cast for producers/delegated to a proxy.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		voterName := toAccount(args[0], "voter name")

		noProxy := zsw.AccountName("")
		var noVotes []zsw.AccountName
		pushEOSCActions(context.Background(), api,
			system.NewVoteProducer(
				voterName,
				noProxy,
				noVotes...,
			),
		)

		fmt.Printf("Consider using `eosc vote status %s` to confirm it has been applied.\n", voterName)
	},
}

func init() {
	voteCmd.AddCommand(voteCancelAllCmd)
}
