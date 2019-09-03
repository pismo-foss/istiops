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
	Istio      Client
}

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

func (d *DestinationRule) Validate(s Shift) error {

	return nil

}

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
				// remove item from slice
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
			logger.Info(fmt.Sprintf("subset '%s'already created", newSubset), d.TrackingId)
		}
	}

	return nil
}

func (d *DestinationRule) Clear(s Shift) error {

	return nil
}

func (d *DestinationRule) List(selector map[string]string) (*IstioRouteList, error) {
	logger.Info(fmt.Sprintf("Getting destinationRules which matches label-selector '%s'", selector), d.TrackingId)

	stringified, err := Stringify(d.TrackingId, selector)
	if err != nil {
		return &IstioRouteList{}, err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: stringified,
	}

	drs, err := d.Istio.Versioned.NetworkingV1alpha3().DestinationRules(d.Namespace).List(listOptions)
	if err != nil {
		return nil, err
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
	_, err := d.Istio.Versioned.NetworkingV1alpha3().DestinationRules(d.Namespace).Update(destinationRule)
	if err != nil {
		return err
	}
	return nil
}
