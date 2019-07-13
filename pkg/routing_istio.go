package pkg

import (
	"crypto/sha256"
	"fmt"
	"github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/pismo/istiops/utils"
	v1alpha32 "istio.io/api/networking/v1alpha3"
	"reflect"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IstioOperationsInterface set IstiOps interface for handling routing
type IstioOperationsInterface interface {
	Headers(cid string, labels map[string]string, headers map[string]string) error
	Percentage(cid string, labels map[string]string, percentage int32) error
}

// GetAllVirtualServices returns all istio resources 'virtualservices'
func GetAllVirtualServices(cid string, namespace string, listOptions metav1.ListOptions) (virtualServiceList *v1alpha3.VirtualServiceList, error error) {
	utils.Info(fmt.Sprintf("Getting all virtualservices..."), cid)
	vss, err := istioClient.NetworkingV1alpha3().VirtualServices(namespace).List(listOptions)
	if err != nil {
		return nil, err
	}

	return vss, nil
}

// GetVirtualService returns a single virtualService object given a name & namespace
func GetVirtualService(cid string, name string, namespace string, getOptions metav1.GetOptions) (virtualService *v1alpha3.VirtualService, error error) {
	utils.Info(fmt.Sprintf("Getting virtualService '%s' to update...", name), cid)
	vs, err := istioClient.NetworkingV1alpha3().VirtualServices(namespace).Get(name, getOptions)
	if err != nil {
		return nil, err
	}
	return vs, nil
}

// GetAllVirtualservices returns all istio resources 'virtualservices'
func GetAllDestinationRules(cid string, namespace string, listOptions metav1.ListOptions) (destinationRuleList *v1alpha3.DestinationRuleList, error error) {
	utils.Info(fmt.Sprintf("Getting all destinationrules..."), cid)
	drs, err := istioClient.NetworkingV1alpha3().DestinationRules(namespace).List(listOptions)
	if err != nil {
		return nil, err
	}

	return drs, nil
}

// GetDestinationRules returns a single destinationRule object given a name & namespace
func GetDestinationRule(cid string, name string, namespace string, getOptions metav1.GetOptions) (destinationRule *v1alpha3.DestinationRule, error error) {
	utils.Info(fmt.Sprintf("Getting destinationRule '%s' to update...", name), cid)
	dr, err := istioClient.NetworkingV1alpha3().DestinationRules(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return dr, nil
}

// UpdateVirtualService updates a specific virtualService given an updated object
func UpdateVirtualService(cid string, namespace string, virtualService *v1alpha3.VirtualService) error {
	utils.Info(fmt.Sprintf("Updating rule for virtualService '%s'...", virtualService.Name), cid)
	_, err := istioClient.NetworkingV1alpha3().VirtualServices(namespace).Update(virtualService)
	if err != nil {
		return err
	}
	return nil
}

// UpdateDestinationRule updates a specific virtualService given an updated object
func UpdateDestinationRule(cid string, namespace string, destinationRule *v1alpha3.DestinationRule) error {
	utils.Info(fmt.Sprintf("Updating rule for destinationRule '%s'...", destinationRule.Name), cid)
	_, err := istioClient.NetworkingV1alpha3().DestinationRules(namespace).Update(destinationRule)
	if err != nil {
		return err
	}
	return nil
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


// GetResourcesToUpdate returns a slice of all DestinationRules and/or VirtualServices (based on given labelSelectors to a posterior update
func GetResourcesToUpdate(cid string, v IstioValues, labels map[string]string) (istioResources []*IstioResources, error error) {
	listOptions := metav1.ListOptions{}
	vss, err := GetAllVirtualServices(cid, v.Namespace, listOptions)
	if err != nil {
		return nil, err
	}

	drs, err := GetAllDestinationRules(cid, v.Namespace, listOptions)
	if err != nil {
		return nil, err
	}

	// iterate every cluster destinationRule

	var resourcesToUpdate []*IstioResources
	resourcesToUpdate = nil

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
									DestinationRule: IstioMatchedDestinationRule{
										subset.Name,
										dr,
									},
									VirtualService: IstioMatchedVirtualService{
										route.Destination.Subset,
										vs,
									},
								})
							}
						}
					}
				}
			}
		}
	}

	if resourcesToUpdate == nil {
		utils.Info(fmt.Sprintf("Couldn't find any istio resources based on given labelsSelector to update. "), cid)
	}

	return resourcesToUpdate, nil
}

// RemoveSubsetRule removes an entry from Subsets[] slice based on given subsetIndex
func RemoveSubsetRule(subsets []*v1alpha32.Subset, subsetIndex int) ([]*v1alpha32.Subset, error) {
	copy(subsets[subsetIndex:], subsets[subsetIndex+1:])
	subsets[len(subsets)-1] = &v1alpha32.Subset{}

	return subsets[:len(subsets)-1], nil
}

// RemoveVirtualServiceHttpRoute removes an entry from Subsets[] slice based on given subsetIndex
func RemoveVirtualServiceHttpRoute(http []*v1alpha32.HTTPRoute, httpRouteIndex int) ([]*v1alpha32.HTTPRoute, error){
	copy(http[httpRouteIndex:], http[httpRouteIndex+1:])
	http[len(http)-1] = &v1alpha32.HTTPRoute{}

	return http[:len(http)-1], nil
}

// CreateNewVirtualServiceHttpRoute returns an existent VirtualService with a new basic HTTP route appended to it
func CreateNewVirtualServiceHttpRoute(cid string, virtualService *v1alpha3.VirtualService, matchHeaders map[string]string, destinationHost string, destinationPort int64) error{

	utils.Info(fmt.Sprintf("Creating new http route for virtualService '%s'...", virtualService.Name), cid)
	fmt.Println(matchHeaders)
	fmt.Println(destinationHost)
	fmt.Println(destinationPort)
	for _, httpValue := range virtualService.Spec.Http {
		for _, httpMatch := range httpValue.Match {
			fmt.Println(httpMatch)
		}
	}

	return nil
}

// Percentage set percentage as routing-match strategy for istio resources
func (v IstioValues) Percentage(cid string, labels map[string]string, percentage int32) error {
	//fmt.Println(v.Name, v.Build, labels)
	//for _, httpRules := range resource.VirtualService.Item.Spec.Http {
	//	for _, matchValue := range httpRules.Route {
	//		if matchValue.Destination.Subset == resource.DestinationRule.Name {
	//			matchValue.Weight = 30
	//			fmt.Println(matchValue)
	//			err := UpdateVirtualService(cid, resource.VirtualService.Subset, v.Namespace, &resource.VirtualService.Item)
	//			if err != nil {
	//				utils.Fatal(fmt.Sprintf("Could not update virtualService '%s' due to error '%s'", resource.VirtualService.Item.Name, err), cid)
	//			}
	//
	//			fmt.Println("Set destinationrule subset to virtualservice")
	//		}
	//	}
	//}
	return nil
}

// Headers set headers as routing-match strategy for istio resources
func (v IstioValues) Headers(cid string, labels map[string]string, headers map[string]string) error {
	replacer := strings.NewReplacer(".", "", "-", "", "/", "")
	simplifiedVersion := replacer.Replace(v.Version)
	simplifiedVersion = strings.ToLower(simplifiedVersion)
	subsetRuleName := fmt.Sprintf("%s-%d", simplifiedVersion, v.Build)

	resourcesToUpdate, err := GetResourcesToUpdate(cid, v, labels)
	if err != nil {
		return err
	}
	//
	for _, resource := range resourcesToUpdate {
		for subsetKey, subset := range resource.DestinationRule.Item.Spec.Subsets {
			// If rule already exists recreate it
			if subset.Name == subsetRuleName {
				cleanedSubsets, err := RemoveSubsetRule(resource.DestinationRule.Item.Spec.Subsets, subsetKey)
				if err != nil {
					utils.Fatal(fmt.Sprintf("Could not recreate subsetRule '%s'", subset.Name), cid)
				}
				resource.DestinationRule.Item.Spec.Subsets = cleanedSubsets
			}
		}

		// Create DestinationRule entry for specified labels & apply it
		resource.DestinationRule.Item.Spec.Subsets = append(
			resource.DestinationRule.Item.Spec.Subsets,
			&v1alpha32.Subset{
				Name:   subsetRuleName,
				Labels: headers,
			})
		err := UpdateDestinationRule(cid, v.Namespace, &resource.DestinationRule.Item)
		if err != nil {
			utils.Fatal(fmt.Sprintf("Could not update destinationRule '%s' due to error '%s'", resource.DestinationRule.Name, err), cid)
		}

		// Search for virtualservice's rule which matches subset name to append headers routing to it
		for _, httpRules := range resource.VirtualService.Item.Spec.Http {
			for _, matchValue := range httpRules.Route {
				if matchValue.Destination.Subset == resource.DestinationRule.Name {
					errT := CreateNewVirtualServiceHttpRoute(cid, &resource.VirtualService.Item, map[string]string{}, "api-host", 8080)
					if errT != nil {
						utils.Fatal(fmt.Sprintf("'%s'", err), cid)
					}
					matchValue.Weight = 30
					fmt.Println(matchValue)
					err := UpdateVirtualService(cid, v.Namespace, &resource.VirtualService.Item)
					if err != nil {
						utils.Fatal(fmt.Sprintf("Could not update virtualService '%s' due to error '%s'", resource.VirtualService.Item.Name, err), cid)
					}

					fmt.Println("Set destinationrule subset to virtualservice")
				}
			}
		}
	}
	return nil
}
