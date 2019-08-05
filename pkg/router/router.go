package router

import "istio.io/api/networking/v1alpha3"

type DestinationRule struct {
	CID  string
	item *v1alpha3.DestinationRule
}

type VirtualService struct {
	CID  string
	item *v1alpha3.VirtualService
}

type Route struct {
	Destination *v1alpha3.RouteDestination
}

type Router interface {
	Add(route Route) error
	Update(route Route) error
	Delete(route Route) error
}
