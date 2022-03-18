package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"

	zsw "github.com/zhongshuwen/zswchain-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var toolsProducerJSONCmd = &cobra.Command{
	Use:   "producerjson [account] [file.json]",
	Short: "Publish a producer json file to a producerjson-compatible contract.",
	Long: `Publish a producer json file to a producerjson-compatible contract.

Reference: https://github.com/greymass/producerjson
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		producerAccount := toAccount(args[0], "sold account")
		fileName := args[1]

		cnt, err := ioutil.ReadFile(fileName)
		errorCheck("reading json file", err)

		// TODO: eventually do some validation on the producerjson content
		var packme map[string]interface{}
		errorCheck("file contains invalid json", json.Unmarshal(cnt, &packme))
		packedCnt, err := json.Marshal(packme)
		errorCheck("packing json more tightly", err)

		type producerJSONSet struct {
			Owner zsw.AccountName `json:"owner"`
			JSON  string          `json:"json"`
		}

		pushEOSCActions(context.Background(), api, &zsw.Action{
			Account: zsw.AccountName(viper.GetString("tools-producerjson-cmd-target-contract")),
			Name:    zsw.ActionName("set"),
			Authorization: []zsw.PermissionLevel{
				{Actor: producerAccount, Permission: zsw.PermissionName("active")},
			},
			ActionData: zsw.NewActionData(producerJSONSet{
				Owner: producerAccount,
				JSON:  string(packedCnt),
			}),
		})

	},
}

func init() {
	toolsCmd.AddCommand(toolsProducerJSONCmd)

	toolsProducerJSONCmd.Flags().StringP("target-contract", "", "producerjson", "Target producerjson contract")
}
