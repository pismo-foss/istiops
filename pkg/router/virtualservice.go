package router

import (
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
)

type VirtualService struct {
	Item  *v1alpha32.VirtualService
	Route *Route
}

func (v *VirtualService) Add(route *Route) error {
	return nil
}

func (v *VirtualService) Update(route *Route) error {
	return nil

}

func (v *VirtualService) Delete(route *Route) error {
	return nil
}
