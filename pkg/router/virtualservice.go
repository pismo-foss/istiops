package router

import (
	"fmt"

	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"github.com/pismo/istiops/utils"
	"github.com/pkg/errors"
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
	// remover metadata de struct
	Metadata VsMetadata
	Istio    *versioned.Clientset
}

func (v *VirtualService) Validate(s *Shift) error {
	// validar cada problema individualmente pra gerar erros especificos
	if s.Traffic.Weight != 0 && len(s.Traffic.RequestHeaders) >= 0 {
		return errors.New("a route needs to be served with a 'weight' or 'request headers', not both")
	}

	return nil

}

func (v *VirtualService) Update(s *Shift) error {
	subsetName := fmt.Sprintf("%s-%v-%s", v.Metadata.Name, v.Metadata.Build, v.Metadata.Namespace)

	StringifyLabelSelector, err := utils.StringifyLabelSelector(v.Metadata.TrackingId, s.Selector.Labels)
	if err != nil {
		return err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: StringifyLabelSelector,
	}

	// nao mandar dados direto pro metodo, aqui nao faria sentido enviar so o client do istio? tem muita informacao sendo enviada e nao utilizada
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

		if !subsetExists {
			// create new subset
			newHttpRoute, err := CreateNewRoute(subsetName, v, s)
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

func (v *VirtualService) Clear(s *Shift) error {
	StringifyLabelSelector, err := utils.StringifyLabelSelector(v.Metadata.TrackingId, s.Selector.Labels)
	if err != nil {
		return err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: StringifyLabelSelector,
	}

	vss, err := GetAllVirtualServices(v, s, listOptions)
	if err != nil {
		return err
	}

	for _, vs := range vss.Items {
		for httpRuleKey, httpRuleValue := range vs.Spec.Http {
			for _, httpRoute := range httpRuleValue.Route {
				if httpRoute.Weight <= 0 {
					utils.Info(fmt.Sprintf("The subset '%s' will be removed due to a non-active weight rule attached", httpRoute.Destination.Subset), v.Metadata.TrackingId)
					vs.Spec.Http = append(vs.Spec.Http[:httpRuleKey], vs.Spec.Http[httpRuleKey+1:]...)
				}
			}
		}

		// In case of all rules had being removed, refuse to continue
		if len(vs.Spec.Http) == 0 {
			return errors.New(fmt.Sprintf("the clear command will result in a resource '%s' without any rules which is not accepted by istio", vs.Name))
		}

		utils.Info(fmt.Sprintf("Clearing all virtualService routes from '%s' except the URI or Weighted ones...", vs.Name), v.Metadata.TrackingId)
		err := UpdateVirtualService(v, &vs)
		if err != nil {
			return err
		}
	}

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
