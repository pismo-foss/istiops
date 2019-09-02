package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	weightCmd.PersistentFlags().StringP("namespace", "n", "default", "kubernetes' cluster namespace")
	weightCmd.PersistentFlags().StringP("destination", "d", "", "* destination's hostname ('api.domain.io' or 'k8s-service')")
	weightCmd.PersistentFlags().Uint32P("port", "p", 0, "* destination's port")
	weightCmd.PersistentFlags().StringP("label-selector", "l", "", "* labels selector to filter istio' resources")
	weightCmd.PersistentFlags().Uint32P("weight", "w", 0, "* weight (percentage) of routing")

	_ = weightCmd.MarkPersistentFlagRequired("destination")
	_ = weightCmd.MarkPersistentFlagRequired("port")
	_ = weightCmd.MarkPersistentFlagRequired("label-selector")
	_ = weightCmd.MarkPersistentFlagRequired("weight")
}

var weightCmd = &cobra.Command{
	Use:   "weight",
	Short: "Set istio's traffic weight rules",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(">")
		fmt.Println(cmd.Flag("namespace").Value)

		//_ = cmd.Usage()
		//os.Exit(1)
	},
}
