package router

import v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"

type Router interface {
	Validate(s *Shift) error
	Update(s *Shift) error
	Delete(s *Shift) error
}

type Shift struct {
	Port     uint32
	Hostname string
	Selector *Selector
	Traffic  *Traffic
}

type Traffic struct {
	PodSelector    map[string]string
	RequestHeaders map[string]string
	Weight         int32
}

type Selector struct {
	Labels map[string]string
}

type IstioRouteList struct {
	VirtualServiceList   *v1alpha32.VirtualServiceList
	DestinationRulesList *v1alpha32.DestinationRuleList
}
