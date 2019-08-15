package router

import (
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/heptio/contour/apis/generated/clientset/versioned"
)

type VirtualService struct {
	Item  *v1alpha32.VirtualService
	Istio *versioned.Clientset
}
