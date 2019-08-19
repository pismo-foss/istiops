package router

type Router interface {
	Validate(r Route) error
	Update(r Route) error
	Delete(r Route) error
}

type Route struct {
	Port     uint32
	Hostname string
	Selector *Selector
	Headers  map[string]string
	Weight   int32
}

type Selector struct {
	ResourceSelector map[string]string
	PodSelector      map[string]string
}
