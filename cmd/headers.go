package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	headersCmd.PersistentFlags().StringP("namespace", "n", "default", "kubernetes' cluster namespace")
	headersCmd.PersistentFlags().StringP("destination", "d", "", "* destination's hostname ('api.domain.io' or 'k8s-service')")
	headersCmd.PersistentFlags().Uint32P("port", "p", 0, "* destination's port")
	headersCmd.PersistentFlags().StringP("label-selector", "l", "", "* labels selector to filter istio' resources")
	headersCmd.PersistentFlags().Uint32P("headers", "H", 0, "* request headers for routing")

	_ = headersCmd.MarkPersistentFlagRequired("namespace")
	_ = headersCmd.MarkPersistentFlagRequired("destination")
	_ = headersCmd.MarkPersistentFlagRequired("port")
	_ = headersCmd.MarkPersistentFlagRequired("label-selector")
	_ = headersCmd.MarkPersistentFlagRequired("headers")
}

var headersCmd = &cobra.Command{
	Use:   "headers",
	Short: "Route to certain pods given request-headers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(">")
		fmt.Println(cmd.Flag("namespace").Value)

		//_ = cmd.Usage()
		//os.Exit(1)
	},
}
