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
	Istio    *versioned.Clientset
}

func (v *DestinationRule) Validate(s *Shift) error {
	newSubset := fmt.Sprintf("%s-%v-%s", v.Name, v.Build, v.Namespace)

	StringifyLabelSelector, err := StringifyLabelSelector(v.TrackingId, s.Selector.Labels)
	if err != nil {
		return err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: StringifyLabelSelector,
	}

	drs, err := v.List(listOptions)
	if err != nil {
		return err
	}

	for _, dr := range drs.DestinationRulesList.Items {
		logger.Info(fmt.Sprintf("Validating destinationRule '%s'", dr.Name), v.TrackingId)
		for _, subsetValue := range dr.Spec.Subsets {
			if subsetValue.Name == newSubset {
				// remove item from slice
				return errors.New(fmt.Sprintf("Found already existent subset '%s', refusing to update", subsetValue.Name))
			}
		}
	}

	return nil

}

func (v *DestinationRule) Update(s *Shift) error {
	StringifyLabelSelector, err := StringifyLabelSelector(v.TrackingId, s.Selector.Labels)
	if err != nil {
		return err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: StringifyLabelSelector,
	}

	drs, err := v.List(listOptions)
	if err != nil {
		fmt.Println("null drs")
		return err
	}

	for _, dr := range drs.DestinationRulesList.Items {
		newSubset := &v1alpha3.Subset{
			Name:   fmt.Sprintf("%s-%v-%s", v.Name, v.Build, v.Namespace),
			Labels: s.Traffic.PodSelector,
		}
		updatedDr, err := createSubset(dr, newSubset)
		if err != nil {
			logger.Error(fmt.Sprintf("could not create subset due to error '%s'", err), v.TrackingId)
			return err
		}

		err = UpdateDestinationRule(v, updatedDr)
		if err != nil {
			logger.Error(fmt.Sprintf("could not update destinationRule '%s' due to error '%s'", updatedDr.Name, err), v.TrackingId)
			return err
		}

	}
	return nil

}

func (v *DestinationRule) Clear(s *Shift) error {
	StringifyLabelSelector, err := StringifyLabelSelector(v.TrackingId, s.Selector.Labels)
	if err != nil {
		return err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: StringifyLabelSelector,
	}

	drs, err := v.List(listOptions)
	if err != nil {
		return err
	}

	for _, dr := range drs.DestinationRulesList.Items {
		dr.Spec.Subsets = []*v1alpha3.Subset{}

		logger.Info(fmt.Sprintf("Clearing all destinationRules subsets from '%s'...", dr.Name), v.TrackingId)
		err := UpdateDestinationRule(v, &dr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *DestinationRule) List(opts metav1.ListOptions) (*IstioRouteList, error) {
	drs, err := v.Istio.NetworkingV1alpha3().DestinationRules(v.Namespace).List(opts)
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

func createSubset(dr v1alpha32.DestinationRule, newSubset *v1alpha3.Subset) (*v1alpha32.DestinationRule, error) {

	dr.Spec.Subsets = append(dr.Spec.Subsets, newSubset)

	return &dr, nil
}

// UpdateDestinationRule updates a specific virtualService given an updated object
func UpdateDestinationRule(dr *DestinationRule, destinationRule *v1alpha32.DestinationRule) error {
	logger.Info(fmt.Sprintf("Updating rule for destinationRule '%s'...", destinationRule.Name), dr.TrackingId)
	_, err := dr.Istio.NetworkingV1alpha3().DestinationRules(dr.Namespace).Update(destinationRule)
	if err != nil {
		return err
	}
	return nil
}
