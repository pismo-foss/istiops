package router

import (
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"istio.io/api/networking/v1alpha3"
)

type VirtualService struct {
	Istio *versioned.Clientset
}

func (v *VirtualService) Validate(route *Route) (v1alpha3.VirtualService, error) {
	return v1alpha3.VirtualService{}, nil

}

func (v *VirtualService) Update(route *Route) error {
	return nil

}

func (v *VirtualService) Delete(route *Route) error {
	return nil

}
