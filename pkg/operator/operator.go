package operator

import "github.com/pismo/istiops/pkg/router"

type Operator interface {
	Get(ir *router.Route) error
	Create(ir *router.Route) error
	Delete(ir *router.Route) error
	Update(ir *router.Route) error
	Clear(map[string]string) error
}
