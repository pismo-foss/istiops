package operator

import (
	"fmt"
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	versionedClientFake "github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/fake"

	"io/ioutil"
	"istio.io/api/networking/v1alpha3"
	"log"
	"os"
	"testing"

	"github.com/pismo/istiops/pkg/router"
)

var fClient *versionedClientFake.Clientset
var namespace string
var mockedDestinationRuleName string
var mockedVirtualServiceName string


type MockedResources struct {
	Metadata MockedMetadata
	Istio    *versionedClientFake.Clientset
}

func (MockedResources) Validate(s *router.Shift) error {
	panic("implement me")
}

func (MockedResources) Update(s *router.Shift) error {
	panic("implement me")
}

func (MockedResources) Delete(s *router.Shift) error {
	panic("implement me")
}

func (MockedResources) Clear(s *router.Shift) error {
	panic("implement me")
}

type MockedMetadata struct {
	TrackingId string
	Name       string
	Namespace  string
	Build      uint32
}

func TestMain(m *testing.M) {
	// discard stdout logs if not being run with '-v' flag
	log.SetOutput(ioutil.Discard)
	TearUp()
	result := m.Run()
	os.Exit(result)

}

func TearUp() {
	fClient = versionedClientFake.NewSimpleClientset()
	namespace = "test-namespace"
	mockedDestinationRuleName = "api-unit-test-destinationrule"
	mockedVirtualServiceName = "api-unit-test-virtualservice"

	var err error

	_, err = CreateMockedDestinationRule(
		mockedDestinationRuleName,
		namespace,
		map[string]string{
			"app":     "api-xpto",
			"version": "2.1.3",
		})
	if err != nil {
		log.Fatal("Could not create mocked DestinationRule")
	}

	_, err = CreateMockedVirtualService(mockedVirtualServiceName, namespace)
	if err != nil {
		log.Fatal("Could not create mocked VirtualService")
	}
}

func CreateMockedDestinationRule(resourceName string, resourceNamespace string, subsetLabels map[string]string) (mockedDestinationRule *v1alpha32.DestinationRule, error error) {
	mockedDr := &v1alpha32.DestinationRule{}
	mockedDr.Name = resourceName
	mockedDr.Namespace = resourceNamespace
	mockedDr.Labels = map[string]string{
		"app":     "api-xpto",
		"version": "2.1.3",
	}

	mockedDr.Spec.Subsets = append(mockedDr.Spec.Subsets, &v1alpha3.Subset{
		Name:   "subset-test",
		Labels: subsetLabels,
	})

	_, err := fClient.NetworkingV1alpha3().DestinationRules(resourceNamespace).Create(mockedDr)

	return mockedDr, err
}

func CreateMockedVirtualService(resourceName string, resourceNamespace string) (mockedVirtualService *v1alpha32.VirtualService, error error) {
	mockedVs := &v1alpha32.VirtualService{}
	mockedVs.Name = resourceName
	mockedVs.Namespace = resourceNamespace
	mockedVs.Labels = map[string]string{
		"app":     "api-xpto",
		"version": "2.1.3",
	}

	mockedVs.Spec.Hosts = []string{"api-unit-test.domain.io"}
	mockedVs.Spec.Gateways = []string{"unit-test-gateway"}

	mockedVs.Spec.Http = append(mockedVs.Spec.Http, &v1alpha3.HTTPRoute{})

	defaultMatch := &v1alpha3.HTTPMatchRequest{Uri: &v1alpha3.StringMatch{MatchType: &v1alpha3.StringMatch_Regex{Regex: ".+"}}}
	defaultDestination := &v1alpha3.HTTPRouteDestination{Destination: &v1alpha3.Destination{Host: "api-xpto", Subset: "subset-test", Port: &v1alpha3.PortSelector{Port: &v1alpha3.PortSelector_Number{Number: 8080}}}}
	defaultRoute := &v1alpha3.HTTPRoute{}
	defaultRoute.Match = append(defaultRoute.Match, defaultMatch)
	defaultRoute.Route = append(defaultRoute.Route, defaultDestination)

	mockedVs.Spec.Http = append(mockedVs.Spec.Http, defaultRoute)

	_, err := fClient.NetworkingV1alpha3().VirtualServices(resourceNamespace).Create(mockedVs)

	return mockedVs, err
}

func TestCreate(t *testing.T) {

	mockedTrackingId := "54ec4fd3-879b-404f-9812-c6b97f663b8d"
	mockedMetadataName := "api-xpto"
	mockedMetadataNamespace := "default"
	mockedBuild := uint32(35)

	m := MockedMetadata{
		TrackingId: mockedTrackingId,
		Name:       mockedMetadataName,
		Namespace:  mockedMetadataNamespace,
		Build:      mockedBuild,
	}

	var mockedDr router.Router
	var mockedVs router.Router

	mockedDr = &MockedResources{
		Metadata: m,
		Istio:    fClient,
	}

	mockedVs = &MockedResources{
		Metadata: m,
		Istio:    fClient,
	}

	shift := &router.Shift{
		Port:     5000,
		Hostname: "api.domain.io",
		Selector: &router.Selector{
			Labels: map[string]string{"environment": "pipeline-go"},
		},
		Traffic: &router.Traffic{
			PodSelector: map[string]string{
				"app":     "api",
				"version": "1.3.2",
				"build":   "24",
			},
			RequestHeaders: map[string]string{
				"x-version": "PR-141",
				"x-cid":     "12312-123121-1212-1231-12131",
			},
			Weight: 0,
		},
	}

	var op Operator
	op = &Istiops{
		Shift:    shift,
		DrRouter: mockedDr,
		VsRouter: mockedVs,
	}

	fmt.Println(op.Update(shift))
}
