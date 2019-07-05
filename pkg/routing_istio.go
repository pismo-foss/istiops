package pkg

import (
	"fmt"
	"reflect"

	"github.com/pismo/istiops/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IstioOperationsInterface interface {
	Headers(cid string, labels map[string]string, headers map[string]string) error
	Percentage(cid string, labels map[string]string, percentage int32) error
}

func (v IstioValues) Headers(cid string, labels map[string]string, headers map[string]string) error {
	fmt.Println(v.Name, v.Build, labels)
	return nil
}

func (v IstioValues) Percentage(cid string, labels map[string]string, percentage int32) error {
	utils.Info(fmt.Sprintf("Getting all virtualservices..."), cid)
	vss, err := istioClient.NetworkingV1alpha3().VirtualServices(v.Namespace).List(v1.ListOptions{
		FieldSelector: "",
	})
	if err != nil {
		return err
	}

	utils.Info(fmt.Sprintf("Getting all destinationrules..."), cid)
	drs, err := istioClient.NetworkingV1alpha3().DestinationRules(v.Namespace).List(v1.ListOptions{})
	if err != nil {
		return err
	}

	// iterate every cluster destinationRule

	resourcesToUpdate := map[string]*IstioResource{}

	for _, dr := range drs.Items {
		destinationRuleName := fmt.Sprintf("%s", dr.Name)

		// checking if destination_rule key is already created for resourcesToUpdate
		utils.Debug(fmt.Sprintf("Checking subset rules for Destination Rule '%s'...", destinationRuleName), cid)
		for _, subset := range dr.Spec.Subsets {
			// checking if the DR subset map (subset.Labels) matches the one provided by Interface client (labels)
			if reflect.DeepEqual(subset.Labels, labels) {
				// find virtualservices which have subset.Name from DestinationRule
				if _, status := resourcesToUpdate[destinationRuleName]; status != true {
					resourcesToUpdate[destinationRuleName] = &IstioResource{"api-xpto-destination-rule", []string{}}
				}
				utils.Info(fmt.Sprintf("Found rule '%s' from Destination Rule '%s' which matches provided label selector!", subset.Name, destinationRuleName), cid)
				resourcesToUpdate[destinationRuleName].Items = append(resourcesToUpdate[destinationRuleName].Items, subset.Name)

				for _, vs := range vss.Items {
					virtualServiceName := fmt.Sprintf("%s", vs.Name)

					utils.Debug(fmt.Sprintf("Checking subset rules for virtualservice '%s'...", virtualServiceName), cid)
					for _, match := range vs.Spec.Http {
						for _, route := range match.Route {
							if route.Destination.Subset == subset.Name {
								fmt.Println("Found virtualservice match!")
								if _, status := resourcesToUpdate[virtualServiceName]; status != true {
									resourcesToUpdate[virtualServiceName] = &IstioResource{"api-xpto-destination-rule", []string{}}
								}
								resourcesToUpdate[virtualServiceName].Items = append(resourcesToUpdate[virtualServiceName].Items, route.Destination.Subset)
							}
						}
					}
				}

			}
		}
	}

	fmt.Println(resourcesToUpdate)

	return nil
}
