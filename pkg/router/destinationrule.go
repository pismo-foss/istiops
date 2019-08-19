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

type DestinationRuleRoute struct {
	Metadata Metadata
	Istio    *versioned.Clientset
}

func (v *DestinationRuleRoute) Validate(route *Route) (v1alpha3.DestinationRule, error) {
	fmt.Println("validating dr")
	return v1alpha3.DestinationRule{}, nil

}

func (v *DestinationRuleRoute) Update(route *Route) error {
	StringifyLabelSelector, err := utils.StringifyLabelSelector(v.Metadata.TrackingId, route.Selector.Labels)
	if err != nil {
		fmt.Println("null drs")
	}

	listOptions := metav1.ListOptions{
		LabelSelector: StringifyLabelSelector,
	}

	drs, err := GetAllDestinationRules(v, route, listOptions)
	if err != nil {
		fmt.Println("null drs")
	}

	for _, dr := range drs.Items {
		newSubset := &v1alpha3.Subset{
			Name:   fmt.Sprintf("%s-%v-%s", v.Metadata.Name, v.Metadata.Build, v.Metadata.Namespace),
			Labels: route.Traffic.PodSelector,
		}
		updatedDr, err := createSubset(dr, newSubset)
		if err != nil {
			utils.Fatal(fmt.Sprintf("could not create subset due to error '%s'", err), v.Metadata.TrackingId)
		}

		err = UpdateDestinationRule(v, updatedDr)
		if err != nil {
			utils.Fatal(fmt.Sprintf("could not update destinationRule '%s' due to error '%s'", updatedDr.Name, err), v.Metadata.TrackingId)
		}

	}
	return nil

}

func (v *DestinationRuleRoute) Delete(route *Route) error {
	return nil

}

// GetAllVirtualservices returns all istio resources 'virtualservices'
func GetAllDestinationRules(drRoute *DestinationRuleRoute, route *Route, listOptions metav1.ListOptions) (*v1alpha32.DestinationRuleList, error) {
	utils.Info(fmt.Sprintf("Finding destinationRules which matches selector '%s'...", listOptions.LabelSelector), drRoute.Metadata.Namespace)
	drs, err := drRoute.Istio.NetworkingV1alpha3().DestinationRules(drRoute.Metadata.Namespace).List(listOptions)
	if err != nil {
		return nil, err
	}

	utils.Info(fmt.Sprintf("Found a total of '%d' destinationRules", len(drs.Items)), drRoute.Metadata.TrackingId)
	return drs, nil
}

func createSubset(dr v1alpha32.DestinationRule, newSubset *v1alpha3.Subset) (*v1alpha32.DestinationRule, error) {
	for _, subsetValue := range dr.Spec.Subsets {
		if subsetValue.Name == newSubset.Name {
			// remove item from slice
			return nil, errors.New(fmt.Sprintf("Found already existent subset '%s', refusing to update", subsetValue.Name))
		}
	}

	dr.Spec.Subsets = append(dr.Spec.Subsets, newSubset)

	return &dr, nil
}

// UpdateDestinationRule updates a specific virtualService given an updated object
func UpdateDestinationRule(drRoute *DestinationRuleRoute, destinationRule *v1alpha32.DestinationRule) error {
	utils.Info(fmt.Sprintf("Updating rule for destinationRule '%s'...", destinationRule.Name), drRoute.Metadata.TrackingId)
	_, err := drRoute.Istio.NetworkingV1alpha3().DestinationRules(drRoute.Metadata.Namespace).Update(destinationRule)
	if err != nil {
		return err
	}
	return nil
}
