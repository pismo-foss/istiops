package router

type Router interface {
	Validate(route Route) error
	Update(route Route) error
	Delete(route Route) error
}
