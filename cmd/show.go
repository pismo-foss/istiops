package cmd

import (
	"fmt"
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/pismo/istiops/pkg/logger"
	"github.com/pismo/istiops/pkg/router"
	"github.com/spf13/cobra"
)

func init() {
	showCmd.PersistentFlags().StringP("namespace", "n", "default", "kubernetes' cluster namespace")
	showCmd.PersistentFlags().StringP("label-selector", "l", "", "* labels selector to filter istio' resources")
	showCmd.PersistentFlags().StringP("output", "o", "", "stdout format, can be 'summarized' or 'beautified'")

	_ = showCmd.MarkPersistentFlagRequired("label-selector")
	_ = showCmd.MarkPersistentFlagRequired("output")
}

func beautified(vss []v1alpha32.VirtualService) {
	for _, vs := range vss {
		fmt.Println("--")
		fmt.Println(vs.Name)
		fmt.Println("Hosts: ", vs.Spec.Hosts)
		for _, httpValue := range vs.Spec.Http {
			for _, httpMatch := range httpValue.Match {
				if httpMatch.Uri != nil {
					fmt.Println("* Match")
					fmt.Println("  \\_", httpMatch.Uri)
				}

				if len(httpMatch.Headers) > 0 {
					fmt.Println("* Match")
					fmt.Println("  \\_ Headers")
					for headerKey, headerValue := range httpMatch.Headers {
						fmt.Println(fmt.Sprintf("      - %s: %s", headerKey, headerValue.GetExact()))
					}
				}
			}

			fmt.Println("      \\_ Destination")
			for _, httpRoute := range httpValue.Route {
				fmt.Println("         -", httpRoute.Destination.Subset)
				if httpRoute.Weight != 0 {
					fmt.Println(fmt.Sprintf("             \\_ %d %% of requests", httpRoute.Weight))
				}
			}
		}
	}
}

func summarized(vss []v1alpha32.VirtualService) {
	for _, vs := range vss {
		fmt.Println("--")
		fmt.Println(vs.Name, vs.Spec.Http)
	}
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current istio's traffic rules",
	Run: func(cmd *cobra.Command, args []string) {
		namespace = fmt.Sprintf("%s", cmd.Flag("namespace").Value)
		output := fmt.Sprintf("%s", cmd.Flag("output").Value)

		if output != "summarized" && output != "beautified" {
			logger.Fatal(fmt.Sprintf("--output must be 'summarized' or 'beautified'"), trackingId)
		}

		mappedLabelSelector, err := router.Mapify(trackingId, fmt.Sprintf("%s", cmd.Flag("label-selector").Value))
		if err != nil {
			fmt.Println(err)
		}

		drR = &router.DestinationRule{
			TrackingId: trackingId,
			Namespace:  namespace,
			Istio:      client,
		}

		vsR = &router.VirtualService{
			TrackingId: trackingId,
			Namespace:  namespace,
			Istio:      client,
		}

		shift := router.Shift{
			Selector: mappedLabelSelector,
		}

		op := operator(drR, vsR)
		vss, err := op.Get(shift.Selector)
		if err != nil {
			logger.Fatal(fmt.Sprintf("%s", err), trackingId)
		}

		if output == "beautified" {
			beautified(vss)
		}

		if output == "summarized" {
			summarized(vss)
		}
	},
}
