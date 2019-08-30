package router

import (
	"fmt"
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/pismo/istiops/pkg/logger"
	"github.com/pkg/errors"
	"istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VirtualService struct {
	TrackingId string
	Name       string
	Namespace  string
	Build      uint32
	Istio      Client
}

func (v *VirtualService) Clear(s *Shift) error {
	return nil
}

func (v *VirtualService) Create(s *Shift) (*IstioRules, error) {
	subsetName := fmt.Sprintf("%s-%v-%s", v.Name, v.Build, v.Namespace)

	logger.Info(fmt.Sprintf("Creating new http route for subset '%s'...", subsetName), v.TrackingId)
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
		logger.Info(fmt.Sprintf("Setting request header's match rule '%s' for '%s'...", s.Traffic.RequestHeaders, subsetName), v.TrackingId)
		newRoute.Match = append(newRoute.Match, newMatch)
	}

	newRoute.Route = append(newRoute.Route, defaultDestination)

	ir := IstioRules{
		MatchDestination: newRoute,
	}

	return &ir, nil
}

func (v *VirtualService) Validate(s *Shift) error {
	if s.Traffic.Weight != 0 && len(s.Traffic.RequestHeaders) > 0 {
		return errors.New("a route needs to be served with a 'weight' or 'request headers', not both")
	}

	return nil

}

func (v *VirtualService) Update(s *Shift) error {
	subsetName := fmt.Sprintf("%s-%v-%s", v.Name, v.Build, v.Namespace)

	StringifyLabelSelector, err := StringifyLabelSelector(v.TrackingId, s.Selector.Labels)
	if err != nil {
		return err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: StringifyLabelSelector,
	}

	vss, err := v.List(listOptions)
	if err != nil {
		return err
	}

	for _, vs := range vss.VList.Items {
		routeExists := false
		for _, httpValue := range vs.Spec.Http {
			for _, routeValue := range httpValue.Route {
				// if subset already exists
				if routeValue.Destination.Subset == subsetName {
					routeExists = true
					logger.Info("Found existent rule created for virtualService, skipping creation", v.TrackingId)
				}
			}
		}

		if !routeExists {
			// create new route
			newHttpRoute, err := v.Create(s)
			if err != nil {
				return err
			}

			vs.Spec.Http = append(vs.Spec.Http, newHttpRoute.MatchDestination)
		}

		if routeExists {
			if s.Traffic.Weight > 0 && s.Traffic.Weight < 100 {
				fmt.Println("It's time to balance traffic")
				err := percentage(vs, subsetName, s)
				if err != nil {
					return err
				}
			}

			// set default rule to 100%
			if s.Traffic.Weight == 100 {
				fmt.Println("it's time to set 100% traffic")
			}
		}

		err := UpdateVirtualService(v, &vs)
		if err != nil {
			return err
		}
	}

	return nil

}

func percentage(vs v1alpha32.VirtualService, subset string, s *Shift) error {
	// Finding master route (URI match)
	masterRouteCounter := 0
	for httpKey, httpValue := range vs.Spec.Http {
		for _, matchValue := range httpValue.Match {

			// reconstruct master route to attend a balanced traffic between versions
			if matchValue.Uri.GetRegex() == ".+" {
				masterRouteCounter += 1

				fmt.Println(matchValue.Uri.GetRegex())
				fmt.Println(vs.Spec.Http[httpKey])
				fmt.Println(vs.Spec.Http[httpKey].Route)

				currentWeight := s.Traffic.Weight - 100
				fmt.Println(s.Traffic.Weight)

				currentDestination := &v1alpha3.HTTPRouteDestination{
					Weight: currentWeight,
					Destination: &v1alpha3.Destination{
						Host:   s.Hostname,
						Subset: subset,
						Port: &v1alpha3.PortSelector{
							Port: &v1alpha3.PortSelector_Number{
								Number: s.Port,
							},
						},
					},
				}

				newDestination := &v1alpha3.HTTPRouteDestination{
					Weight: s.Traffic.Weight,
					Destination: &v1alpha3.Destination{
						Host:   s.Hostname,
						Subset: subset,
						Port: &v1alpha3.PortSelector{
							Port: &v1alpha3.PortSelector_Number{
								Number: s.Port,
							},
						},
					},
				}

				vs.Spec.Http[httpKey].Route = []*v1alpha3.HTTPRouteDestination{}
				vs.Spec.Http[httpKey].Route = append(vs.Spec.Http[httpKey].Route, currentDestination)
				vs.Spec.Http[httpKey].Route = append(vs.Spec.Http[httpKey].Route, newDestination)
				fmt.Println(vs.Spec.Http[httpKey].Route)

				if len(vs.Spec.Http[httpKey].Route) != 2 {
					return errors.New("more than 2 destination for route")
				}
			}
		}

	}

	if masterRouteCounter != 1 {
		return errors.New("multiple master routes (URI: .+)")
	}

	// create a master route rule if
	if masterRouteCounter == 0 {

	}

	return nil
}

func master(vs v1alpha32.VirtualService, s *Shift) error {
	return nil
}

func (v *VirtualService) List(opts metav1.ListOptions) (*IstioRouteList, error) {
	vss, err := v.Istio.Versioned.NetworkingV1alpha3().VirtualServices(v.Namespace).List(opts)
	if err != nil {
		return nil, err
	}

	if len(vss.Items) <= 0 {
		return nil, errors.New(fmt.Sprintf("could not find any virtualServices which matched label-selector '%v'", opts.LabelSelector))
	}

	irl := &IstioRouteList{
		VList: vss,
	}

	return irl, nil
}

// UpdateDestinationRule updates a specific virtualService given an updated object
func UpdateVirtualService(vs *VirtualService, virtualService *v1alpha32.VirtualService) error {
	logger.Info(fmt.Sprintf("Updating route for virtualService '%s'...", virtualService.Name), vs.TrackingId)
	_, err := vs.Istio.Versioned.NetworkingV1alpha3().VirtualServices(vs.Namespace).Update(virtualService)
	if err != nil {
		return err
	}
	return nil
}
