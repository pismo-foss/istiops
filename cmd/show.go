package cmd

import (
	"fmt"
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

func beautified(irl router.IstioRouteList) {
	// filtering virtualServices
	for _, vs := range irl.VList.Items {
		fmt.Println("--")
		fmt.Println("Resource: ", vs.Name)
		fmt.Println("")
		fmt.Println("client -> request to -> ", vs.Spec.Hosts)
		for _, httpValue := range vs.Spec.Http {
			for _, httpMatch := range httpValue.Match {
				if httpMatch.Uri != nil {
					//fmt.Println("* Match")
					fmt.Println("  \\_", httpMatch.Uri)
				}

				if len(httpMatch.Headers) > 0 {
					//fmt.Println("* Match")
					fmt.Println("  \\_ Headers")
					for headerKey, headerValue := range httpMatch.Headers {
						fmt.Println(fmt.Sprintf("      - %s: %s", headerKey, headerValue.GetExact()))
					}
				}
			}

			fmt.Println("      \\_ Destination [k8s service]")
			for _, httpRoute := range httpValue.Route {
				fmt.Println(fmt.Sprintf("         - %s:%d", httpRoute.Destination.Host, httpRoute.Destination.Port.GetNumber()))

				if httpRoute.Weight != 0 {
					fmt.Println(fmt.Sprintf("           \\_ %d %% of requests for pods with labels", httpRoute.Weight))
					for _, dr := range irl.DList.Items {
						for _, subset := range dr.Spec.Subsets {
							if subset.Name == httpRoute.Destination.Subset {
								for labelKey, labelValue := range subset.Labels {
									fmt.Println(fmt.Sprintf("               |- %s: %s", labelKey, labelValue))
								}
							}
						}
					}
					fmt.Println("                ---")
				}
			}
			fmt.Println("")
		}
	}
}

func summarized(irl router.IstioRouteList) {
	for _, vs := range irl.VList.Items {
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
		irl, err := op.Get(shift.Selector)
		if err != nil {
			logger.Fatal(fmt.Sprintf("%s", err), trackingId)
		}

		if output == "beautified" {
			beautified(irl)
		}

		if output == "summarized" {
			summarized(irl)
		}
	},
}
