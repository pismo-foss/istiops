package cmd

import (
	"github.com/pismo/istiops/utils"

	"github.com/spf13/cobra"
)

func init() {
	trafficCmd.AddCommand(cleanRulesCmd)

	cleanRulesCmd.Flags().StringP("label-selector", "l", "", "LabelSelector. Ex: app=api-foo,build=3")
	cleanRulesCmd.MarkFlagRequired("label-selector")
}

var cleanRulesCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all rules except the main one (URI set)",
	Run: func(cmd *cobra.Command, args []string) {
		labelSelector, err := cmd.Flags().GetString("label-selector")
		if labelSelector == "" {
			utils.Fatal("empty label", CID)
		}
		if err != nil {
			utils.Fatal("Failed when getting label selector", CID)
		}

		istioResult := istiops.ClearRules(CID, map[string]string{"environment": "pipeline-go"})
		if istioResult != nil {
			utils.Fatal("Failed", CID)
		}
	},
}
