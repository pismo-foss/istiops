package router

import (
	"istio.io/api/networking/v1alpha3"
)

type Route struct {
	Port     uint32
	Hostname string
	Selector map[string]string
	Headers  map[string]string
	Weight   int32
}

type Router interface {
	Validate(route Route) error
	Update(route Route) error
	Delete(route Route) error
}

type Subset struct {
	Subset *v1alpha3.Subset
}

type TrafficShift struct {
	Headers map[string]string
	Percent int32
}
