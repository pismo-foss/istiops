package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {
	trafficCmd.AddCommand(showCmd)
	trafficCmd.AddCommand(rulesClearCmd)
	trafficCmd.AddCommand(shiftCmd)
}

var trafficCmd = &cobra.Command{
	Use:   "traffic",
	Short: "Manage istio's traffic rules",
	Run: func(cmd *cobra.Command, args []string) {
		kubeContext, _ := rootCmd.Flags().GetString("context")
		kubeConfigPath, _ := rootCmd.Flags().GetString("kubeconfig")
		inCluster, _ := rootCmd.Flags().GetBool("in-cluster")
		clientSetup(kubeContext, kubeConfigPath, inCluster)

		_ = cmd.Usage()
		os.Exit(1)
	},
}
