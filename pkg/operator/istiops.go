package operator

import (
	"fmt"
	"github.com/pismo/istiops/pkg/router"
	"github.com/pismo/istiops/utils"
)

type Istiops struct {
	Shift    *router.Shift
	VsRouter *router.VirtualService
	DrRouter *router.DestinationRule
}

func (ips *Istiops) Get(r *router.Shift) error {
	return nil
}

func (ips *Istiops) Create(r *router.Shift) error {

	VsRouter := ips.VsRouter
	err := VsRouter.Validate(r)
	if err != nil {
		return err
	}

	DrRouter := ips.DrRouter
	err = DrRouter.Validate(r)
	if err != nil {
		return err
	}

	return nil
}

func (ips *Istiops) Delete(r *router.Shift) error {
	fmt.Println("Initializing something")
	fmt.Println("", "cid")
	return nil
}

func (ips *Istiops) Update(r *router.Shift) error {
	if len(r.Selector.Labels) == 0 || len(r.Traffic.PodSelector) == 0 {
		utils.Fatal(fmt.Sprintf("Selectors must not be empty otherwise istiops won't be able to find any resources."), "")
	}

	DrRouter := ips.DrRouter
	VsRouter := ips.VsRouter
	var err error

	err = DrRouter.Validate(r)
	if err != nil {
		utils.Fatal(fmt.Sprintf("%s", err), ips.DrRouter.Metadata.TrackingId)
	}
	err = DrRouter.Update(r)
	if err != nil {
		utils.Fatal(fmt.Sprintf("%s", err), ips.DrRouter.Metadata.TrackingId)
	}

	err = VsRouter.Validate(r)
	if err != nil {
		utils.Fatal(fmt.Sprintf("%s", err), ips.VsRouter.Metadata.TrackingId)
	}
	err = VsRouter.Update(r)
	if err != nil {
		utils.Fatal(fmt.Sprintf("%s", err), ips.VsRouter.Metadata.TrackingId)
	}

	if r.Traffic.Weight > 0 {
		// update router to serve percentage
		if err != nil {
			utils.Fatal(fmt.Sprintf("Could no create resource due to an error '%s'", err), "")
		}
	}

	return nil
}

// ClearRules will remove any destination & virtualService rules except the main one (provided by client).
// Ex: URI or Prefix
func (ips *Istiops) Clear(s *router.Shift) error {
	DrRouter := ips.DrRouter
	VsRouter := ips.VsRouter
	var err error

	err = DrRouter.Clear(ips.Shift)
	if err != nil {
		utils.Fatal(fmt.Sprintf("%s", err), ips.DrRouter.Metadata.TrackingId)
	}

	err = VsRouter.Clear(ips.Shift)
	if err != nil {
		utils.Fatal(fmt.Sprintf("%s", err), ips.DrRouter.Metadata.TrackingId)
	}

	// Clean dr rules ?

	return nil
}

func Percentage(istioResources *router.IstioRouteList, s *router.Shift) (err error) {

	var matchedHeaders int
	var matchedUriSubset []string
	for _, vs := range istioResources.VirtualServiceList.Items {
		matchedUriSubset = []string{}
		for _, httpValue := range vs.Spec.Http {
			for matchKey, matchValue := range httpValue.Match {
				// Find a URI match to serve as final routing
				if matchValue.Uri != nil {
					matchedUriSubset = append(matchedUriSubset, httpValue.Route[matchKey].Destination.Subset)
					httpValue.Route[matchKey].Weight = 100 - s.Traffic.Weight
					fmt.Println(httpValue.Route[matchKey].Destination.Subset)
				}

				// Find the correct match-rule among all headers based on given input (ir.Weight.Headers)
				matchedHeaders = 0
				for headerKey, headerValue := range s.Traffic.RequestHeaders {
					if _, ok := matchValue.Headers[headerKey]; ok {
						if matchValue.Headers[headerKey].GetExact() == headerValue {
							matchedHeaders += 1
						}
					}
				}

				// In case of a Rule matches all headers' input, set weight between URI & Headers
				if matchedHeaders == len(s.Traffic.RequestHeaders) {
					utils.Info(fmt.Sprintf("Configuring weight to '%v' from '%s' in subset '%s'",
						s.Traffic.Weight, vs.Name, httpValue.Route[matchKey].Destination.Subset,
					), "cid")
					httpValue.Route[matchKey].Weight = s.Traffic.Weight
				}
			}
		}

		if len(matchedUriSubset) == 0 {
			utils.Fatal(fmt.Sprintf("Could not find any URI match in '%s' for final routing.", vs.Name), "cid")
		}

		if len(matchedUriSubset) > 1 {
			utils.Fatal(fmt.Sprintf("Found more than one URI match in '%s'. A unique URI match is expected instead.", vs.Name), "")
		}

		//err := UpdateVirtualService(ips, &vs)
		//if err != nil {
		//	utils.Fatal(fmt.Sprintf("Could not update virtualService '%s' due to an error '%s'", vs.Name, err), ips.Metadata.TrackingId)
		//}

	}

	return nil

}
