package operator

import "github.com/pismo/istiops/pkg/router"

type Operator interface {
	Create(s *router.Shift) error
	Delete(s *router.Shift) error
	Update(s *router.Shift) error
	Clear(s *router.Shift) error
}
