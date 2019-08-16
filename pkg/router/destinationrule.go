package router

import (
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"istio.io/api/networking/v1alpha3"
)

type DestinationRuleRoute struct {
	Istio *versioned.Clientset
}

func (v *DestinationRuleRoute) Validate(route *Route) (v1alpha3.DestinationRule, error) {
	return v1alpha3.DestinationRule{}, nil

}

func (v *DestinationRuleRoute) Update(route *Route) error {
	return nil

}

func (v *DestinationRuleRoute) Delete(route *Route) error {
	return nil

}
