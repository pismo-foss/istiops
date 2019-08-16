package router

import (
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"istio.io/api/networking/v1alpha3"
)

type VirtualServiceRoute struct {
	Istio *versioned.Clientset
}

func (v *VirtualServiceRoute) Validate(route *Route) (v1alpha3.VirtualService, error) {
	return v1alpha3.VirtualService{}, nil

}

func (v *VirtualServiceRoute) Update(route *Route) error {
	return nil

}

func (v *VirtualServiceRoute) Delete(route *Route) error {
	return nil

}
