package router

import "istio.io/api/networking/v1alpha3"

type VirtualService struct {
	CID  string
	item *v1alpha3.VirtualService
}

func (v *VirtualService) Add(route Route) error {
	return nil

}

func (v *VirtualService) Update(route Route) error {
	return nil

}

func (v *VirtualService) Delete(route Route) error {
	return nil

}
