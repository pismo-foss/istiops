package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	trafficCmd.AddCommand(rulesClearCmd)

	trafficCmd.PersistentFlags().StringP("namespace", "n", "default", "kubernetes' cluster namespace")
	trafficCmd.PersistentFlags().StringP("destination", "d", "", "* destination's hostname ('api.domain.io' or 'k8s-service')")
	trafficCmd.PersistentFlags().Uint32P("port", "p", 0, "* destination's port")
	trafficCmd.PersistentFlags().StringP("label-selector", "l", "", "* labels selector to filter istio' resources")
	trafficCmd.PersistentFlags().StringP("headers", "e", "", "request headers to filter routing destination")
	trafficCmd.PersistentFlags().Uint32P("weight", "w", 0, "weight (percentage) of routing")

	_ = trafficCmd.MarkPersistentFlagRequired("destination")
	_ = trafficCmd.MarkPersistentFlagRequired("port")
	_ = trafficCmd.MarkPersistentFlagRequired("label-selector")
}

var trafficCmd = &cobra.Command{
	Use:   "traffic",
	Short: "Manage istio's traffic rules",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(">")
		fmt.Println(cmd.Flag("namespace").Value)

		//_ = cmd.Usage()
		//os.Exit(1)
	},
}
