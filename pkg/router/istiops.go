import "istio.io/api/networking/v1alpha3"

type Istiops struct {
	client
	istio
	blabla
}

type Route struct {
	Destination *v1alpha3.RouteDestination
	Labels      *Labels
}

func (i Istiops) Validate(r Route) {

	r.dr.istioDr = i.client.Get(r.dr.Name)

}
