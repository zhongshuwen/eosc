package cmd

import (
	"context"
	"fmt"

	"github.com/zhongshuwen/zswchain-go/system"
	"github.com/spf13/cobra"
)

var voteProxyCmd = &cobra.Command{
	Use:   "proxy [voter name] [proxy name]",
	Short: "Proxy your vote strength to a proxy.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		voterName := toAccount(args[0], "voter name")
		proxyName := toAccount(args[1], "proxy name")

		fmt.Printf("Voter [%s] voting for proxy: %s\n", voterName, proxyName)

		pushEOSCActions(context.Background(), api,
			system.NewVoteProducer(
				voterName,
				proxyName,
			),
		)
	},
}

func init() {
	voteCmd.AddCommand(voteProxyCmd)
}
