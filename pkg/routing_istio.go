package pkg

import (
	"crypto/sha256"
	"fmt"
	"reflect"

	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"

	"github.com/pismo/istiops/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IstioOperationsInterface set IstiOps interface for handling routing
type IstioOperationsInterface interface {
	Headers(cid string, labels map[string]string, headers map[string]string) error
	Percentage(cid string, labels map[string]string, percentage int32) error
}

// GetAllVirtualServices returns all istio resources 'virtualservices'
func GetAllVirtualServices(cid string, namespace string) (virtualServiceList *v1alpha32.VirtualServiceList, error error) {
	utils.Info(fmt.Sprintf("Getting all virtualservices..."), cid)
	vss, err := istioClient.NetworkingV1alpha3().VirtualServices(namespace).List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return vss, nil
}

// GetAllVirtualservices returns all istio resources 'virtualservices'
func GetAllDestinationRules(cid string, namespace string) (destinationRuleList *v1alpha32.DestinationRuleList, error error) {
	utils.Info(fmt.Sprintf("Getting all destinationrules..."), cid)
	drs, err := istioClient.NetworkingV1alpha3().DestinationRules(namespace).List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return drs, nil
}

// GenerateShaFromMap returns a slice of hashes (sha256) for every key:value in given map[string]string
func GenerateShaFromMap(mapToHash map[string]string) ([]string, error) {
	var mapHashes []string

	for k, v := range mapToHash {
		keyValue := fmt.Sprintf("%s=%s", k, v)
		sha256 := sha256.Sum256([]byte(keyValue))
		mapHashes = append(mapHashes, fmt.Sprintf("%x", sha256))
	}

	return mapHashes, nil
}

// CompareMapsKeyPairsHash compares two string maps using sha256, if both have identical content (no order oriented) it will return true
func CompareMapsKeyPairsHash(mapOne map[string]string, mapTwo map[string]string) bool {

	mapHashOne, err := GenerateShaFromMap(mapOne)
	if err != nil {
		return false
	}

	mapHashTwo, err := GenerateShaFromMap(mapTwo)
	if err != nil {
		return false
	}

	fmt.Println(mapHashOne)
	fmt.Println(mapHashTwo)

	// now we are going to compare it self each content for both slices
	// ... perhaps using `key` instead of iterating N slice positions?

	return true
}

// Percentage set percentage as routing-match strategy for istio resources
func (v IstioValues) Percentage(cid string, labels map[string]string, percentage int32) error {
	fmt.Println(v.Name, v.Build, labels)
	return nil
}

// Headers set headers as routing-match strategy for istio resources
func (v IstioValues) Headers(cid string, labels map[string]string, headers map[string]string) error {
	vss, err := GetAllVirtualServices(cid, v.Namespace)
	if err != nil {
		return err
	}

	drs, err := GetAllDestinationRules(cid, v.Namespace)
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
			if isMapEqual := CompareMapsKeyPairsHash(subset.Labels, labels); isMapEqual {
				fmt.Println("Is equal! %s", isMapEqual)
			}

			if reflect.DeepEqual(subset.Labels, labels) {
				// find virtualservices which have subset.Name from DestinationRule
				if _, status := resourcesToUpdate[destinationRuleName]; status != true {
					resourcesToUpdate[destinationRuleName] = &IstioResource{subset.Name, []string{}}
				}
				utils.Info(fmt.Sprintf("Found rule '%s' from Destination Rule '%s' which matches provided label selector!", subset.Name, destinationRuleName), cid)
				resourcesToUpdate[destinationRuleName].Items = append(resourcesToUpdate[destinationRuleName].Items, subset.Name)

				for _, vs := range vss.Items {
					virtualServiceName := fmt.Sprintf("%s", vs.Name)

					utils.Debug(fmt.Sprintf("Checking subset rules for virtualservice '%s'...", virtualServiceName), cid)
					for _, match := range vs.Spec.Http {
						for _, route := range match.Route {
							if route.Destination.Subset == subset.Name {
								if _, status := resourcesToUpdate[virtualServiceName]; status != true {
									resourcesToUpdate[virtualServiceName] = &IstioResource{route.Destination.Subset, []string{}}
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
