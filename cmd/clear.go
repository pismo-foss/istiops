package cmd

import (
	"fmt"

	"github.com/pismo/istiops/pkg/logger"
	"github.com/pismo/istiops/pkg/router"
	"github.com/spf13/cobra"
)

func init() {
	rulesClearCmd.PersistentFlags().StringP("namespace", "n", "default", "kubernetes' cluster namespace")
	rulesClearCmd.PersistentFlags().StringP("label-selector", "l", "", "* labels selector to filter istio' resources")
	rulesClearCmd.PersistentFlags().StringP("mode", "m", "soft", "if 'hard' all canary rules will be cleaned otherwise only canary rules with no pods will be cleaned")

	_ = rulesClearCmd.MarkPersistentFlagRequired("namespace")
	_ = rulesClearCmd.MarkPersistentFlagRequired("label-selector")
}

var rulesClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Removes all rules & routes except the master-one",
	Run: func(cmd *cobra.Command, args []string) {
		kubeContext, _ := rootCmd.Flags().GetString("context")
		kubeConfigPath, _ := rootCmd.Flags().GetString("kubeconfig")
		clientSetup(kubeContext, kubeConfigPath)

		namespace := cmd.Flag("namespace").Value.String()
		if namespace == "" {
			namespace = "default"
		} else {
			namespace = cmd.Flag("namespace").Value.String()
		}

		mappedLabelSelector, err := router.Mapify(trackingId, fmt.Sprintf("%s", cmd.Flag("label-selector").Value))
		if err != nil {
			logger.Fatal(fmt.Sprintf("%s", err), "cmd")
		}

		clearMode := cmd.Flag("mode").Value.String()
		if clearMode != "hard" {
			clearMode = "soft"
		} else {
			// enforce any value which is not 'hard' to it
			clearMode = "hard"
		}

		drR := &router.DestinationRule{
			TrackingId: trackingId,
			Namespace:  namespace,
			Istio:      clients.Istio,
			KubeClient: clients.Kubernetes,
		}

		vsR := &router.VirtualService{
			TrackingId: trackingId,
			Namespace:  namespace,
			Istio:      clients.Istio,
			KubeClient: clients.Kubernetes,
		}

		shift := router.Shift{
			Selector: mappedLabelSelector,
		}

		op := operator(drR, vsR)
		err = op.Clear(shift, clearMode)
		if err != nil {
			logger.Fatal(fmt.Sprintf("%s", err), "cmd")
		}
	},
}
