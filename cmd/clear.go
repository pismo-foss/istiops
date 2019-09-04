package cmd

import (
	"fmt"
	"github.com/pismo/istiops/pkg/router"
	"github.com/spf13/cobra"
)

func init() {
	rulesClearCmd.PersistentFlags().StringP("namespace", "n", "default", "kubernetes' cluster namespace")
	rulesClearCmd.PersistentFlags().StringP("label-selector", "l", "", "* labels selector to filter istio' resources")

	_ = rulesClearCmd.MarkPersistentFlagRequired("namespace")
	_ = rulesClearCmd.MarkPersistentFlagRequired("label-selector")
}

var rulesClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Removes all rules & routes",
	Run: func(cmd *cobra.Command, args []string) {

		namespace := cmd.Flag("namespace").Value.String()
		if namespace == "" {
			namespace = "default"
		} else {
			namespace = cmd.Flag("namespace").Value.String()
		}

		mappedLabelSelector, err := router.Mapify(trackingId, fmt.Sprintf("%s", cmd.Flag("label-selector").Value))
		if err != nil {
			fmt.Println(err)
		}

		drR := &router.DestinationRule{
			TrackingId: trackingId,
			Namespace:  namespace,
			Istio:      client,
		}

		vsR := &router.VirtualService{
			TrackingId: trackingId,
			Namespace:  namespace,
			Istio:      client,
		}

		shift := router.Shift{
			Selector: mappedLabelSelector,
		}

		op := operator(drR, vsR)
		err = op.Clear(shift)
		if err != nil {
			fmt.Println(err)
		}
	},
}
