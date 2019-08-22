package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

//func init() {
//	trafficCmd.AddCommand(cleanRulesCmd)
//	trafficCmd.AddCommand(setHeadersCmd)
//}

var trafficCmd = &cobra.Command{
	Use:   "traffic",
	Short: "Manage istio's traffic rules",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Usage()
		os.Exit(1)
	},
}

func init() {
	RootCmd.AddCommand(trafficCmd)
}
