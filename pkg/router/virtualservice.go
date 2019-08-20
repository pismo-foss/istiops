package router

import (
	"fmt"
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"github.com/pismo/istiops/utils"
	"istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VsMetadata struct {
	TrackingId string
	Name       string
	Namespace  string
	Build      uint32
}

type VirtualService struct {
	Metadata VsMetadata
	Istio    *versioned.Clientset
}

func (v *VirtualService) Validate(s *Shift) (v1alpha3.VirtualService, error) {
	fmt.Println("validating vr")
	return v1alpha3.VirtualService{}, nil

}

func (v *VirtualService) Update(s *Shift) error {
	fmt.Println("updating virtualservice")
	subsetName := fmt.Sprintf("%s-%v-%s", v.Metadata.Name, v.Metadata.Build, v.Metadata.Namespace)

	StringifyLabelSelector, err := utils.StringifyLabelSelector(v.Metadata.TrackingId, s.Selector.Labels)
	if err != nil {
		fmt.Println("null drs")
		return err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: StringifyLabelSelector,
	}

	vss, err := GetAllVirtualServices(v, s, listOptions)
	for _, vs := range vss.Items {
		subsetExists := false
		for _, httpValue := range vs.Spec.Http {
			for _, routeValue := range httpValue.Route {
				// if subset already exists
				if routeValue.Destination.Host == subsetName {
					subsetExists = true
					return err
				}
			}
		}

		if ! subsetExists {
			// create new subset
			utils.Info(fmt.Sprintf("Creating new route"), v.Metadata.TrackingId)
			newHttpRoute, err := CreateNewRoute(subsetName, v, s)
			fmt.Println(newHttpRoute)
			if err != nil {
				return err
			}

			vs.Spec.Http = append(vs.Spec.Http, newHttpRoute)
		}

		err := UpdateVirtualService(v, &vs)
		if err != nil {
			return err
		}

	}

	return nil

}

func (v *VirtualService) Delete(s *Shift) error {
	return nil

}

// GetAllVirtualServices returns all istio resources 'virtualservices'
func GetAllVirtualServices(vsRoute *VirtualService, s *Shift, listOptions metav1.ListOptions) (*v1alpha32.VirtualServiceList, error) {
	utils.Info(fmt.Sprintf("Getting all virtualservices..."), vsRoute.Metadata.TrackingId)

	vss, err := vsRoute.Istio.NetworkingV1alpha3().VirtualServices(vsRoute.Metadata.Namespace).List(listOptions)
	if err != nil {
		return nil, err
	}

	return vss, nil
}

func CreateNewRoute(subsetName string, vsRoute *VirtualService, s *Shift) (*v1alpha3.HTTPRoute, error) {
	utils.Info(fmt.Sprintf("Creating new http route for subset '%s'...", subsetName), vsRoute.Metadata.TrackingId)
	newMatch := &v1alpha3.HTTPMatchRequest{
		Headers: map[string]*v1alpha3.StringMatch{},
	}

	// append user labels to exact match
	for headerKey, headerValue := range s.Traffic.RequestHeaders {
		newMatch.Headers[headerKey] = &v1alpha3.StringMatch{
			MatchType: &v1alpha3.StringMatch_Exact{
				Exact: headerValue,
			},
		}
	}

	defaultDestination := &v1alpha3.HTTPRouteDestination{
		Destination: &v1alpha3.Destination{
			Host:   s.Hostname,
			Subset: subsetName,
			Port: &v1alpha3.PortSelector{
				Port: &v1alpha3.PortSelector_Number{
					Number: s.Port,
				},
			},
		},
	}

	newRoute := &v1alpha3.HTTPRoute{}

	if len(s.Traffic.RequestHeaders) > 0 {
		utils.Info(fmt.Sprintf("Setting request header's match rule '%s' for '%s'...", s.Traffic.RequestHeaders, subsetName), vsRoute.Metadata.TrackingId)
		newRoute.Match = append(newRoute.Match, newMatch)
	}

	newRoute.Route = append(newRoute.Route, defaultDestination)

	return newRoute, nil
}

// UpdateDestinationRule updates a specific virtualService given an updated object
func UpdateVirtualService(vs *VirtualService, virtualService *v1alpha32.VirtualService) error {
	utils.Info(fmt.Sprintf("Updating route for virtualService '%s'...", virtualService.Name), vs.Metadata.TrackingId)
	_, err := vs.Istio.NetworkingV1alpha3().VirtualServices(vs.Metadata.Namespace).Update(virtualService)
	if err != nil {
		return err
	}
	return nil
}
