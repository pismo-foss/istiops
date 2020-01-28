package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/gookit/color"
	"github.com/pismo/istiops/pkg/logger"
	"github.com/pismo/istiops/pkg/router"
	"github.com/spf13/cobra"
	"istio.io/api/networking/v1alpha3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func init() {
	showCmd.PersistentFlags().StringP("namespace", "n", "default", "kubernetes' cluster namespace")
	showCmd.PersistentFlags().StringP("label-selector", "l", "", "* labels selector to filter istio' resources")
	showCmd.PersistentFlags().StringP("output", "o", "", "stdout format can be 'json', 'yaml' or 'pretty'")

	_ = showCmd.MarkPersistentFlagRequired("label-selector")
}

type Subset struct {
	Name   string
	Labels map[string]string
}

type Deployment struct {
	Name      string
	Namespace string
	Pods      int32
}

type Destination struct {
	Service    string
	Weight     int32
	Subset     Subset
	Routable   bool
	Deployment Deployment
}

type Routes struct {
	Match        []*v1alpha3.HTTPMatchRequest
	Destinations []Destination
}

type Resource struct {
	Name      string
	Namespace string
	Hosts     []string
	Routes    []*Routes
}

func structured(trackingId string, namespace string, irl router.IstioRouteList, kClient kubernetes.Clientset) []Resource {
	var r Resource
	var resourceList []Resource

	for _, vs := range irl.VList.Items {
		r = Resource{}

		r.Name = vs.Name
		r.Namespace = vs.Namespace
		r.Hosts = vs.Spec.Hosts

		for _, httpValue := range vs.Spec.Http {
			route := &Routes{}

			for _, matchValue := range httpValue.Match {
				route.Match = append(route.Match, matchValue)
			}

			// handle destination
			var currentWeight int32
			for _, httpRoute := range httpValue.Route {
				jr := Destination{}
				jr.Service = fmt.Sprintf("%s:%d", httpRoute.Destination.Host, httpRoute.Destination.Port.GetNumber())

				if httpRoute.Weight == 0 {
					currentWeight = 100
				} else {
					currentWeight = httpRoute.Weight
				}

				subsetExists := false
				jr.Routable = true
				for _, dr := range irl.DList.Items {

					// validate if subset is valid and routable
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
							jr.Subset.Name = js.Name
						}
					}

					if !subsetExists {
						jr.Subset.Name = httpRoute.Destination.Subset
						jr.Routable = false
					}

					// validate if there are any pods to be routed
					labelString, err := router.Stringify(trackingId, jr.Subset.Labels)
					dep, err := kClient.AppsV1().Deployments(namespace).List(v1.ListOptions{
						LabelSelector: labelString,
					})
					if err != nil {
						return []Resource{}
					}

					if len(dep.Items) == 1 {
						depItem := dep.Items[0]
						jr.Deployment.Name = depItem.Name
						jr.Deployment.Namespace = depItem.Namespace
						jr.Deployment.Pods = depItem.Status.ReadyReplicas
					}

				}

				jr.Weight = currentWeight
				route.Destinations = append(route.Destinations, jr)
			}

			r.Routes = append(r.Routes, route)
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

	fmt.Print(string(jsonData))

}

func yamlfy(resourceList []Resource) {
	var yamlData []byte
	yamlData, err := yaml.Marshal(resourceList)
	if err != nil {
		logger.Fatal(fmt.Sprintf("%s", err), trackingId)
	}

	fmt.Print(string(yamlData))
}

func beautified(resourceList []Resource) {
	for _, vs := range resourceList {
		fmt.Println("")
		fmt.Println("Resource: ", vs.Name)
		fmt.Println("Namespace: ", vs.Namespace)
		fmt.Println("client -> request to -> ", vs.Hosts)

		for _, route := range vs.Routes {
			for _, httpMatch := range route.Match {
				if httpMatch.Uri != nil {
					color.Green.Println("  \\_", httpMatch.Uri)
				}

				if len(httpMatch.Headers) > 0 {
					color.Cyan.Println("  \\_ Header")
					for headerKey, headerValue := range httpMatch.Headers {
						var escapedHeaderValue string

						if headerValue.GetExact() != "" {
							escapedHeaderValue = fmt.Sprintf("exact: %s", headerValue.GetExact())
						}

						if headerValue.GetRegex() != "" {
							escapedHeaderValue = fmt.Sprintf("regex: %s", headerValue.GetRegex())
						}

						color.Cyan.Println("      |- ", headerKey)
						color.Cyan.Println("      |- ", escapedHeaderValue)
					}
				}
			}

			// handle destinations
			fmt.Println("       \\_ Destination [k8s service]")
			for _, httpRoute := range route.Destinations {
				fmt.Println(fmt.Sprintf("         - %s [%s]", httpRoute.Service, httpRoute.Deployment.Name))

				if httpRoute.Deployment.Pods > 0 {
					color.Green.Println("            |- active pods: ", httpRoute.Deployment.Pods)
				} else {
					color.Red.Println("            |- NON-EXISTENT ACTIVE PODS:", httpRoute.Deployment.Pods)
				}

				fmt.Println(fmt.Sprintf("            \\_ %d %% of requests for pods with labels", httpRoute.Weight))

				for labelKey, labelValue := range httpRoute.Subset.Labels {
					fmt.Println(fmt.Sprintf("               |- %s: %s", labelKey, labelValue))
				}

				if !httpRoute.Routable {
					color.Red.Println("               |- NON-EXISTENT SUBSET", httpRoute.Subset.Name)
				}

			}
		}
		fmt.Println("--")
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
		if output == "" {
			output = "pretty"
		}

		if output != "yaml" && output != "json" && output != "pretty" {
			logger.Fatal(fmt.Sprintf("--output must be 'yaml', 'json' or 'pretty'"), trackingId)
		}

		mappedLabelSelector, err := router.Mapify(trackingId, fmt.Sprintf("%s", cmd.Flag("label-selector").Value))
		if err != nil {
			logger.Fatal(fmt.Sprintf("%s", err), "cmd")
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
		irl, err := op.Get(shift.Selector)
		if err != nil {
			logger.Fatal(fmt.Sprintf("%s", err), trackingId)
		}

		logger.Debug("Listing all current active routing rules", trackingId)
		resourceList := structured(trackingId, namespace, irl, *clients.Kubernetes)

		if output == "pretty" {
			beautified(resourceList)
		}

		if output == "yaml" {
			yamlfy(resourceList)
		}

		if output == "json" {
			jsonfy(resourceList)
		}
	},
}
