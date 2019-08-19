package router

import (
	"fmt"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"istio.io/api/networking/v1alpha3"
)

type VirtualServiceRoute struct {
	Metadata Metadata
	Istio    *versioned.Clientset
}

func (v *VirtualServiceRoute) Validate(route *Route) (v1alpha3.VirtualService, error) {
	fmt.Println("validating vr")
	return v1alpha3.VirtualService{}, nil

}

func (v *VirtualServiceRoute) Update(route *Route) error {
	fmt.Println("updating virtualservice")
	return nil

}

func (v *VirtualServiceRoute) Delete(route *Route) error {
	return nil

}

// GetAllVirtualServices returns all istio resources 'virtualservices'
//func GetAllVirtualServices(vsRoute router.VirtualServiceRoute, ips *Istiops, listOptions metav1.ListOptions) (*v1alpha32.VirtualServiceList, error) {
//	utils.Info(fmt.Sprintf("Finding virtualServices which matches selector '%s'...", listOptions.LabelSelector), ips.TrackingId)
//	vss, err := vsRoute.Istio.NetworkingV1alpha3().VirtualServices(ips.Namespace).List(listOptions)
//	if err != nil {
//		return nil, err
//	}
//
//	utils.Info(fmt.Sprintf("Found a total of '%d' virtualServices", len(vss.Items)), ips.TrackingId)
//	return vss, nil
//}
