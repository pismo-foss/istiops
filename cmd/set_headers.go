package cmd

import (
	"fmt"
	"github.com/pismo/istiops/utils"
	"github.com/spf13/cobra"
)

func init() {
	setHeadersCmd.Flags().StringP("destination", "d", "", "Destination hostname. Ex: api-foo-service || api-foo.domain.io")
	setHeadersCmd.Flags().Uint32P("port", "p", 0, "Destination port. Ex: 8080")
	setHeadersCmd.Flags().StringP("headers", "H", "", "LabelSelector. Ex: app=api-foo,build=3")
	setHeadersCmd.Flags().StringP("label-selector", "l", "", "LabelSelector. Ex: app=api-foo,build=3")

	setHeadersCmd.MarkFlagRequired("destination")
	setHeadersCmd.MarkFlagRequired("port")
	setHeadersCmd.MarkFlagRequired("headers")
	setHeadersCmd.MarkFlagRequired("label-selector")
}

var setHeadersCmd = &cobra.Command{
	Use:   "headers",
	Short: "Create new header's canary release",
	Run: func(cmd *cobra.Command, args []string) {
		labelSelector, err := cmd.Flags().GetString("label-selector")
		if labelSelector == "" {
			utils.Fatal("empty label", CID)
		}
		if err != nil {
			utils.Fatal("Failed when getting label selector", CID)
		}

		headers, err := cmd.Flags().GetString("headers")
		if err != nil {
			utils.Fatal("Failed when getting headers", CID)
		}

		mapifiedLabels, err := utils.MapifyLabels(CID, headers)
		fmt.Println(mapifiedLabels)

		//_, istioResult := istiops.SetHeaders(CID, mapifiedLabels)
		//if istioResult != nil {
		//	utils.Fatal("Failed", CID)
		//}
	},
}
