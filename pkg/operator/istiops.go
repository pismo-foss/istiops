package operator

import (
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/pismo/istiops/pkg/router"
	"github.com/pkg/errors"
)

type Router interface {
	Create(shift router.Shift) (*router.IstioRules, error)
	Validate(shift router.Shift) error
	Update(shift router.Shift) error
	Clear(shift router.Shift) error
	List(selector map[string]string) (*router.IstioRouteList, error)
}

type Istiops struct {
	DrRouter Router
	VsRouter Router
}

func (ips *Istiops) Get(selector map[string]string) ([]v1alpha32.VirtualService, error) {
	VsRouter := ips.VsRouter
	ivl, err := VsRouter.List(selector)
	if err != nil {
		return []v1alpha32.VirtualService{}, err
	}

	if len(ivl.VList.Items) == 0 {
		return []v1alpha32.VirtualService{}, errors.New("empty virtualServices")
	}
	return ivl.VList.Items, nil
}

func (ips *Istiops) Update(shift router.Shift) error {
	if len(shift.Selector) == 0 {
		return errors.New("label-selector must exists in need to find resources")
	}

	if len(shift.Traffic.PodSelector) == 0 {
		return errors.New("pod-selector must exists in need to find traffic destination")
	}

	DrRouter := ips.DrRouter
	VsRouter := ips.VsRouter
	var err error

	err = DrRouter.Validate(shift)
	if err != nil {
		return err
	}
	err = VsRouter.Validate(shift)
	if err != nil {
		return err
	}
	err = DrRouter.Update(shift)
	if err != nil {
		return err
	}
	err = VsRouter.Update(shift)
	if err != nil {
		return err
	}

	return nil
}

// ClearRules will remove any destination & virtualService rules except the main one (provided by client).
// Ex: URI or Prefix
func (ips *Istiops) Clear(shift router.Shift) error {
	DrRouter := ips.DrRouter
	VsRouter := ips.VsRouter
	var err error

	// in this scenario virtualService must be cleaned before the DestinationRule
	err = VsRouter.Clear(shift)
	if err != nil {
		return err
	}

	err = DrRouter.Clear(shift)
	if err != nil {
		return err
	}

	return nil
}
