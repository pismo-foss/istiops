package operator

import (
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/pismo/istiops/pkg/router"
	"github.com/pkg/errors"
)

type Istiops struct {
	DrRouter Router
	VsRouter Router
}

func (ips *Istiops) Get(r router.Shift) ([]v1alpha32.VirtualService, error) {
	VsRouter := ips.VsRouter
	ivl, err := VsRouter.List(r)
	if err != nil {
		return []v1alpha32.VirtualService{}, err
	}

	if len(ivl.VList.Items) == 0 {
		return []v1alpha32.VirtualService{}, errors.New("empty virtualServices")
	}
	return ivl.VList.Items, nil
}

func (ips *Istiops) Update(r router.Shift) error {
	if len(r.Selector.Labels) == 0 {
		return errors.New("label-selector must exists in need to find resources")
	}

	if len(r.Traffic.PodSelector) == 0 {
		return errors.New("pod-selector must exists in need to find traffic destination")
	}

	DrRouter := ips.DrRouter
	VsRouter := ips.VsRouter
	var err error

	err = DrRouter.Validate(r)
	if err != nil {
		return err
	}
	err = VsRouter.Validate(r)
	if err != nil {
		return err
	}
	err = DrRouter.Update(r)
	if err != nil {
		return err
	}
	err = VsRouter.Update(r)
	if err != nil {
		return err
	}

	return nil
}

// ClearRules will remove any destination & virtualService rules except the main one (provided by client).
// Ex: URI or Prefix
func (ips *Istiops) Clear(s router.Shift) error {
	DrRouter := ips.DrRouter
	VsRouter := ips.VsRouter
	var err error

	// in this scenario virtualService must be cleaned before the DestinationRule
	err = VsRouter.Clear(s)
	if err != nil {
		return err
	}

	err = DrRouter.Clear(s)
	if err != nil {
		return err
	}

	return nil
}
