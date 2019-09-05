package operator

import (
	"github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	versionedClientFake "github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (m MockedResources) Create(shift router.Shift) (*router.IstioRules, error) {
	return &router.IstioRules{}, nil
}

func (m MockedResources) List(selector map[string]string) (*router.IstioRouteList, error) {
	return &router.IstioRouteList{
		// initialize both Lists with an empty item, to pass in cases with "if len(list) == 0"
		VList: &v1alpha3.VirtualServiceList{
			TypeMeta: v1.TypeMeta{},
			ListMeta: v1.ListMeta{},
			Items: []v1alpha3.VirtualService{
				{},
			},
		},
		DList: &v1alpha3.DestinationRuleList{
			TypeMeta: v1.TypeMeta{},
			ListMeta: v1.ListMeta{},
			Items: []v1alpha3.DestinationRule{
				{},
			},
		},
	}, nil
}

func (m MockedResources) Clear(shift router.Shift) error { return nil }

func (m MockedResources) Validate(shift router.Shift) error { return nil }

func (m MockedResources) Update(shift router.Shift) error { return nil }

// Tests scenarios for interface Istiops mocked

// It will test the Get() interface's method in the simplest scenario
func TestGet_Unit(t *testing.T) {

	var dr Router
	dr = &MockedResources{}

	var vs Router
	vs = &MockedResources{}

	var op Operator
	op = &Istiops{
		DrRouter: dr,
		VsRouter: vs,
	}

	irl, err := op.Get(map[string]string{})
	assert.Equal(t, router.IstioRouteList{
		VList: &v1alpha3.VirtualServiceList{
			TypeMeta: v1.TypeMeta{},
			ListMeta: v1.ListMeta{},
			Items: []v1alpha3.VirtualService{
				{},
			},
		},
		DList: &v1alpha3.DestinationRuleList{
			TypeMeta: v1.TypeMeta{},
			ListMeta: v1.ListMeta{},
			Items: []v1alpha3.DestinationRule{
				{},
			},
		},
	}, irl)
	assert.NoError(t, err)
}

// It will test the Clear() interface's method in the simplest scenario
func TestClear_Unit(t *testing.T) {

	var dr Router
	dr = &MockedResources{}

	var vs Router
	vs = &MockedResources{}

	shift := router.Shift{}

	var op Operator
	op = &Istiops{
		DrRouter: dr,
		VsRouter: vs,
	}

	err := op.Clear(shift)
	assert.NoError(t, err)
}

// It will test the Update() interface's method in the simplest scenario
func TestUpdate_Unit(t *testing.T) {

	var dr Router
	dr = &MockedResources{}

	var vs Router
	vs = &MockedResources{}

	shift := router.Shift{
		Selector: map[string]string{
			"app": "api-domain",
		},
		Traffic: router.Traffic{
			PodSelector: map[string]string{
				"version": "2.1.3",
			},
		},
	}

	var op Operator
	op = &Istiops{
		DrRouter: dr,
		VsRouter: vs,
	}

	err := op.Update(shift)
	assert.NoError(t, err)
}

// It will test the Update() interface's method in the scenario when pod-selector is empty
func TestUpdate_Unit_EmptyPodSelector(t *testing.T) {

	var dr Router
	dr = &MockedResources{}

	var vs Router
	vs = &MockedResources{}

	shift := router.Shift{
		Selector: map[string]string{
			"app": "api-domain",
		},
	}

	var op Operator
	op = &Istiops{
		DrRouter: dr,
		VsRouter: vs,
	}

	err := op.Update(shift)
	assert.EqualError(t, err, "pod-selector must exists in need to find traffic destination")
}

// It will test the Update() interface's method in the scenario when label-selector is empty
func TestUpdate_Unit_EmptyLabelSelector(t *testing.T) {

	var dr Router
	dr = &MockedResources{}

	var vs Router
	vs = &MockedResources{}

	shift := router.Shift{}

	var op Operator
	op = &Istiops{
		DrRouter: dr,
		VsRouter: vs,
	}

	err := op.Update(shift)
	assert.EqualError(t, err, "label-selector must exists in need to find resources")
}
