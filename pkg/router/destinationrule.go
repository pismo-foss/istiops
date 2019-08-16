package router

import (
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"istio.io/api/networking/v1alpha3"
)

type DestinationRule struct {
	Istio *versioned.Clientset
}

func (v *DestinationRule) Validate(route *Route) (v1alpha3.DestinationRule, error) {
	return v1alpha3.DestinationRule{}, nil

}

func (v *DestinationRule) Update(route *Route) error {
	return nil

}

func (v *DestinationRule) Delete(route *Route) error {
	return nil

}
