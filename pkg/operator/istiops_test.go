package operator

import (
	"errors"
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	versionedClientFake "github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
var mockedDr router.Router
var mockedVs router.Router

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

type MockedResources struct {
	TrackingId string
	Name       string
	Namespace  string
	Build      uint32
	Istio      *versionedClientFake.Clientset
}

func (m MockedResources) Create(s *router.Shift) (*router.IstioRules, error) {
	return &router.IstioRules{}, nil
}

func (m MockedResources) List(opts metav1.ListOptions) (*router.IstioRouteList, error) {
	return &router.IstioRouteList{}, nil
}

func (m MockedResources) Clear(s *router.Shift) error { return nil }

func (m MockedResources) Validate(s *router.Shift) error {
	return errors.New("a route needs to be served with a 'weight' or 'request headers', not both")
}

func (m MockedResources) Update(s *router.Shift) error { return nil }

func TestCreate(t *testing.T) {

	createErrorCases := []struct {
		fop  Operator
		s    *router.Shift
		want string
	}{
		{&Istiops{
			DrRouter: &MockedResources{
				TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
				Name:       "api-xpto",
				Namespace:  "default",
				Build:      2,
				Istio:      fClient,
			},
			VsRouter: &MockedResources{
				TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
				Name:       "api-xpto",
				Namespace:  "default",
				Build:      2,
				Istio:      fClient,
			},
		}, &router.Shift{
			Port:     5000,
			Hostname: "api.domain.io",
			Selector: &router.Selector{
				//Labels: map[string]string{"environment": "pipeline-go"},
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
		}, "label-selector must exists in need to find resources",
		},
		{&Istiops{
			DrRouter: &MockedResources{
				TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
				Name:       "api-xpto",
				Namespace:  "default",
				Build:      2,
				Istio:      fClient,
			},
			VsRouter: &MockedResources{
				TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
				Name:       "api-xpto",
				Namespace:  "default",
				Build:      2,
				Istio:      fClient,
			},
		},
			&router.Shift{
				Port:     5000,
				Hostname: "api.domain.io",
				Selector: &router.Selector{
					Labels: map[string]string{"environment": "pipeline-go"},
				},
				Traffic: &router.Traffic{
					//PodSelector: map[string]string{
					//	"app":     "api",
					//	"version": "1.3.2",
					//	"build":   "24",
					//},
					RequestHeaders: map[string]string{
						"x-version": "PR-141",
						"x-cid":     "12312-123121-1212-1231-12131",
					},
					Weight: 0,
				},
			}, "pod-selector must exists in need to find traffic destination",
		},
		{&Istiops{
			DrRouter: &MockedResources{
				TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
				Name:       "api-xpto",
				Namespace:  "default",
				Build:      2,
				Istio:      fClient,
			},
			VsRouter: &MockedResources{
				TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
				Name:       "api-xpto",
				Namespace:  "default",
				Build:      2,
				Istio:      fClient,
			},
		},
			&router.Shift{
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
					Weight: 10,
				},
			}, "a route needs to be served with a 'weight' or 'request headers', not both",
		},
	}

	for _, tCase := range createErrorCases {
		got := tCase.fop.Update(tCase.s)
		assert.EqualError(t, got, tCase.want)
	}


}

func TestUpdate(t *testing.T) {

	updateErrorCases := []struct {
		fop  Operator
		s    *router.Shift
		want string
	}{
		{&Istiops{
			DrRouter: &MockedResources{
				TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
				Name:       "api-xpto",
				Namespace:  "default",
				Build:      2,
				Istio:      fClient,
			},
			VsRouter: &MockedResources{
				TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
				Name:       "api-xpto",
				Namespace:  "default",
				Build:      2,
				Istio:      fClient,
			},
		}, &router.Shift{
			Port:     5000,
			Hostname: "api.domain.io",
			Selector: &router.Selector{
				//Labels: map[string]string{"environment": "pipeline-go"},
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
		}, "label-selector must exists in need to find resources",
		},
		{&Istiops{
			DrRouter: &MockedResources{
				TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
				Name:       "api-xpto",
				Namespace:  "default",
				Build:      2,
				Istio:      fClient,
			},
			VsRouter: &MockedResources{
				TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
				Name:       "api-xpto",
				Namespace:  "default",
				Build:      2,
				Istio:      fClient,
			},
		},
			&router.Shift{
				Port:     5000,
				Hostname: "api.domain.io",
				Selector: &router.Selector{
					Labels: map[string]string{"environment": "pipeline-go"},
				},
				Traffic: &router.Traffic{
					//PodSelector: map[string]string{
					//	"app":     "api",
					//	"version": "1.3.2",
					//	"build":   "24",
					//},
					RequestHeaders: map[string]string{
						"x-version": "PR-141",
						"x-cid":     "12312-123121-1212-1231-12131",
					},
					Weight: 0,
				},
			}, "pod-selector must exists in need to find traffic destination",
		},
		{&Istiops{
			DrRouter: &MockedResources{
				TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
				Name:       "api-xpto",
				Namespace:  "default",
				Build:      2,
				Istio:      fClient,
			},
			VsRouter: &MockedResources{
				TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
				Name:       "api-xpto",
				Namespace:  "default",
				Build:      2,
				Istio:      fClient,
			},
		},
			&router.Shift{
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
					Weight: 10,
				},
			}, "a route needs to be served with a 'weight' or 'request headers', not both",
		},
	}

	for _, tCase := range updateErrorCases {
		got := tCase.fop.Update(tCase.s)
		assert.EqualError(t, got, tCase.want)
	}

}

func TestClear(t *testing.T) {

	clearCases := []struct {
		fop  Operator
		s    *router.Shift
		want error
	}{
		{
			&Istiops{
				DrRouter: &MockedResources{
					TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
					Name:       "api-xpto",
					Namespace:  "default",
					Build:      2,
					Istio:      fClient,
				},
				VsRouter: &MockedResources{
					TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
					Name:       "api-xpto",
					Namespace:  "default",
					Build:      2,
					Istio:      fClient,
				},
			},
			&router.Shift{
				Port:     5000,
				Hostname: "api.domain.io",
				Selector: &router.Selector{
					//Labels: map[string]string{"environment": "pipeline-go"},
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
			}, nil,
		},
	}

	for _, tCase := range clearCases {
		got := tCase.fop.Clear(tCase.s)
		assert.NoError(t, got, tCase.want)
	}

}
