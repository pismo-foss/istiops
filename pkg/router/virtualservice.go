package router

import (
	"fmt"
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"github.com/pismo/istiops/utils"
	"istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VsMetadata struct {
	TrackingId string
	Name       string
	Namespace  string
	Build      uint32
}

type VirtualService struct {
	Metadata VsMetadata
	Istio    *versioned.Clientset
}

func (v *VirtualService) Validate(s *Shift) (v1alpha3.VirtualService, error) {
	fmt.Println("validating vr")
	return v1alpha3.VirtualService{}, nil

}

func (v *VirtualService) Update(s *Shift) error {
	fmt.Println("updating virtualservice")
	subsetName := fmt.Sprintf("%s-%v-%s", v.Metadata.Name, v.Metadata.Build, v.Metadata.Namespace)

	StringifyLabelSelector, err := utils.StringifyLabelSelector(v.Metadata.TrackingId, s.Selector.Labels)
	if err != nil {
		fmt.Println("null drs")
		return err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: StringifyLabelSelector,
	}

	vss, err := GetAllVirtualServices(v, s, listOptions)
	for _, vs := range vss.Items {
		for _, httpValue := range vs.Spec.Http {
			for _, routeValue := range httpValue.Route {
				if routeValue.Destination.Host == subsetName {
					utils.Error(fmt.Sprintf("Updating virtualservice rule '%s'", routeValue.Destination.Host), v.Metadata.TrackingId)
					return err
				}
			}
		}
	}

	return nil

}

func (v *VirtualService) Delete(s *Shift) error {
	return nil

}

// GetAllVirtualServices returns all istio resources 'virtualservices'
func GetAllVirtualServices(vsRoute *VirtualService, s *Shift, listOptions metav1.ListOptions) (*v1alpha32.VirtualServiceList, error) {
	utils.Info(fmt.Sprintf("Getting all virtualservices..."), vsRoute.Metadata.TrackingId)

	vss, err := vsRoute.Istio.NetworkingV1alpha3().VirtualServices(vsRoute.Metadata.Namespace).List(listOptions)
	if err != nil {
		return nil, err
	}

	return vss, nil
}
