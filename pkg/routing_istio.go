package pkg

import (
	"crypto/sha256"
	"fmt"
	"github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"reflect"

	"github.com/pismo/istiops/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IstioOperationsInterface set IstiOps interface for handling routing
type IstioOperationsInterface interface {
	Headers(cid string, labels map[string]string, headers map[string]string) error
	Percentage(cid string, labels map[string]string, percentage int32) error
}

// GetAllVirtualServices returns all istio resources 'virtualservices'
func GetAllVirtualServices(cid string, namespace string) (virtualServiceList *v1alpha3.VirtualServiceList, error error) {
	utils.Info(fmt.Sprintf("Getting all virtualservices..."), cid)
	vss, err := istioClient.NetworkingV1alpha3().VirtualServices(namespace).List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return vss, nil
}

// GetVirtualService returns a single virtualService object given a name & namespace
func GetVirtualService(cid string, name string, namespace string) (virtualService *v1alpha3.VirtualService, error error) {
	utils.Info(fmt.Sprintf("Getting virtualService '%s' to update...", name), cid)
	vs, err := istioClient.NetworkingV1alpha3().VirtualServices(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return vs, nil
}

// GetAllVirtualservices returns all istio resources 'virtualservices'
func GetAllDestinationRules(cid string, namespace string) (destinationRuleList *v1alpha3.DestinationRuleList, error error) {
	utils.Info(fmt.Sprintf("Getting all destinationrules..."), cid)
	drs, err := istioClient.NetworkingV1alpha3().DestinationRules(namespace).List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return drs, nil
}

// UpdateVirtualService updates a specific virtualService given an updated object
func UpdateVirtualService(cid string, name string, namespace string, virtualService *v1alpha3.VirtualService) error {
	utils.Info(fmt.Sprintf("Updating rules for virtualService '%s'...", name), cid)
	_, err := istioClient.NetworkingV1alpha3().VirtualServices(namespace).Update(virtualService)
	if err != nil {
		return err
	}
	return nil
}

// GetDestinationRules returns a single destinationRule object given a name & namespace
func GetDestinationRule(cid string, name string, namespace string) (destinationRule *v1alpha3.DestinationRule, error error) {
	utils.Info(fmt.Sprintf("Getting destinationRule '%s' to update...", name), cid)
	dr, err := istioClient.NetworkingV1alpha3().DestinationRules(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return dr, nil
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

func GetResourcesToUpdate(cid string, v IstioValues, labels map[string]string) (istioResources []*IstioResources, error error) {
	vss, err := GetAllVirtualServices(cid, v.Namespace)
	if err != nil {
		return nil, err
	}

	drs, err := GetAllDestinationRules(cid, v.Namespace)
	if err != nil {
		return nil, err
	}

	// iterate every cluster destinationRule

	var resourcesToUpdate []*IstioResources

	for _, dr := range drs.Items {
		destinationRuleName := fmt.Sprintf("%s", dr.Name)

		// checking if destination_rule key is already created for resourcesToUpdate
		utils.Debug(fmt.Sprintf("Checking subset rules for Destination Rule '%s'...", destinationRuleName), cid)
		for _, subset := range dr.Spec.Subsets {

			// checking if the DR subset map (subset.Labels) matches the one provided by Interface client (labels)
			if reflect.DeepEqual(subset.Labels, labels) {
				// find virtualservices which have subset.Name from DestinationRule
				utils.Info(fmt.Sprintf("Found rule '%s' from Destination Rule '%s' which matches provided label selector!", subset.Name, destinationRuleName), cid)

				for _, vs := range vss.Items {
					virtualServiceName := fmt.Sprintf("%s", vs.Name)

					utils.Debug(fmt.Sprintf("Checking subset rules for virtualservice '%s'...", virtualServiceName), cid)
					for _, match := range vs.Spec.Http {
						for _, route := range match.Route {
							if route.Destination.Subset == subset.Name {
								// In case of a non-existent key 'virtualServiceName', create it
								resourcesToUpdate = append(resourcesToUpdate, &IstioResources{
									DestinationRule: dr,
									VirtualService: vs,
								})
							}
						}
					}
				}
			}
		}
	}

	return resourcesToUpdate, nil
}

// Percentage set percentage as routing-match strategy for istio resources
func (v IstioValues) Percentage(cid string, labels map[string]string, percentage int32) error {
	//fmt.Println(v.Name, v.Build, labels)

	return nil
}

// Headers set headers as routing-match strategy for istio resources
func (v IstioValues) Headers(cid string, labels map[string]string, headers map[string]string) error {
	resourcesToUpdate, err := GetResourcesToUpdate(cid, v, labels)
	for _, value := range resourcesToUpdate {
		fmt.Println(value)
	}
	if err != nil {
		return err
	}
	//
	//for resourceName, resource := range resourcesToUpdate {
	//
	//	if fmt.Sprintf("%s", resource.Resource) == "VirtualService" {
	//		vs, err := GetVirtualService(cid, resourceName, v.Namespace)
	//		if err != nil {
	//			utils.Fatal(fmt.Sprintf("Could not find virtualService '%s' due to error '%s'", resourceName, err), cid)
	//		}
	//
	//		for _, value := range vs.Spec.Http {
	//			for _, matchValue := range value.Match {
	//				if matchValue.Uri != nil {
	//					fmt.Println(matchValue.Uri)
	//					matchValue.Uri = &v1alpha3.StringMatch{MatchType: &v1alpha3.StringMatch_Regex{Regex: "////"}}
	//
	//					err := UpdateVirtualService(cid, resourceName, vs.Namespace, vs)
	//					if err != nil {
	//						utils.Fatal(fmt.Sprintf("Could not update virtualService '%s' due to error '%s'", resourceName, err), cid)
	//					}
	//				}
	//				fmt.Println(matchValue.Uri)
	//			}
	//		}
	//	}
	//
	//	if fmt.Sprintf("%s", resource.Resource) == "DestinationRule" {
	//		dr, err := GetDestinationRule(cid, resourceName, v.Namespace)
	//		if err != nil {
	//			utils.Fatal(fmt.Sprintf("Could not find destinationRule '%s' due to error '%s'", resourceName, err), cid)
	//		}
	//
	//		for _, value := range dr.Spec.Subsets {
	//			fmt.Println(value.Name)
	//		}
	//	}
	//
	//	fmt.Println("")
	//}
	return nil
}
