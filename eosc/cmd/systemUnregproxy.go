package cmd

import (
	"context"

	"github.com/zhongshuwen/zswchain-go/system"
	"github.com/spf13/cobra"
)

var systemUnregisterProxyCmd = &cobra.Command{
	Use:   "unregproxy [account_name]",
	Short: "Unregister account as voting proxy.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()
		accountName := toAccount(args[0], "account name")
		pushEOSCActions(context.Background(), api, system.NewRegProxy(accountName, false))
	},
}

func init() {
	systemCmd.AddCommand(systemUnregisterProxyCmd)
}
