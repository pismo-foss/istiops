package operator

import (
	"fmt"
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/pismo/istiops/pkg/client"
	"github.com/pismo/istiops/utils"
	"istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IstioOperator struct {
	TrackingId string
	Name       string
	Namespace  string
	Build      int8
	Client     *client.ClientSet
}

type IstioRoute struct {
	Port     uint32
	Hostname string
	Selector Selector
	Headers  map[string]string
	Weight   int32
}

type IstioRouteList struct {
	VirtualServiceList   *v1alpha32.VirtualServiceList
	DestinationRulesList *v1alpha32.DestinationRuleList
}

type Istiops interface {
	Create(ir *IstioRoute)
	Delete(ir *IstioRoute)
	Update(ir *IstioRoute) error
	Clear(Selector) error
}

type Selector struct {
	Labels map[string]string
}

type Headers struct {
	Labels map[string]string
}

func (ips *IstioOperator) Create(ir *IstioRoute) {
	fmt.Println("creating")

	istioNetworking := ips.Client.Istio.NetworkingV1alpha3()
	dr := &v1alpha32.DestinationRule{}
	dr.Name = ips.Name

	_, err := istioNetworking.DestinationRules(ips.Namespace).Create(dr)

	if err != nil {
		utils.Fatal(fmt.Sprintf("Could no create resource due to an error '%s'", err), ips.TrackingId)
	}

	fmt.Println("Creating something")
}

func (ips *IstioOperator) Delete(ir *IstioRoute) {
	fmt.Println("Initializing something")
	fmt.Println("", ips.TrackingId)
}

func (ips *IstioOperator) Update(ir *IstioRoute) error {
	if len(ir.Selector.Labels) == 0 {
		utils.Fatal(fmt.Sprintf("Labels must not be empty otherwise istiops won't be able to find any resources."), ips.TrackingId)
	}

	// Getting destination rules
	istioResources, err := GetResourcesToUpdate(ips, ir.Selector)
	if err != nil {
		utils.Fatal(fmt.Sprintf("Could not get istio resources to be updated due to an error '%s'", err), ips.TrackingId)
	}

	for _, dr := range istioResources.DestinationRulesList.Items {
		for _, subset := range dr.Spec.Subsets {
			matchedHeaders := 0
			for headerKey, headerValue := range ir.Headers {
				if _, ok := subset.Labels[headerKey]; ok {
					if subset.Labels[headerKey] == headerValue {
						matchedHeaders += 1
					}
				}
			}

			if matchedHeaders == len(ir.Headers) {
				utils.Fatal(fmt.Sprintf("Found existent subset for the current labels '%s'. Creating a new one for '%s'.", ir.Headers, dr.Name), ips.TrackingId)
			}

			// if no destinationRule matches provided headers, create it
			//if matchedHeaders != len(ir.Headers) {
			//	utils.Info(fmt.Sprintf("Could not find subset rules which matches provided headers '%s'. Creating a new one for '%s'.", ir.Headers, dr.Name), ips.TrackingId)
			//	newSubset := &v1alpha3.Subset{
			//		Name:   fmt.Sprintf("%s-%v-%s", ips.Name, ips.Build, ips.Namespace),
			//		Labels: ir.Selector.Labels,
			//	}
			//	_, err := appendSubset(ips.TrackingId, &dr, newSubset)
			//	if err != nil {
			//		utils.Fatal("", ips.TrackingId)
			//	}
			//
			//	err = UpdateDestinationRule(ips, updatedDr)
			//	if err != nil {
			//		utils.Fatal(fmt.Sprintf("%s", err), ips.TrackingId)
			//	}
			//}
		}
	}

	if ir.Weight > 0 {
		err = Percentage(ips, istioResources, ir)
		if err != nil {
			utils.Fatal(fmt.Sprintf("Could no create resource due to an error '%s'", err), ips.TrackingId)
		}
	}

	return nil
}

func appendSubset(trackingId string, dr *v1alpha32.DestinationRule, newSubset *v1alpha3.Subset) (*v1alpha32.DestinationRule, error) {
	for _, subsetValue := range dr.Spec.Subsets {
		if subsetValue.Name == newSubset.Name {
			// remove item from slice
			utils.Fatal(fmt.Sprintf("Found already existent subset '%s', refusing to update", subsetValue.Name), trackingId)
		}
	}

	dr.Spec.Subsets = append(dr.Spec.Subsets, newSubset)

	return dr, nil
}

// UpdateDestinationRule updates a specific virtualService given an updated object
func UpdateDestinationRule(ips *IstioOperator, destinationRule *v1alpha32.DestinationRule) error {
	utils.Info(fmt.Sprintf("Updating rule for destinationRule '%s'...", destinationRule.Name), ips.TrackingId)
	_, err := ips.Client.Istio.NetworkingV1alpha3().DestinationRules(ips.Namespace).Update(destinationRule)
	if err != nil {
		return err
	}
	return nil
}

// ClearRules will remove any destination & virtualService rules except the main one (provided by client).
// Ex: URI or Prefix
func (ips *IstioOperator) Clear(labels Selector) error {
	resources, err := GetResourcesToUpdate(ips, labels)
	if err != nil {
		return err
	}

	// Clean vs rules
	for _, vs := range resources.VirtualServiceList.Items {
		var cleanedRoutes []*v1alpha3.HTTPRoute
		for httpRuleKey, httpRuleValue := range vs.Spec.Http {
			for _, matchRuleValue := range httpRuleValue.Match {
				if matchRuleValue.Uri != nil {
					// remove rule with no Uri from HTTPRoute list to a posterior update
					cleanedRoutes = append(cleanedRoutes, vs.Spec.Http[httpRuleKey])
				}
			}
		}

		vs.Spec.Http = cleanedRoutes
		err := UpdateVirtualService(ips, &vs)
		if err != nil {
			utils.Fatal(fmt.Sprintf("Could not update virtualService '%s' due to error '%s'", vs.Name, err), ips.TrackingId)
		}
	}

	// Clean dr rules ?

	return nil
}

// GetResourcesToUpdate returns a slice of all DestinationRules and/or VirtualServices (based on given labelSelectors to a posterior update
func GetResourcesToUpdate(ips *IstioOperator, labelSelector Selector) (*IstioRouteList, error) {
	StringifyLabelSelector, _ := utils.StringifyLabelSelector(ips.TrackingId, labelSelector.Labels)

	listOptions := metav1.ListOptions{
		LabelSelector: StringifyLabelSelector,
	}

	matchedDrs, err := GetAllDestinationRules(ips, listOptions)
	if err != nil {
		utils.Fatal(fmt.Sprintf("%s", err), ips.TrackingId)
		return nil, err
	}

	matchedVss, err := GetAllVirtualServices(ips, listOptions)
	if err != nil {
		utils.Fatal(fmt.Sprintf("%s", err), ips.TrackingId)
		return nil, err
	}

	if len(matchedDrs.Items) == 0 || len(matchedVss.Items) == 0 {
		utils.Fatal(fmt.Sprintf("Couldn't find any istio resources based on given labelSelector '%s' to update. ", StringifyLabelSelector), ips.TrackingId)
		return nil, err
	}

	matchedResourcesList := &IstioRouteList{
		matchedVss,
		matchedDrs,
	}

	return matchedResourcesList, nil
}

// GetAllVirtualServices returns all istio resources 'virtualservices'
func GetAllVirtualServices(ips *IstioOperator, listOptions metav1.ListOptions) (virtualServiceList *v1alpha32.VirtualServiceList, error error) {
	utils.Info(fmt.Sprintf("Finding virtualServices which matches selector '%s'...", listOptions.LabelSelector), ips.TrackingId)
	vss, err := ips.Client.Istio.NetworkingV1alpha3().VirtualServices(ips.Namespace).List(listOptions)
	if err != nil {
		return nil, err
	}

	utils.Info(fmt.Sprintf("Found a total of '%d' virtualServices", len(vss.Items)), ips.TrackingId)
	return vss, nil
}

// GetAllVirtualservices returns all istio resources 'virtualservices'
func GetAllDestinationRules(ips *IstioOperator, listOptions metav1.ListOptions) (destinationRuleList *v1alpha32.DestinationRuleList, error error) {
	utils.Info(fmt.Sprintf("Finding destinationRules which matches selector '%s'...", listOptions.LabelSelector), ips.TrackingId)
	drs, err := ips.Client.Istio.NetworkingV1alpha3().DestinationRules(ips.Namespace).List(listOptions)
	if err != nil {
		return nil, err
	}

	utils.Info(fmt.Sprintf("Found a total of '%d' destinationRules", len(drs.Items)), ips.TrackingId)
	return drs, nil
}

// UpdateVirtualService updates a specific virtualService given an updated object
func UpdateVirtualService(ips *IstioOperator, virtualService *v1alpha32.VirtualService) error {
	utils.Info(fmt.Sprintf("Updating rule for virtualService '%s'...", virtualService.Name), ips.TrackingId)
	_, err := ips.Client.Istio.NetworkingV1alpha3().VirtualServices(ips.Namespace).Update(virtualService)
	if err != nil {
		return err
	}

	utils.Info("VirtualService successfully updated", ips.TrackingId)
	return nil
}

func Percentage(ips *IstioOperator, istioResources *IstioRouteList, ir *IstioRoute) (err error) {

	var matchedHeaders int
	var matchedUriSubset []string
	for _, vs := range istioResources.VirtualServiceList.Items {
		matchedUriSubset = []string{}
		for _, httpValue := range vs.Spec.Http {
			for matchKey, matchValue := range httpValue.Match {
				// Find a URI match to serve as final routing
				if matchValue.Uri != nil {
					matchedUriSubset = append(matchedUriSubset, httpValue.Route[matchKey].Destination.Subset)
					httpValue.Route[matchKey].Weight = 100 - ir.Weight
					fmt.Println(httpValue.Route[matchKey].Destination.Subset)
				}

				// Find the correct match-rule among all headers based on given input (ir.Weight.Headers)
				matchedHeaders = 0
				for headerKey, headerValue := range ir.Headers {
					if _, ok := matchValue.Headers[headerKey]; ok {
						if matchValue.Headers[headerKey].GetExact() == headerValue {
							matchedHeaders += 1
						}
					}
				}

				// In case of a Rule matches all headers' input, set weight between URI & Headers
				if matchedHeaders == len(ir.Headers) {
					utils.Info(fmt.Sprintf("Configuring weight to '%v' from '%s' in subset '%s'",
						ir.Weight, vs.Name, httpValue.Route[matchKey].Destination.Subset,
					), ips.TrackingId)
					httpValue.Route[matchKey].Weight = ir.Weight
				}
			}
		}

		if len(matchedUriSubset) == 0 {
			utils.Fatal(fmt.Sprintf("Could not find any URI match in '%s' for final routing.", vs.Name), ips.TrackingId)
		}

		if len(matchedUriSubset) > 1 {
			utils.Fatal(fmt.Sprintf("Found more than one URI match in '%s'. A unique URI match is expected instead.", vs.Name), ips.TrackingId)
		}

		err := UpdateVirtualService(ips, &vs)
		if err != nil {
			utils.Fatal(fmt.Sprintf("Could not update virtualService '%s' due to an error '%s'", vs.Name, err), ips.TrackingId)
		}

	}

	return nil

}
