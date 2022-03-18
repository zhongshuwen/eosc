package cmd

import (
	"context"

	"github.com/zhongshuwen/zswchain-go/system"
	"github.com/spf13/cobra"
)

var systemSetcodeCmd = &cobra.Command{
	Use:   "setcode [account name] [wasm file]",
	Short: "Set code only on an account.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		accountName := toAccount(args[0], "account name")
		wasmFile := args[1]

		action, err := system.NewSetCode(accountName, wasmFile)
		errorCheck("loading wasm file", err)

		pushEOSCActions(context.Background(), api, action)
	},
}

func init() {
	systemCmd.AddCommand(systemSetcodeCmd)
}
