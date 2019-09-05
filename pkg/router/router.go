package router

import (
	"errors"
	"fmt"
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/fake"
	"istio.io/api/networking/v1alpha3"
	"strings"
)

type Client struct {
	Versioned *versioned.Clientset
	Fake      *fake.Clientset
}

type Shift struct {
	Port     uint32
	Hostname string
	Selector map[string]string
	Traffic  Traffic
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
	VList *v1alpha32.VirtualServiceList
	DList *v1alpha32.DestinationRuleList
}

// StringifyLabelSelector returns a k8s selector string based on given map. Ex: "key=value,key2=value2"
func Stringify(cid string, labelSelector map[string]string) (string, error) {

	var labelsPair []string

	for key, value := range labelSelector {
		labelsPair = append(labelsPair, fmt.Sprintf("%s=%s", key, value))
	}

	if len(labelsPair) == 0 {
		return "", errors.New("got an empty labelSelector")
	}

	return strings.Join(labelsPair[:], ","), nil
}

func Mapify(cid string, labelSelector string) (map[string]string, error) {
	mapLabels := map[string]string{}

	if labelSelector == "" {
		return nil, errors.New("got an empty labelSelector string")
	}

	if !strings.Contains(labelSelector, "=") {
		return nil, errors.New("missing '=' operator for labelSelector")
	}

	splitedLabels := strings.Split(labelSelector, ",")
	for _, value := range splitedLabels {
		parsedLabels := strings.Split(value, "=")
		mapLabels[parsedLabels[0]] = parsedLabels[1]
	}

	if len(mapLabels) == 0 {
		return nil, errors.New("empty label selector")
	}

	return mapLabels, nil
}
