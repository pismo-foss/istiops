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
	Namespace  string
	Client     *client.ClientSet
}

type IstioResources struct {
	DestinationRuleName string
	VirtualServiceName  string
}

type IstioResourcesList struct {
	VirtualServiceList   *v1alpha32.VirtualServiceList
	DestinationRulesList *v1alpha32.DestinationRuleList
}

type Istiops interface {
	Create(ir *IstioResources)
	Delete(ir *IstioResources)
	Validate(cid string)
	Clear(LabelSelector) error
}

type LabelSelector struct {
	Labels map[string]string
}

func (ips *IstioOperator) Create(ir *IstioResources) {
	fmt.Println("creating")

	istioNetworking := ips.Client.Istio.NetworkingV1alpha3()
	dr := &v1alpha32.DestinationRule{}
	dr.Name = ir.DestinationRuleName

	_, err := istioNetworking.DestinationRules(ips.Namespace).Create(dr)

	if err != nil {
		utils.Fatal(fmt.Sprintf("Could no create resource due to an error '%s'", err), ips.TrackingId)
	}

	fmt.Println("Creating something")
}

func (ips *IstioOperator) Delete(ir *IstioResources) {
	fmt.Println("Initializing something")
	fmt.Println("", ips.TrackingId)
}

func (ips *IstioOperator) Validate(cid string) {
	fmt.Println("Initializing something")
	fmt.Println("", cid)
}

// ClearRules will remove any destination & virtualService rules except the main one (provided by client).
// Ex: URI or Prefix
func (ips *IstioOperator) Clear(labels LabelSelector) error {
	resources, err := GetResourcesToUpdate(ips, labels)
	if err != nil {
		return err
	}

	// Clean vs rules
	for _, vs := range resources.VirtualServiceList.Items {
		var cleanedRoutes []*v1alpha3.HTTPRoute
		fmt.Println(vs.Name)
		for httpRuleKey, httpRuleValue := range vs.Spec.Http {
			for _, matchRuleValue := range httpRuleValue.Match {
				if matchRuleValue.Uri != nil {
					// remove rule with no Uri from HTTPRoute list to a posterior update
					cleanedRoutes = append(cleanedRoutes, vs.Spec.Http[httpRuleKey])
				}
			}
		}

		vs.Spec.Http = cleanedRoutes
		fmt.Println("Update!")
		err := UpdateVirtualService(ips, &vs)
		if err != nil {
			utils.Fatal(fmt.Sprintf("Could not update virtualService '%s' due to error '%s'", vs.Name, err), ips.TrackingId)
		}
	}

	// Clean dr rules ?

	return nil
}

// GetResourcesToUpdate returns a slice of all DestinationRules and/or VirtualServices (based on given labelSelectors to a posterior update
func GetResourcesToUpdate(ips *IstioOperator, labelSelector LabelSelector) (matchedResourcesList *IstioResourcesList, error error) {
	stringfiedLabelSelector, _ := utils.StringfyLabelSelector(ips.TrackingId, labelSelector.Labels)

	listOptions := metav1.ListOptions{
		LabelSelector: stringfiedLabelSelector,
	}

	matchedDrs, err := GetAllDestinationRules(ips, listOptions)
	fmt.Println(err)
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
		utils.Fatal(fmt.Sprintf("Couldn't find any istio resources based on given labelSelector '%s' to update. ", stringfiedLabelSelector), ips.TrackingId)
		return nil, err
	}

	matchedResourcesList = &IstioResourcesList{
		matchedVss,
		matchedDrs,
	}

	return matchedResourcesList, nil
}

// GetAllVirtualServices returns all istio resources 'virtualservices'
func GetAllVirtualServices(ips *IstioOperator, listOptions metav1.ListOptions) (virtualServiceList *v1alpha32.VirtualServiceList, error error) {
	utils.Info(fmt.Sprintf("Getting all virtualservices..."), ips.TrackingId)
	vss, err := ips.Client.Istio.NetworkingV1alpha3().VirtualServices(ips.Namespace).List(listOptions)
	if err != nil {
		return nil, err
	}

	return vss, nil
}

// GetAllVirtualservices returns all istio resources 'virtualservices'
func GetAllDestinationRules(ips *IstioOperator, listOptions metav1.ListOptions) (destinationRuleList *v1alpha32.DestinationRuleList, error error) {
	utils.Info(fmt.Sprintf("Getting all destinationrules..."), ips.TrackingId)
	drs, err := ips.Client.Istio.NetworkingV1alpha3().DestinationRules(ips.Namespace).List(listOptions)
	if err != nil {
		return nil, err
	}

	return drs, nil
}

// UpdateVirtualService updates a specific virtualService given an updated object
func UpdateVirtualService(ips *IstioOperator, virtualService *v1alpha32.VirtualService) error {
	utils.Info(fmt.Sprintf("Updating rule for virtualService '%s'...", virtualService.Name), ips.TrackingId)
	_, err := ips.Client.Istio.NetworkingV1alpha3().VirtualServices(ips.Namespace).Update(virtualService)
	if err != nil {
		return err
	}
	return nil
}
