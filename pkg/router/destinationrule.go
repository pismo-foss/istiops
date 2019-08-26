package router

import (
	"errors"
	"fmt"

	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"github.com/pismo/istiops/utils"
	"istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DrMetadata struct {
	TrackingId string
	Name       string
	Namespace  string
	Build      uint32
}

type DestinationRule struct {
	Metadata DrMetadata
	Istio    *versioned.Clientset
}

func (v *DestinationRule) Validate(s *Shift) error {
	newSubset := fmt.Sprintf("%s-%v-%s", v.Metadata.Name, v.Metadata.Build, v.Metadata.Namespace)

	StringifyLabelSelector, err := utils.StringifyLabelSelector(v.Metadata.TrackingId, s.Selector.Labels)
	if err != nil {
		return err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: StringifyLabelSelector,
	}

	drs, err := GetAllDestinationRules(v, listOptions)
	if err != nil {
		return err
	}

	for _, dr := range drs.Items {
		utils.Info(fmt.Sprintf("Validating destinationRule '%s'", dr.Name), v.Metadata.TrackingId)
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
	StringifyLabelSelector, err := utils.StringifyLabelSelector(v.Metadata.TrackingId, s.Selector.Labels)
	if err != nil {
		return err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: StringifyLabelSelector,
	}

	drs, err := GetAllDestinationRules(v, listOptions)
	if err != nil {
		fmt.Println("null drs")
		return err
	}

	for _, dr := range drs.Items {
		newSubset := &v1alpha3.Subset{
			Name:   fmt.Sprintf("%s-%v-%s", v.Metadata.Name, v.Metadata.Build, v.Metadata.Namespace),
			Labels: s.Traffic.PodSelector,
		}
		updatedDr, err := createSubset(dr, newSubset)
		if err != nil {
			utils.Error(fmt.Sprintf("could not create subset due to error '%s'", err), v.Metadata.TrackingId)
			return err
		}

		err = UpdateDestinationRule(v, updatedDr)
		if err != nil {
			utils.Error(fmt.Sprintf("could not update destinationRule '%s' due to error '%s'", updatedDr.Name, err), v.Metadata.TrackingId)
			return err
		}

	}
	return nil

}

func (v *DestinationRule) Clear(s *Shift) error {
	StringifyLabelSelector, err := utils.StringifyLabelSelector(v.Metadata.TrackingId, s.Selector.Labels)
	if err != nil {
		return err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: StringifyLabelSelector,
	}

	drs, err := GetAllDestinationRules(v, listOptions)
	if err != nil {
		return err
	}

	for _, dr := range drs.Items {
		dr.Spec.Subsets = []*v1alpha3.Subset{}

		utils.Info(fmt.Sprintf("Clearing all destinationRules subsets from '%s'...", dr.Name), v.Metadata.TrackingId)
		err := UpdateDestinationRule(v, &dr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *DestinationRule) Delete(s *Shift) error {
	return nil

}

// GetAllDestinationRules returns all istio resources 'virtualservices'
func GetAllDestinationRules(dr *DestinationRule, listOptions metav1.ListOptions) (*v1alpha32.DestinationRuleList, error) {
	utils.Info(fmt.Sprintf("Finding destinationRules which matches selector '%s'...", listOptions.LabelSelector), dr.Metadata.TrackingId)

	drs, err := dr.Istio.NetworkingV1alpha3().DestinationRules(dr.Metadata.Namespace).List(listOptions)
	if err != nil {
		return nil, err
	}

	if len(drs.Items) <= 0 {
		return nil, errors.New("could not find any destinationRules")
	}

	utils.Info(fmt.Sprintf("Found a total of '%d' destinationRules to work it", len(drs.Items)), dr.Metadata.TrackingId)
	return drs, nil
}

func createSubset(dr v1alpha32.DestinationRule, newSubset *v1alpha3.Subset) (*v1alpha32.DestinationRule, error) {

	dr.Spec.Subsets = append(dr.Spec.Subsets, newSubset)

	return &dr, nil
}

// UpdateDestinationRule updates a specific virtualService given an updated object
func UpdateDestinationRule(dr *DestinationRule, destinationRule *v1alpha32.DestinationRule) error {
	utils.Info(fmt.Sprintf("Updating rule for destinationRule '%s'...", destinationRule.Name), dr.Metadata.TrackingId)
	_, err := dr.Istio.NetworkingV1alpha3().DestinationRules(dr.Metadata.Namespace).Update(destinationRule)
	if err != nil {
		return err
	}
	return nil
}
