// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"context"

	"github.com/zhongshuwen/zswchain-go"
	"github.com/zhongshuwen/zswchain-go/system"
	"github.com/spf13/cobra"
)

var systemDeleteauthCmd = &cobra.Command{
	Use:   `deleteauth [account] [permission_name]`,
	Short: "Removes a permission currently set on an account. See --help for more details.",
	Long: `Removes a permission currently set on an account.

This undoes the action of updateauth. Please refer to the updateauth help for more details.
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		account := toAccount(args[0], "account")
		permissionName := toName(args[1], "permission_name")

		api := getAPI()
		pushEOSCActions(context.Background(), api, system.NewDeleteAuth(account, zsw.PermissionName(permissionName)))
	},
}

func init() {
	systemCmd.AddCommand(systemDeleteauthCmd)
}
