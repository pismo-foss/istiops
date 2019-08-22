package operator

import (
	"fmt"
	"github.com/pismo/istiops/pkg/router"
	"github.com/pismo/istiops/utils"
)

type Istiops struct {
	Shift    *router.Shift
	VsRouter *router.VirtualService
	DrRouter *router.DestinationRule
}

func (ips *Istiops) Get(r *router.Shift) error {
	return nil
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

func (ips *Istiops) Delete(r *router.Shift) error {
	fmt.Println("Initializing something")
	fmt.Println("", "cid")
	return nil
}

func (ips *Istiops) Update(r *router.Shift) error {
	if len(r.Selector.Labels) == 0 || len(r.Traffic.PodSelector) == 0 {
		utils.Fatal(fmt.Sprintf("Selectors must not be empty otherwise istiops won't be able to find any resources."), "")
	}

	DrRouter := ips.DrRouter
	VsRouter := ips.VsRouter
	var err error

	err = DrRouter.Validate(r)
	if err != nil {
		utils.Fatal(fmt.Sprintf("%s", err), ips.DrRouter.Metadata.TrackingId)
	}
	err = DrRouter.Update(r)
	if err != nil {
		utils.Fatal(fmt.Sprintf("%s", err), ips.DrRouter.Metadata.TrackingId)
	}

	err = VsRouter.Validate(r)
	if err != nil {
		utils.Fatal(fmt.Sprintf("%s", err), ips.VsRouter.Metadata.TrackingId)
	}
	err = VsRouter.Update(r)
	if err != nil {
		utils.Fatal(fmt.Sprintf("%s", err), ips.VsRouter.Metadata.TrackingId)
	}

	if r.Traffic.Weight > 0 {
		// update router to serve percentage
		if err != nil {
			utils.Fatal(fmt.Sprintf("Could no create resource due to an error '%s'", err), "")
		}
	}

	return nil
}

// ClearRules will remove any destination & virtualService rules except the main one (provided by client).
// Ex: URI or Prefix
func (ips *Istiops) Clear(s *router.Shift) error {
	DrRouter := ips.DrRouter
	VsRouter := ips.VsRouter
	var err error

	err = DrRouter.Clear(ips.Shift)
	if err != nil {
		utils.Fatal(fmt.Sprintf("%s", err), ips.DrRouter.Metadata.TrackingId)
	}

	err = VsRouter.Clear(ips.Shift)
	if err != nil {
		utils.Fatal(fmt.Sprintf("%s", err), ips.DrRouter.Metadata.TrackingId)
	}

	// Clean dr rules ?

	return nil
}
