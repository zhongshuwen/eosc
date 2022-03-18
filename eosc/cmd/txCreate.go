// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"

	zsw "github.com/zhongshuwen/zswchain-go"
	"github.com/zhongshuwen/zswchain-go/system"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var txCreateCmd = &cobra.Command{
	Use:   "create [contract] [action] [payload]",
	Short: "Create a transaction with a single action",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		contract := toAccount(args[0], "contract")
		action := toActionName(args[1], "action")
		payload := args[2]

		forceUnique := viper.GetBool("tx-create-cmd-force-unique")

		var dump map[string]interface{}
		err := json.Unmarshal([]byte(payload), &dump)
		errorCheck("[payload] is not valid JSON", err)

		api := getAPI()
		actionBinary, err := api.ABIJSONToBin(ctx, contract, zsw.Name(action), dump)
		errorCheck("unable to retrieve action binary from JSON via API", err)

		actions := []*zsw.Action{
			&zsw.Action{
				Account:    contract,
				Name:       action,
				ActionData: zsw.NewActionDataFromHexData([]byte(actionBinary)),
			}}

		var contextFreeActions []*zsw.Action
		if forceUnique {
			contextFreeActions = append(contextFreeActions, newNonceAction())
		}

		pushEOSCActionsAndContextFreeActions(ctx, api, contextFreeActions, actions)
	},
}

func newNonceAction() *zsw.Action {
	return &zsw.Action{
		Account: zsw.AN("eosio.null"),
		Name:    zsw.ActN("nonce"),
		ActionData: zsw.NewActionData(system.Nonce{
			Value: hex.EncodeToString(generateRandomNonce()),
		}),
	}
}

func generateRandomNonce() []byte {
	// Use 48 bits of entropy to generate a valid random
	nonce := make([]byte, 6)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		errorCheck("unable to correctly generate nonce", err)
	}

	return nonce
}

func init() {
	txCmd.AddCommand(txCreateCmd)

	txCreateCmd.Flags().BoolP("force-unique", "f", false, "force the transaction to be unique. this will consume extra bandwidth and remove any protections against accidently issuing the same transaction multiple times")
}
