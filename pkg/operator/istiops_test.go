package operator

import (
	"errors"
	versionedClientFake "github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/pismo/istiops/pkg/router"
)

var fClient *versionedClientFake.Clientset

func TestMain(m *testing.M) {
	// discard stdout logs if not being run with '-v' flag
	log.SetOutput(ioutil.Discard)
	TearUp()
	result := m.Run()
	os.Exit(result)

}

func TearUp() {
	fClient = versionedClientFake.NewSimpleClientset()
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
