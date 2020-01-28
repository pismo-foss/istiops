package router

import (
	"fmt"
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/pismo/istiops/pkg/logger"
	"github.com/pkg/errors"
	"istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
)

type VirtualService struct {
	TrackingId string
	Name       string
	Namespace  string
	Build      uint32
	Istio      IstioClientInterface
	KubeClient KubeClientInterface
}

// Clear will remove any virtualService's routes which are not master ones given a k8s labelSelector
func (v *VirtualService) Clear(s Shift, m string) error {
	//subsetName := fmt.Sprintf("%s-%v-%s", v.Name, v.Build, v.Namespace)

	//dr := DestinationRule{
	//	TrackingId: v.TrackingId,
	//	Name:       v.Name,
	//	Namespace:  v.Namespace,
	//	Build:      v.Build,
	//	Istio:      v.Istio,
	//	KubeClient: v.KubeClient,
	//}
	//
	//dss, err := dr.List(s.Selector)
	//if err != nil {
	//	return err
	//}

	vss, err := v.List(s.Selector)
	if err != nil {
		return err
	}

	// generating a cleaned list of routes with only route-master (URI: .+) included
	for _, vs := range vss.VList.Items {
		var cleanedRules []*v1alpha3.HTTPRoute
		cleanedRules = []*v1alpha3.HTTPRoute{}

		if m == "hard" {
			logger.Info(fmt.Sprintf("removing all virtualservice '%s' rules except the master-route one (Regex: .+) due to '%s' flag", vs.Name, m), v.TrackingId)
			for httpKey, httpValue := range vs.Spec.Http {
				for _, matchValue := range httpValue.Match {
					anyPrefix, _ := regexp.MatchString(`.+`, matchValue.Uri.GetPrefix())
					if anyPrefix {
						cleanedRules = append(cleanedRules, vs.Spec.Http[httpKey])
					}
					if matchValue.Uri.GetRegex() == ".+" {
						cleanedRules = append(cleanedRules, vs.Spec.Http[httpKey])
					}
				}
			}
		}

		if m == "soft" {
			logger.Info(fmt.Sprintf("removing all virtualservice '%s' rules with no pods associated due to '%s' flag", vs.Name, m), v.TrackingId)
			for httpKey, httpValue := range vs.Spec.Http {
				// append canary rules without pods associated - based on destinationRules
				//for _, routeValue := range httpValue.Route {
				//	for _, d := range dss.DList.Items {
				//		for _, subset := range d.Spec.Subsets {
				//
				//
				//		}
				//	}
				//}

				// append the default rule as well
				for _, matchValue := range httpValue.Match {
					anyPrefix, _ := regexp.MatchString(`.+`, matchValue.Uri.GetPrefix())
					if anyPrefix {
						cleanedRules = append(cleanedRules, vs.Spec.Http[httpKey])
					}
					if matchValue.Uri.GetRegex() == ".+" {
						cleanedRules = append(cleanedRules, vs.Spec.Http[httpKey])
					}
				}
			}
		}

		if m != "hard" && m != "soft" {
			return errors.New("empty mode when trying do clear routes. Refusing to continue")
		}

		if len(cleanedRules) == 0 {
			return errors.New("empty routes when cleaning virtualService's rules")
		}

		vs.Spec.Http = cleanedRules

		err := UpdateVirtualService(v, &vs)
		if err != nil {
			return err
		}
	}

	return nil
}

// Create returns a new route to be posterior appended to virtualService
func (v *VirtualService) Create(s Shift) (*IstioRules, error) {
	subsetName := fmt.Sprintf("%s-%v-%s", v.Name, v.Build, v.Namespace)

	logger.Info(fmt.Sprintf("Creating new http route for subset '%s'...", subsetName), v.TrackingId)
	newMatch := &v1alpha3.HTTPMatchRequest{}

	// append user labels to exact match
	if s.Traffic.Regexp {
		newMatch = &v1alpha3.HTTPMatchRequest{
			Headers: map[string]*v1alpha3.StringMatch{},
		}

		for headerKey, headerValue := range s.Traffic.RequestHeaders {
			newMatch.Headers[headerKey] = &v1alpha3.StringMatch{
				MatchType: &v1alpha3.StringMatch_Regex{
					Regex: headerValue,
				},
			}
		}
	}

	if s.Traffic.Exact {
		newMatch = &v1alpha3.HTTPMatchRequest{
			Headers: map[string]*v1alpha3.StringMatch{},
		}

		for headerKey, headerValue := range s.Traffic.RequestHeaders {
			newMatch.Headers[headerKey] = &v1alpha3.StringMatch{
				MatchType: &v1alpha3.StringMatch_Exact{
					Exact: headerValue,
				},
			}
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

	if len(s.Traffic.RequestHeaders) == 0 {
		return &IstioRules{}, errors.New("can't create a new route without request header's match")
	}

	logger.Info(fmt.Sprintf("Setting request header's match rule '%#v' for '%s'...", s.Traffic.RequestHeaders, subsetName), v.TrackingId)
	newRoute.Match = append(newRoute.Match, newMatch)
	newRoute.Route = append(newRoute.Route, defaultDestination)

	ir := IstioRules{
		MatchDestination: newRoute,
	}

	return &ir, nil
}

// Validate checks if VirtualService and Shift objects are correctly filled up
func (v *VirtualService) Validate(s Shift) error {
	if s.Traffic.Weight != 0 && len(s.Traffic.RequestHeaders) > 0 {
		return errors.New("a route needs to be served with a 'weight' or 'request headers', not both")
	}

	if s.Traffic.Weight == 0 && len(s.Traffic.RequestHeaders) == 0 {
		return errors.New("could not update route without 'weight' or 'headers'")
	}

	return nil

}

/* Update a virtualService with an existent route based on Shift object
or just create a new one (based on Create() method). The traffic will be balanced
based on Shift object with the inclusion of Weight or RequestHeaders attributes
*/
func (v *VirtualService) Update(s Shift) error {
	subsetName := fmt.Sprintf("%s-%v-%s", v.Name, v.Build, v.Namespace)

	vss, err := v.List(s.Selector)
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
				}
			}
		}

		if !routeExists {
			// create new route
			newHttpRoute, err := v.Create(s)
			if err != nil {
				return err
			}

			// ensure that http headers match will be the first element of vs.Spec.Http due to istio's rules precedence
			var auxHttp []*v1alpha3.HTTPRoute
			auxHttp = []*v1alpha3.HTTPRoute{}
			auxHttp = append(auxHttp, newHttpRoute.MatchDestination)
			for _, httpValue := range vs.Spec.Http {
				auxHttp = append(auxHttp, httpValue)
			}

			vs.Spec.Http = auxHttp
		}

		if routeExists {
			logger.Info("Found existent rule created for virtualService, skipping creation", v.TrackingId)

			// If a canary rule already exists, just warn it
			if len(s.Traffic.RequestHeaders) > 0 {
				logger.Warn(fmt.Sprintf("Already existent canary rule for build '%v', refusing to update it", v.Build), v.TrackingId)
			}

			// If a weight rule already exists, just update it
			if s.Traffic.Weight > 0 {
				httpRoutes, err := Percentage(v.TrackingId, subsetName, vs.Spec.Http, s)
				if err != nil {
					return err
				}

				httpRoutesNoHeaders, err := RemoveOutdatedRoutes(v.TrackingId, subsetName, httpRoutes)
				if err != nil {
					return err
				}

				vs.Spec.Http = httpRoutesNoHeaders
			}

		}

		err := UpdateVirtualService(v, &vs)
		if err != nil {
			return err
		}
	}

	return nil

}

// List will return all virtualServices which matches a k8s labelSelector
func (v *VirtualService) List(selector map[string]string) (*IstioRouteList, error) {
	logger.Debug(fmt.Sprintf("Getting virtualServices which matches label-selector '%s'", selector), v.TrackingId)
	stringified, err := Stringify(v.TrackingId, selector)
	if err != nil {
		return &IstioRouteList{}, err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: stringified,
	}

	vss, err := v.Istio.NetworkingV1alpha3().VirtualServices(v.Namespace).List(listOptions)
	if err != nil {
		return nil, err
	}

	if vss == nil {
		return nil, errors.New("empty virtualService list")
	}
	if len(vss.Items) <= 0 {
		return nil, errors.New(fmt.Sprintf("could not find any virtualServices which matched label-selector '%v'", listOptions.LabelSelector))
	}

	irl := &IstioRouteList{
		VList: vss,
	}

	return irl, nil
}

// UpdateDestinationRule updates a specific virtualService given an updated object
func UpdateVirtualService(vs *VirtualService, virtualService *v1alpha32.VirtualService) error {
	logger.Info(fmt.Sprintf("Updating route for virtualService '%s'...", virtualService.Name), vs.TrackingId)
	_, err := vs.Istio.NetworkingV1alpha3().VirtualServices(vs.Namespace).Update(virtualService)
	if err != nil {
		return err
	}
	return nil
}

// Balance returns a RouteDestination with balanced weight to be posterior appended to a virtualService
func Balance(currentSubset string, newSubset string, s Shift) ([]*v1alpha3.HTTPRouteDestination, error) {
	var routeBalanced []*v1alpha3.HTTPRouteDestination

	routeBalanced = []*v1alpha3.HTTPRouteDestination{}

	// if weight must be balanced between two subsets
	if s.Traffic.Weight < 100 {
		currentWeight := 100 - s.Traffic.Weight

		currentDestination := &v1alpha3.HTTPRouteDestination{
			Weight: currentWeight,
			Destination: &v1alpha3.Destination{
				Host:   s.Hostname,
				Subset: currentSubset,
				Port: &v1alpha3.PortSelector{
					Port: &v1alpha3.PortSelector_Number{
						Number: s.Port,
					},
				},
			},
		}

		routeBalanced = append(routeBalanced, currentDestination)
	}

	newDestination := &v1alpha3.HTTPRouteDestination{
		Weight: s.Traffic.Weight,
		Destination: &v1alpha3.Destination{
			Host:   s.Hostname,
			Subset: newSubset,
			Port: &v1alpha3.PortSelector{
				Port: &v1alpha3.PortSelector_Number{
					Number: s.Port,
				},
			},
		},
	}

	routeBalanced = append(routeBalanced, newDestination)

	return routeBalanced, nil
}

// Remove returns a slice of Routes without an element given an index
func Remove(slice []*v1alpha3.HTTPRoute, index int) []*v1alpha3.HTTPRoute {
	return append(slice[:index], slice[index+1:]...)
}

// RemoveOutdatedRoutes returns a slice without any route which matches the given subset
func RemoveOutdatedRoutes(trackingId string, subset string, httpRoute []*v1alpha3.HTTPRoute) ([]*v1alpha3.HTTPRoute, error) {
	var noMasterRoutes []*v1alpha3.HTTPRoute
	var masterRoute *v1alpha3.HTTPRoute
	var cleanedRoutes []*v1alpha3.HTTPRoute

	// get a HTTPRoute without a master route to be posterior cleaned
	for httpKey, httpValue := range httpRoute {
		// remove request header route based on subset to avoid non-used rules persisted
		for _, matchValue := range httpValue.Match {
			if matchValue.Uri.GetRegex() == ".+" {
				masterRoute = httpRoute[httpKey]
				noMasterRoutes = Remove(httpRoute, httpKey)
			}
		}
	}

	cleanedRoutes = noMasterRoutes

	// for the cleaned HTTPRoute, every route with the given subset must be removed
	for httpKey, httpValue := range noMasterRoutes {
		for _, routeValue := range httpValue.Route {
			if routeValue.Destination.Subset == subset {
				logger.Info(fmt.Sprintf("removing outdated rule for '%s' subset in order to weight routing", subset), trackingId)
				cleanedRoutes = Remove(noMasterRoutes, httpKey)
			}
		}
	}

	cleanedRoutes = append(cleanedRoutes, masterRoute)

	if len(cleanedRoutes) == 0 {
		return nil, errors.New("empty routes when removing outdated subsets")
	}

	if cleanedRoutes[len(cleanedRoutes)-1] == nil {
		return nil, errors.New("got nil routes when removing outdated subsets")
	}

	if cleanedRoutes[len(cleanedRoutes)-1].Match[0].Uri.GetRegex() != ".+" {
		return nil, errors.New(fmt.Sprintf("non master-route as last element in routes: %s", cleanedRoutes))
	}

	return cleanedRoutes, nil
}

// Percentage returns a []HTTPRoute with weight routing set to be posterior appended to a virtualService
func Percentage(trackingId string, subset string, httpRoute []*v1alpha3.HTTPRoute, s Shift) ([]*v1alpha3.HTTPRoute, error) {
	// Finding master route (URI match)
	var masterRouteCounter int
	var masterIndex int

	// work with the need of cleaning old headers for the same subset

	if len(httpRoute) == 0 {
		return nil, errors.New("empty routes")
	}

	// work with percentage rules
	for httpKey, httpValue := range httpRoute {
		for _, matchValue := range httpValue.Match {

			// reconstruct master route to attend a balanced traffic between versions
			if matchValue.Uri.GetRegex() == ".+" {
				logger.Info(fmt.Sprintf("Updating master route to balance canary traffic"), trackingId)

				masterRouteCounter += 1
				masterIndex = httpKey

				newSubset := httpValue.Route[0].Destination.Subset

				balancedRoute, err := Balance(newSubset, subset, s)
				if err != nil {
					return nil, err
				}

				httpRoute[httpKey].Route = balancedRoute

				if len(httpRoute[httpKey].Route) > 2 {
					return nil, errors.New("more than 2 destination for route")
				}
			}
		}
	}

	// setting URI Master route to the last element of []*Routes due to istio's traffic rule precedence
	tempMasterRoute := httpRoute[masterIndex]
	httpRoute = Remove(httpRoute, masterIndex)
	httpRoute = append(httpRoute, tempMasterRoute)

	if masterRouteCounter > 1 {
		return nil, errors.New("multiple master routes found")
	}

	// create a master route rule if does not exists
	if masterRouteCounter == 0 {
		logger.Info(fmt.Sprintf("Could not find a master route 'Regex: .+', creating with 100%% of weight..."), trackingId)
		routeMaster := &v1alpha3.HTTPRoute{}
		routeMasterMatch := &v1alpha3.HTTPMatchRequest{Uri: &v1alpha3.StringMatch{MatchType: &v1alpha3.StringMatch_Regex{Regex: ".+"}}}

		routeMasterDestination := &v1alpha3.HTTPRouteDestination{
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

		routeMaster.Match = append(routeMaster.Match, routeMasterMatch)
		routeMaster.Route = append(routeMaster.Route, routeMasterDestination)
		httpRoute = append(httpRoute, routeMaster)
	}

	return httpRoute, nil
}

// ValidateVirtualServiceList checks for inconsistencies in IstioRouteList.VList
func ValidateVirtualServiceList(irl *IstioRouteList) error {
	if irl.VList == nil {
		return errors.New("empty virtualServices list")
	}

	if len(irl.VList.Items) == 0 {
		return errors.New("empty virtualServices")
	}

	return nil
}
