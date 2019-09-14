package operator

import (
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

// Get will return a list of istio resources: destinationRules & virtualServices
func (ips *Istiops) Get(selector map[string]string) (router.IstioRouteList, error) {
	DrRouter := ips.DrRouter
	dsl, err := DrRouter.List(selector)
	if err != nil {
		return router.IstioRouteList{}, err
	}

	err = router.ValidateDestinationRuleList(dsl)
	if err != nil {
		return router.IstioRouteList{}, err
	}

	VsRouter := ips.VsRouter
	vsl, err := VsRouter.List(selector)
	if err != nil {
		return router.IstioRouteList{}, err
	}

	err = router.ValidateVirtualServiceList(vsl)
	if err != nil {
		return router.IstioRouteList{}, err
	}

	ivl := router.IstioRouteList{}
	ivl.DList = dsl.DList
	ivl.VList = vsl.VList

	return ivl, nil
}

// Update will update (and create if not exists) a route rule based on given Shift struct
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

// ClearRules will remove any destination & virtualService route rules except the main one (provided by client).
func (ips *Istiops) Clear(shift router.Shift) error {
	DrRouter := ips.DrRouter
	VsRouter := ips.VsRouter
	var err error

	// to bypass
	shift.Traffic.RequestHeaders = map[string]string{"clear": "true"}

	err = VsRouter.Validate(shift)
	if err != nil {
		return err
	}

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
