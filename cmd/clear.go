package cmd

import (
	"fmt"
	"github.com/pismo/istiops/pkg/router"
	"github.com/spf13/cobra"
)

func init() {
	rulesClearCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "kubernetes' cluster namespace")
	rulesClearCmd.PersistentFlags().StringVarP(&labelSelector, "label-selector", "l", "", "* labels selector to filter istio' resources")

	_ = rulesClearCmd.MarkPersistentFlagRequired("label-selector")
}

var rulesClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Removes all rules & routes",
	Run: func(cmd *cobra.Command, args []string) {
		mappedLabelSelector, err := router.Mapify(trackingId, fmt.Sprintf("%s", cmd.Flag("label-selector").Value))
		if err != nil {
			fmt.Println(err)
		}

		shift.Selector.Labels = mappedLabelSelector

		err = op.Clear(shift)
		if err != nil {
			fmt.Println(err)
		}
	},
}
