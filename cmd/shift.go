package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

func init() {
	shiftCmd.AddCommand(rulesClearCmd)
	shiftCmd.AddCommand(weightCmd)
	shiftCmd.AddCommand(headersCmd)
}

var shiftCmd = &cobra.Command{
	Use:   "shift",
	Short: "Shift istio's traffic rules",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Usage()
		os.Exit(1)
	},
}
