package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	zsw "github.com/zhongshuwen/zswchain-go"
	"github.com/zhongshuwen/zswchain-go/msig"
	"github.com/zhongshuwen/zswchain-go/system"
	"github.com/zhongshuwen/zswchain-go/token"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var toolsSellAccountCmd = &cobra.Command{
	Use:   "sell-account [sold account] [buyer account] [beneficiary account] [amount]",
	Short: "Create a multisig transaction that both parties need to approve in order to do an atomic sale of your account.",
	Long: `Create a multisig transaction that both parties need to approve in order to do an atomic sale of your account.

Transfers both "owner" and "active" authority to a clone of the buyer's account's authority.

MAKE SURE TO INSPECT THE GENERATED MULTISIG TRANSACTION BEFORE APPROVING IT.
`,
	Args: cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		soldAccount := toAccount(args[0], "sold account")
		buyerAccount := toAccount(args[1], "buyer account")
		beneficiaryAccount := toAccount(args[2], "beneficiary account")
		saleAmount := toCoreAsset(args[3], "amount")
		proposalName := viper.GetString("tools-sell-account-cmd-proposal-name")
		memo := viper.GetString("tools-sell-account-cmd-memo")

		api := getAPI()

		soldAccountData, err := api.GetAccount(ctx, soldAccount)
		errorCheck("could not find sold account on chain: "+string(soldAccount), err)

		if len(soldAccountData.Permissions) > 2 {
			fmt.Println("WARNING: your account has more than 2 permissions.")
			fmt.Println("This operation hands off control of `owner` and `active` keys.")
			fmt.Println("Please clean-up your permissions before selling your account.")
			os.Exit(1)
		}

		buyerAccountData, err := api.GetAccount(ctx, buyerAccount)
		errorCheck("could not find buyer's account on chain", err)

		_, err = api.GetAccount(ctx, beneficiaryAccount)
		errorCheck("could not find beneficiary's account on chain", err)

		buyerPermText := viper.GetString("tools-sell-account-cmd-buyer-permission")
		if buyerPermText == "" {
			buyerPermText = string(buyerAccount)
		}
		buyerPerm, err := zsw.NewPermissionLevel(buyerPermText)
		errorCheck(`invalid "buyer-permission"`, err)

		myPermText := viper.GetString("tools-sell-account-cmd-seller-permission")
		if myPermText == "" {
			myPermText = string(soldAccount)
		}
		myPerm, err := zsw.NewPermissionLevel(myPermText)
		errorCheck(`invalid "seller-permission"`, err)

		targetOwnerAuth, err := sellAccountFindAuthority(buyerAccountData, "owner")
		errorCheck("error finding buyer's owner permission", err)
		targetActiveAuth, err := sellAccountFindAuthority(buyerAccountData, "active")
		errorCheck("error finding buyer's owner permission", err)

		infoResp, err := api.GetInfo(ctx)
		errorCheck("couldn't get_info from chain", err)

		tx := zsw.NewTransaction([]*zsw.Action{
			system.NewUpdateAuth(soldAccount, zsw.PermissionName("owner"), zsw.PermissionName(""), targetOwnerAuth, zsw.PermissionName("owner")),
			system.NewUpdateAuth(soldAccount, zsw.PermissionName("active"), zsw.PermissionName("owner"), targetActiveAuth, zsw.PermissionName("active")),
			token.NewTransfer(buyerAccount, beneficiaryAccount, saleAmount, memo),
		}, &zsw.TxOptions{HeadBlockID: infoResp.HeadBlockID})
		tx.SetExpiration(viper.GetDuration("tools-sell-account-cmd-sale-expiration"))

		fmt.Println("Submitting `eosio.msig` proposal:")
		fmt.Printf("  proposer: %s\n", soldAccount)
		fmt.Printf("  proposal_name: %s\n", proposalName)
		fmt.Println("If this transaction is successful, have the other party approve and execute the multisig proposal to an atomic swap.")
		fmt.Println("Review this proposal with:")
		fmt.Printf("  eosc multisig review %s %s", soldAccount, proposalName)
		fmt.Println("")
		msigPermissions := []zsw.PermissionLevel{buyerPerm, myPerm, zsw.PermissionLevel{Actor: soldAccount, Permission: zsw.PermissionName("owner")}}
		pushEOSCActions(ctx, api, msig.NewPropose(soldAccount, zsw.Name(proposalName), msigPermissions, tx))

	},
}

func init() {
	toolsCmd.AddCommand(toolsSellAccountCmd)

	toolsSellAccountCmd.Flags().StringP("memo", "", "", "Memo message to attach to transfer")
	toolsSellAccountCmd.Flags().StringP("proposal-name", "", "sellaccount", "Proposal name to use in the eosio.msig contract")
	toolsSellAccountCmd.Flags().StringP("buyer-permission", "", "", "Permission required of the buyer (to authorized 'eosio.token::transfer')")
	toolsSellAccountCmd.Flags().StringP("seller-permission", "", "", "Permission required of the seller (you, to authorize 'eosio::updateauth')")
	toolsSellAccountCmd.Flags().DurationP("sale-expiration", "", 1*time.Hour, "Expire proposed transaction after this amount of time (30m, 1h, etc..)")
}

func sellAccountFindAuthority(data *zsw.AccountResp, targetPerm string) (zsw.Authority, error) {
	for _, perm := range data.Permissions {
		if perm.PermName == targetPerm {
			return perm.RequiredAuth, nil
		}
	}
	return zsw.Authority{}, fmt.Errorf("permission %q not found in account %q", targetPerm, data.AccountName)
}
