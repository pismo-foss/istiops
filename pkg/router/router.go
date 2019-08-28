package router

import (
	"errors"
	"fmt"
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type Router interface {
	Create(s *Shift) (*IstioRules, error)
	Validate(s *Shift) error
	Update(s *Shift) error
	Clear(s *Shift) error
	List(opts metav1.ListOptions) (*IstioRouteList, error)
}

type Shift struct {
	Port     uint32
	Hostname string
	Selector *Selector
	Traffic  *Traffic
}

type Traffic struct {
	PodSelector    map[string]string
	RequestHeaders map[string]string
	Weight         int32
}

type Selector struct {
	Labels map[string]string
}

type IstioRules struct {
	MatchDestination *v1alpha3.HTTPRoute
	Subset           *v1alpha3.Subset
}

type IstioRouteList struct {
	VirtualServiceList   *v1alpha32.VirtualServiceList
	DestinationRulesList *v1alpha32.DestinationRuleList
}

// StringifyLabelSelector returns a k8s selector string based on given map. Ex: "key=value,key2=value2"
func StringifyLabelSelector(cid string, labelSelector map[string]string) (string, error) {

	var labelsPair []string

	for key, value := range labelSelector {
		labelsPair = append(labelsPair, fmt.Sprintf("%s=%s", key, value))
	}

	if len(labelsPair) == 0 {
		return "", errors.New("got an empty labelSelector")
	}

	return strings.Join(labelsPair[:], ","), nil
}
