package operator

import (
	"github.com/pismo/istiops/pkg/router"
	"github.com/pkg/errors"
)

type Istiops struct {
	DrRouter router.Router
	VsRouter router.Router
}

func (ips *Istiops) Create(r *router.Shift) error {
	VsRouter := ips.VsRouter
	err := VsRouter.Validate(r)
	if err != nil {
		return err
	}

	DrRouter := ips.DrRouter
	err = DrRouter.Validate(r)
	if err != nil {
		return err
	}

	return nil
}

func (ips *Istiops) Update(r *router.Shift) error {
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
func (ips *Istiops) Clear(s *router.Shift) error {
	DrRouter := ips.DrRouter
	VsRouter := ips.VsRouter
	var err error

	err = DrRouter.Clear(s)
	if err != nil {
		return err
	}

	err = VsRouter.Clear(s)
	if err != nil {
		return err
	}

	// Clean dr rules ?

	return nil
}
