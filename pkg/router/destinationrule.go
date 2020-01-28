package router

import (
	"errors"
	"fmt"
	"github.com/pismo/istiops/pkg/logger"

	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DestinationRule struct {
	TrackingId string
	Name       string
	Namespace  string
	Build      uint32
	Istio      IstioClientInterface
	KubeClient KubeClientInterface
}

// Clear will remove any subset which are not used by a virtualService given a k8s labelSelector
func (d *DestinationRule) Clear(s Shift, m string) error {
	v := VirtualService{
		TrackingId: d.TrackingId,
		Name:       d.Name,
		Namespace:  d.Namespace,
		Build:      d.Build,
		Istio:      d.Istio,
	}

	vss, err := v.List(s.Selector)
	if err != nil {
		return err
	}

	drs, err := d.List(s.Selector)
	if err != nil {
		return err
	}

	var cleanedSubsetList []*v1alpha3.Subset

	for _, dr := range drs.DList.Items {
		// validate for each subset it's own existence in virtualServices
		for _, subset := range dr.Spec.Subsets {
			subsetExists := false
			for _, vs := range vss.VList.Items {
				for _, http := range vs.Spec.Http {
					for _, route := range http.Route {
						if subset.GetName() == route.Destination.Subset {
							subsetExists = true
						}
					}
				}
			}

			// create a new subsetList with only the active ones
			if subsetExists {
				logger.Info(fmt.Sprintf("found active subset rule '%s' which will be kept", subset.Name), d.TrackingId)
				cleanedSubsetList = append(cleanedSubsetList, subset)
			} else {
				logger.Info(fmt.Sprintf("found inactive subset rule '%s' to be deleted", subset.Name), d.TrackingId)
			}
		}

		dr.Spec.Subsets = cleanedSubsetList
		err = UpdateDestinationRule(d, &dr)
		if err != nil {
			logger.Error(fmt.Sprintf("could not update destinationRule '%s' due to error '%s'", dr.Name, err), d.TrackingId)
			return err
		}
	}

	return nil
}

// Create returns a new subset to be posterior appended to destinationRules
func (d *DestinationRule) Create(s Shift) (*IstioRules, error) {
	newSubset := &v1alpha3.Subset{
		Name:   fmt.Sprintf("%s-%v-%s", d.Name, d.Build, d.Namespace),
		Labels: s.Traffic.PodSelector,
	}

	irl := &IstioRules{
		Subset: newSubset,
	}
	return irl, nil
}

// Validate checks if DestinationRule and Shift objects are correctly filled up
func (d *DestinationRule) Validate(s Shift) error {
	if d.Name == "" {
		return errors.New("empty 'name' attribute")
	}

	if d.Namespace == "" {
		return errors.New("empty 'namespace' attribute")
	}

	if d.Build == 0 {
		return errors.New("empty 'build' attribute")
	}

	if d.TrackingId == "" {
		return errors.New("empty 'trackingId' attribute")
	}

	if d.Istio == nil {
		return errors.New("nil istioClient object")
	}

	if len(s.Selector) == 0 {
		return errors.New("empty label-selector")
	}

	if s.Port == 0 {
		return errors.New("empty port")
	}

	if s.Port < 1023 {
		return errors.New("port not in range 1024 - 65535")
	}

	if s.Port > 65535 {
		return errors.New("port not in range 1024 - 65535")
	}

	if len(s.Traffic.PodSelector) == 0 {
		return errors.New("empty pod selector")
	}

	if !s.Traffic.Exact && !s.Traffic.Regexp {
		return errors.New("need 'exact' or 'regexp' flags")
	}

	return nil
}

/* Update a destinationRule with an existent subset based on Shift object
or just create a new one (based on Create() method)
*/
func (d *DestinationRule) Update(s Shift) error {
	newSubset := fmt.Sprintf("%s-%v-%s", d.Name, d.Build, d.Namespace)

	drs, err := d.List(s.Selector)
	if err != nil {
		return err
	}

	for _, dr := range drs.DList.Items {
		subsetExists := false
		for _, subsetValue := range dr.Spec.Subsets {
			if subsetValue.Name == newSubset {
				subsetExists = true
			}
		}

		if !subsetExists {
			irl, err := d.Create(s)
			if err != nil {
				logger.Error(fmt.Sprintf("could not create subset due to error '%s'", err), d.TrackingId)
				return err
			}

			dr.Spec.Subsets = append(dr.Spec.Subsets, irl.Subset)

			err = UpdateDestinationRule(d, &dr)
			if err != nil {
				logger.Error(fmt.Sprintf("could not update destinationRule '%s' due to error '%s'", dr.Name, err), d.TrackingId)
				return err
			}
		} else {
			logger.Info(fmt.Sprintf("subset '%s' already created", newSubset), d.TrackingId)
		}
	}

	return nil
}

// List will return all destinationRules which matches a k8s labelSelector
func (d *DestinationRule) List(selector map[string]string) (*IstioRouteList, error) {
	logger.Debug(fmt.Sprintf("Getting destinationRules which matches label-selector '%s'", selector), d.TrackingId)

	stringified, err := Stringify(d.TrackingId, selector)
	if err != nil {
		return &IstioRouteList{}, err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: stringified,
	}

	drs, err := d.Istio.NetworkingV1alpha3().DestinationRules(d.Namespace).List(listOptions)
	if err != nil {
		return &IstioRouteList{}, err
	}

	if len(drs.Items) <= 0 {
		return nil, errors.New(fmt.Sprintf("could not find any destinationRules which matched label-selector '%v'", listOptions.LabelSelector))
	}

	irl := IstioRouteList{
		DList: drs,
	}

	return &irl, nil
}

// UpdateDestinationRule updates a specific virtualService given an updated object
func UpdateDestinationRule(d *DestinationRule, destinationRule *v1alpha32.DestinationRule) error {
	logger.Info(fmt.Sprintf("Updating rule for destinationRule '%s'...", destinationRule.Name), d.TrackingId)
	_, err := d.Istio.NetworkingV1alpha3().DestinationRules(d.Namespace).Update(destinationRule)
	if err != nil {
		return err
	}

	return nil
}

// ValidateDestinationRuleList checks for inconsistencies in IstioRouteList.DList
func ValidateDestinationRuleList(irl *IstioRouteList) error {
	if len(irl.DList.Items) == 0 {
		return errors.New("empty destinationRules")
	}

	return nil
}
