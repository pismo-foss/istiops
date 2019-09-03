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
		namespace = fmt.Sprintf("%s", cmd.Flag("namespace").Value)
		mappedLabelSelector, err := router.Mapify(trackingId, fmt.Sprintf("%s", cmd.Flag("label-selector").Value))
		if err != nil {
			fmt.Println(err)
		}

		dr = &router.DestinationRule{
			TrackingId: trackingId,
			Namespace:  namespace,
			Istio:      client,
		}

		vs = &router.VirtualService{
			TrackingId: trackingId,
			Namespace:  namespace,
			Istio:      client,
		}

		shift := router.Shift{
			Selector: mappedLabelSelector,
		}

		op := operator(dr, vs)
		err = op.Clear(shift)
		if err != nil {
			fmt.Println(err)
		}
	},
}
