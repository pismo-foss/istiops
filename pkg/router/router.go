package router

import "istio.io/api/networking/v1alpha3"

type RouteInterface interface {
	Validate(route Route) error
	Update(route Route) error
	Delete(route Route) error
}

type Route struct {
	Destination *v1alpha3.RouteDestination
}

type Subset struct {
	Subset *v1alpha3.Subset
}

type TrafficShift struct {
	Headers map[string]string
	Weight  uint32
}
