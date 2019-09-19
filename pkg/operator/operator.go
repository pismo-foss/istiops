package operator

import (
	"github.com/pismo/istiops/pkg/router"
)

type Operator interface {
	Get(selector map[string]string) (router.IstioRouteList, error)
	Update(shift router.Shift) error
	Clear(shift router.Shift) error
}
