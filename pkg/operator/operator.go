package operator

import (
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/pismo/istiops/pkg/router"
)

type Operator interface {
	Get(selector map[string]string) ([]v1alpha32.VirtualService, error)
	Update(shift router.Shift) error
	Clear(shift router.Shift) error
}
