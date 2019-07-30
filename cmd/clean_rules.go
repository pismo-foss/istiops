package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	var labels string

	trafficCmd.AddCommand(cleanRulesCmd)

	cleanRulesCmd.Flags().StringVarP(&labels, "label-selector", "l", "", "LabelSelector. Ex: app=api-foo,build=3")
	cleanRulesCmd.MarkFlagRequired("label-selector")
}

var trafficCmd = &cobra.Command{
	Use:   "traffic",
	Short: "Manage istio's traffic rules",
	Long:  `Use it`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("%s", args))
	},
}

var cleanRulesCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all rules except the main one (URI set)",
	Long:  `Use it`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inside subCmd Run with args: %v\n", args)
	},
}
