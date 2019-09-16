package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/gookit/color"
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

type Subset struct {
	Name   string
	Labels map[string]string
}

type Destination struct {
	Destination string
	Weight      int32
	Subset      Subset
}

type Resource struct {
	Name         string
	Namespace    string
	Hosts        []string
	Match        map[string]string
	Destinations []Destination
}

func structured(irl router.IstioRouteList) []Resource {
	var r Resource
	var resourceList []Resource

	for _, vs := range irl.VList.Items {
		r = Resource{}
		r.Match = map[string]string{}

		r.Name = vs.Name
		r.Namespace = vs.Namespace
		r.Hosts = vs.Spec.Hosts

		for _, httpValue := range vs.Spec.Http {
			// filtering virtualServices
			for _, httpMatch := range httpValue.Match {
				if httpMatch.Uri != nil {
					r.Match["regex"] = httpMatch.Uri.GetRegex()
				}

				if len(httpMatch.Headers) > 0 {
					for headerKey, headerValue := range httpMatch.Headers {
						r.Match["headers"] = fmt.Sprintf("%s:%s", headerKey, headerValue)
					}
				}
			}

			// handle destination
			var currentWeight int32
			for _, httpRoute := range httpValue.Route {
				jr := Destination{}
				jr.Destination = fmt.Sprintf("%s:%d", httpRoute.Destination.Host, httpRoute.Destination.Port.GetNumber())

				if httpRoute.Weight == 0 {
					currentWeight = 100
				} else {
					currentWeight = httpRoute.Weight
				}

				subsetExists := false
				for _, dr := range irl.DList.Items {
					for _, subset := range dr.Spec.Subsets {
						js := Subset{}
						js.Labels = map[string]string{}

						if subset.Name == httpRoute.Destination.Subset {
							subsetExists = true
							js.Name = subset.Name

							// append pod labels
							for labelKey, labelValue := range subset.Labels {
								js.Labels[labelKey] = labelValue
							}

							jr.Subset.Labels = js.Labels
						}
					}

					if !subsetExists {
						jr.Subset = Subset{
							Name:   "None",
							Labels: nil,
						}
					}
				}

				jr.Weight = currentWeight
				r.Destinations = append(r.Destinations, jr)
			}
		}

		resourceList = append(resourceList, r)
	}

	return resourceList
}

func jsonfy(resourceList []Resource) {
	var jsonData []byte
	jsonData, err := json.Marshal(resourceList)
	if err != nil {
		logger.Fatal(fmt.Sprintf("%s", err), trackingId)
	}

	fmt.Println(string(jsonData))

}

func yamlfy(resourceList []Resource) {
	var yamlData []byte
	yamlData, err := yaml.Marshal(resourceList)
	if err != nil {
		logger.Fatal(fmt.Sprintf("%s", err), trackingId)
	}

	fmt.Println(string(yamlData))
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
					color.Green.Println("  \\_", httpMatch.Uri)
				}

				if len(httpMatch.Headers) > 0 {
					//fmt.Println("* Match")
					color.Cyan.Println("  \\_ Headers")
					for headerKey, headerValue := range httpMatch.Headers {
						fmt.Println(color.Cyan.Sprintf("      |- %s: %s", headerKey, headerValue.GetExact()))
					}
				}
			}

			// handle destinations
			fmt.Println("       \\_ Destination [k8s service]")
			var currentWeight int32
			for _, httpRoute := range httpValue.Route {
				fmt.Println(fmt.Sprintf("         - %s:%d", httpRoute.Destination.Host, httpRoute.Destination.Port.GetNumber()))

				if httpRoute.Weight == 0 {
					currentWeight = 100
				} else {
					currentWeight = httpRoute.Weight
				}

				fmt.Println(fmt.Sprintf("            \\_ %d %% of requests for pods with labels", currentWeight))
				subsetExists := false
				for _, dr := range irl.DList.Items {
					for _, subset := range dr.Spec.Subsets {
						if subset.Name == httpRoute.Destination.Subset {
							subsetExists = true
							for labelKey, labelValue := range subset.Labels {
								fmt.Println(fmt.Sprintf("               |- %s: %s", labelKey, labelValue))
							}
						}
					}

					if !subsetExists {
						color.Red.Printf("               |- NON-EXISTENT SUBSET '%s'", httpRoute.Destination.Subset)
					}
				}
				fmt.Println("")
			}
			fmt.Println("")
		}
	}
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current istio's traffic rules",
	Run: func(cmd *cobra.Command, args []string) {

		namespace := cmd.Flag("namespace").Value.String()
		if namespace == "" {
			namespace = "default"
		} else {
			namespace = cmd.Flag("namespace").Value.String()
		}

		output := fmt.Sprintf("%s", cmd.Flag("output").Value)

		if output != "yaml" && output != "json" && output != "beautified" {
			logger.Fatal(fmt.Sprintf("--output must be 'yaml', 'json' or 'beautified'"), trackingId)
		}

		mappedLabelSelector, err := router.Mapify(trackingId, fmt.Sprintf("%s", cmd.Flag("label-selector").Value))
		if err != nil {
			fmt.Println(err)
		}

		drR := &router.DestinationRule{
			TrackingId: trackingId,
			Namespace:  namespace,
			Istio:      clients.Istio,
		}

		vsR := &router.VirtualService{
			TrackingId: trackingId,
			Namespace:  namespace,
			Istio:      clients.Istio,
		}

		shift := router.Shift{
			Selector: mappedLabelSelector,
		}

		op := operator(drR, vsR)
		irl, err := op.Get(shift.Selector)
		if err != nil {
			logger.Fatal(fmt.Sprintf("%s", err), trackingId)
		}

		logger.Info("Listing all current active routing rules", trackingId)
		resourceList := structured(irl)

		if output == "beautified" {
			beautified(irl)
		}

		if output == "yaml" {
			yamlfy(resourceList)
		}

		if output == "json" {
			jsonfy(resourceList)
		}
	},
}
