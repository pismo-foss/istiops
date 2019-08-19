package operator

import (
	"fmt"
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/pismo/istiops/pkg/router"
	"github.com/pismo/istiops/utils"
)

type Istiops struct {
	Metadata              *router.Metadata
	VirtualServiceRouter  *router.VirtualServiceRoute
	DestinationRuleRouter *router.DestinationRuleRoute
}

//should be inside router for vs and dr
type IstioRouteList struct {
	VirtualServiceList   *v1alpha32.VirtualServiceList
	DestinationRulesList *v1alpha32.DestinationRuleList
}

func (ips *Istiops) Get(r *router.Route) error {
	return nil
}

func (ips *Istiops) Create(r *router.Route) error {

	VsRouter := ips.VirtualServiceRouter
	vs, err := VsRouter.Validate(r)
	fmt.Println(vs)
	if err != nil {
		return err
	}

	DrRouter := ips.DestinationRuleRouter
	dr, err := DrRouter.Validate(r)
	fmt.Println(dr)
	if err != nil {
		return err
	}

	return nil
}

func (ips *Istiops) Delete(r *router.Route) error {
	fmt.Println("Initializing something")
	fmt.Println("", ips.Metadata.TrackingId)
	return nil
}

func (ips *Istiops) Update(r *router.Route) error {
	if len(r.Selector.Labels) == 0 || len(r.Traffic.PodSelector) == 0 {
		utils.Fatal(fmt.Sprintf("Selectors must not be empty otherwise istiops won't be able to find any resources."), ips.Metadata.TrackingId)
	}

	VsRouter := ips.VirtualServiceRouter
	vs, err := VsRouter.Validate(r)
	fmt.Println(vs)
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	DrRouter := ips.DestinationRuleRouter
	dr, err := DrRouter.Validate(r)
	fmt.Println(dr)
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}
	err = DrRouter.Update(r)
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}
	err = VsRouter.Update(r)
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	if r.Traffic.Weight > 0 {
		// update router to serve percentage
		if err != nil {
			utils.Fatal(fmt.Sprintf("Could no create resource due to an error '%s'", err), ips.Metadata.TrackingId)
		}
	}

	return nil
}

// ClearRules will remove any destination & virtualService rules except the main one (provided by client).
// Ex: URI or Prefix
func (ips *Istiops) Clear(labels map[string]string) error {
	//resources, err := GetResourcesToUpdate(ips, labels)
	//if err != nil {
	//	return err
	//}

	// Clean vs rules
	//for _, vs := range resources.VirtualServiceList.Items {
	//	var cleanedRoutes []*v1alpha3.HTTPRoute
	//	for httpRuleKey, httpRuleValue := range vs.Spec.Http {
	//		for _, matchRuleValue := range httpRuleValue.Match {
	//			if matchRuleValue.Uri != nil {
	//				// remove rule with no Uri from HTTPRoute list to a posterior update
	//				cleanedRoutes = append(cleanedRoutes, vs.Spec.Http[httpRuleKey])
	//			}
	//		}
	//	}
	//
	//	vs.Spec.Http = cleanedRoutes
	//	// update virtualService
	//}

	// Clean dr rules ?

	return nil
}

// UpdateVirtualService updates a specific virtualService given an updated object
func UpdateVirtualService(drRoute router.DestinationRuleRoute, ips *Istiops, virtualService *v1alpha32.VirtualService) error {
	utils.Info(fmt.Sprintf("Updating rule for virtualService '%s'...", virtualService.Name), ips.Metadata.TrackingId)
	_, err := drRoute.Istio.NetworkingV1alpha3().VirtualServices(ips.Metadata.Namespace).Update(virtualService)
	if err != nil {
		return err
	}

	utils.Info("VirtualService successfully updated", ips.Metadata.TrackingId)
	return nil
}

func Percentage(ips *Istiops, istioResources *IstioRouteList, r *router.Route) (err error) {

	var matchedHeaders int
	var matchedUriSubset []string
	for _, vs := range istioResources.VirtualServiceList.Items {
		matchedUriSubset = []string{}
		for _, httpValue := range vs.Spec.Http {
			for matchKey, matchValue := range httpValue.Match {
				// Find a URI match to serve as final routing
				if matchValue.Uri != nil {
					matchedUriSubset = append(matchedUriSubset, httpValue.Route[matchKey].Destination.Subset)
					httpValue.Route[matchKey].Weight = 100 - r.Traffic.Weight
					fmt.Println(httpValue.Route[matchKey].Destination.Subset)
				}

				// Find the correct match-rule among all headers based on given input (ir.Weight.Headers)
				matchedHeaders = 0
				for headerKey, headerValue := range r.Traffic.RequestHeaders {
					if _, ok := matchValue.Headers[headerKey]; ok {
						if matchValue.Headers[headerKey].GetExact() == headerValue {
							matchedHeaders += 1
						}
					}
				}

				// In case of a Rule matches all headers' input, set weight between URI & Headers
				if matchedHeaders == len(r.Traffic.RequestHeaders) {
					utils.Info(fmt.Sprintf("Configuring weight to '%v' from '%s' in subset '%s'",
						r.Traffic.Weight, vs.Name, httpValue.Route[matchKey].Destination.Subset,
					), ips.Metadata.TrackingId)
					httpValue.Route[matchKey].Weight = r.Traffic.Weight
				}
			}
		}

		if len(matchedUriSubset) == 0 {
			utils.Fatal(fmt.Sprintf("Could not find any URI match in '%s' for final routing.", vs.Name), ips.Metadata.TrackingId)
		}

		if len(matchedUriSubset) > 1 {
			utils.Fatal(fmt.Sprintf("Found more than one URI match in '%s'. A unique URI match is expected instead.", vs.Name), ips.Metadata.TrackingId)
		}

		//err := UpdateVirtualService(ips, &vs)
		//if err != nil {
		//	utils.Fatal(fmt.Sprintf("Could not update virtualService '%s' due to an error '%s'", vs.Name, err), ips.Metadata.TrackingId)
		//}

	}

	return nil

}
