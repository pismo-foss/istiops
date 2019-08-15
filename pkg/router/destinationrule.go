package router

import (
	"fmt"

	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/heptio/contour/apis/generated/clientset/versioned"
)

type DestinationRule struct {
	Item  *v1alpha32.DestinationRule
	Istio *versioned.Clientset
}

func (v *DestinationRule) Validate(route *Route) error {
	fmt.Println(v.Item.Name)
	return nil

}

func (v *DestinationRule) Update(route *Route) error {
	return nil

}

func (v *DestinationRule) Delete(route *Route) error {
	return nil

}
