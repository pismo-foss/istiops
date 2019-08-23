package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rulesClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Removes all rules & routes",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Clearing")
	},
}
