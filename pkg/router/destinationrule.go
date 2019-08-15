package router

import (
	"fmt"
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
)

type DestinationRule struct {
	Item *v1alpha32.DestinationRule
}

func (v *DestinationRule) Add(route *Route) error {
	fmt.Println(v.Item.Name)
	return nil

}

func (v *DestinationRule) Update(route *Route) error {
	return nil

}

func (v *DestinationRule) Delete(route *Route) error {
	return nil

}
