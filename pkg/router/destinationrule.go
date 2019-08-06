package router

import "istio.io/api/networking/v1alpha3"

type DestinationRule {
	IstioDestination
	Name
	Namespace
}

type DestinationRule struct {
	CID  string
	item *v1alpha3.DestinationRule
}

func (v *DestinationRule) Add(route Route) error {
	return nil
}

func (v *DestinationRule) Update(route Route) error {
	return nil

}

func (v *DestinationRule) Delete(route Route) error {
	return nil

}
