package router

import (
	"errors"
	"fmt"
	"github.com/pismo/istiops/pkg/logger"

	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DestinationRule struct {
	TrackingId string
	Name       string
	Namespace  string
	Build      uint32
	Istio      *versioned.Clientset
}

func (d *DestinationRule) Create(s *Shift) (*IstioRules, error) {
	newSubset := &v1alpha3.Subset{
		Name:   fmt.Sprintf("%s-%v-%s", d.Name, d.Build, d.Namespace),
		Labels: s.Traffic.PodSelector,
	}

	irl := &IstioRules{
		Subset: newSubset,
	}
	return irl, nil
}

func (d *DestinationRule) Validate(s *Shift) error {
	newSubset := fmt.Sprintf("%s-%v-%s", d.Name, d.Build, d.Namespace)

	StringifyLabelSelector, err := StringifyLabelSelector(d.TrackingId, s.Selector.Labels)
	if err != nil {
		return err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: StringifyLabelSelector,
	}

	drs, err := d.List(listOptions)
	if err != nil {
		return err
	}

	for _, dr := range drs.DestinationRulesList.Items {
		logger.Info(fmt.Sprintf("Validating destinationRule '%s'", dr.Name), d.TrackingId)
		for _, subsetValue := range dr.Spec.Subsets {
			if subsetValue.Name == newSubset {
				// remove item from slice
				return errors.New(fmt.Sprintf("Found already existent subset '%s', refusing to update", subsetValue.Name))
			}
		}
	}

	return nil

}

func (d *DestinationRule) Update(s *Shift) error {
	StringifyLabelSelector, err := StringifyLabelSelector(d.TrackingId, s.Selector.Labels)
	if err != nil {
		return err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: StringifyLabelSelector,
	}

	drs, err := d.List(listOptions)
	if err != nil {
		fmt.Println("null drs")
		return err
	}

	for _, dr := range drs.DestinationRulesList.Items {
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

	}
	return nil

}

func (d *DestinationRule) Clear(s *Shift) error {
	StringifyLabelSelector, err := StringifyLabelSelector(d.TrackingId, s.Selector.Labels)
	if err != nil {
		return err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: StringifyLabelSelector,
	}

	drs, err := d.List(listOptions)
	if err != nil {
		return err
	}

	for _, dr := range drs.DestinationRulesList.Items {
		dr.Spec.Subsets = []*v1alpha3.Subset{}

		logger.Info(fmt.Sprintf("Clearing all destinationRules subsets from '%s'...", dr.Name), d.TrackingId)
		err := UpdateDestinationRule(d, &dr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DestinationRule) List(opts metav1.ListOptions) (*IstioRouteList, error) {
	drs, err := d.Istio.NetworkingV1alpha3().DestinationRules(d.Namespace).List(opts)
	if err != nil {
		return nil, err
	}

	if len(drs.Items) <= 0 {
		return nil, errors.New(fmt.Sprintf("could not find any destinationRules which matched label-selector '%v'", opts.LabelSelector))
	}

	irl := IstioRouteList{
		DestinationRulesList: drs,
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
