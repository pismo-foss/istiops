package router

import (
	"errors"
	"fmt"
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	networkingv1alpha3 "github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/typed/networking/v1alpha3"
	"istio.io/api/networking/v1alpha3"
	appsV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"strings"
)

type IstioClientInterface interface {
	NetworkingV1alpha3() networkingv1alpha3.NetworkingV1alpha3Interface
}

type KubeClientInterface interface {
	AppsV1() appsV1.AppsV1Interface
	CoreV1() coreV1.CoreV1Interface
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
	Exact          bool
	Regexp         bool
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

// Stringify returns a k8s selector string based on given map. Ex: "map[key] = value -> key=value"
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

// Mapify returns a map based on given string. Ex: "key=value -> map[key] = value"
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
