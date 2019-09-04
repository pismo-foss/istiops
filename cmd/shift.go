package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	shiftCmd.PersistentFlags().StringP("namespace", "n", "default", "kubernetes' cluster namespace")
	shiftCmd.PersistentFlags().StringP("destination", "d", "", "* destination's hostname ('api.domain.io' or 'k8s-service')")
	shiftCmd.PersistentFlags().Uint32P("port", "p", 0, "* destination's port")
	shiftCmd.PersistentFlags().StringP("label-selector", "l", "", "* labels selector to filter istio' resources")
	shiftCmd.PersistentFlags().Uint32P("weight", "w", 0, "* weight (percentage) of routing")

	_ = shiftCmd.MarkPersistentFlagRequired("destination")
	_ = shiftCmd.MarkPersistentFlagRequired("port")
	_ = shiftCmd.MarkPersistentFlagRequired("label-selector")
	_ = shiftCmd.MarkPersistentFlagRequired("weight")
}

var shiftCmd = &cobra.Command{
	Use:   "shift",
	Short: "Shift istio's traffic",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(">")
		fmt.Println(cmd.Flag("namespace").Value)

		//_ = cmd.Usage()
		//os.Exit(1)
	},
}
