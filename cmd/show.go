package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	showCmd.PersistentFlags().StringP("namespace", "n", "default", "kubernetes' cluster namespace")
	showCmd.PersistentFlags().StringP("label-selector", "l", "", "* labels selector to filter istio' resources")

	_ = showCmd.MarkPersistentFlagRequired("label-selector")
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show currentistio's traffic rules",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(">")
		fmt.Println(cmd.Flag("namespace").Value)
	},
}
