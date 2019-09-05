package router

import (
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var fakeIstioClient IstioClientInterface

func TestMain(m *testing.M) {
	// discard stdout logs if not being run with '-v' flag
	log.SetOutput(ioutil.Discard)
	result := m.Run()
	os.Exit(result)
}

func TestUpdateDestinationRule_Integrated(t *testing.T) {
	fakeIstioClient = &fake.Clientset{}
	dr := DestinationRule{
		TrackingId: "unit-testing-tracking-id",
		Istio:      fakeIstioClient,
	}

	mockedDestinationRule := &v1alpha32.DestinationRule{}
	mockedDestinationRule.Name = "mocked-destination-rule"
	mockedDestinationRule.Namespace = "default"

	err := UpdateDestinationRule(&dr, mockedDestinationRule)
	assert.NoError(t, err)
}

func TestValidateDestinationRuleList_Unit(t *testing.T) {
	irl := IstioRouteList{
		VList: &v1alpha32.VirtualServiceList{
			Items: []v1alpha32.VirtualService{
				{},
			},
		},
		DList: &v1alpha32.DestinationRuleList{
			Items: []v1alpha32.DestinationRule{
				{},
			},
		},
	}

	err := ValidateDestinationRuleList(&irl)
	assert.NoError(t, err)
}

func TestValidateDestinationRuleList_Unit_EmptyItems(t *testing.T) {
	irl := IstioRouteList{
		VList: &v1alpha32.VirtualServiceList{
			Items: nil,
		},
		DList: &v1alpha32.DestinationRuleList{
			Items: nil,
		},
	}

	err := ValidateDestinationRuleList(&irl)
	assert.EqualError(t, err, "empty destinationRules")
}

func TestDestinationRule_Validate(t *testing.T) {
	fakeIstioClient = &fake.Clientset{}
	dr := DestinationRule{
		TrackingId: "unit-testing-tracking-id",
		Istio:      fakeIstioClient,
	}

	cases := []struct {
		shift Shift
		want        string
	}{
		{
			Shift{
				Port:     8080,
				Hostname: "api-domain",
				Selector: nil,
				Traffic:  Traffic{
					PodSelector: map[string]string{"version":"1.2.3"},
				},
			},
			"empty label-selector",
		},
		{
			Shift{
				Port:     0,
				Hostname: "api-domain",
				Selector: map[string]string{"app":"api-domain"},
				Traffic:  Traffic{
					PodSelector: map[string]string{"version":"1.2.3"},
				},
			},
			"empty port",
		},
		{
			Shift{
				Port:     1000,
				Hostname: "api-domain",
				Selector: map[string]string{"app":"api-domain"},
				Traffic:  Traffic{
					PodSelector: map[string]string{"version":"1.2.3"},
				},
			},
			"port not in range 1024 - 65535",
		},
		{
			Shift{
				Port:     66000,
				Hostname: "api-domain",
				Selector: map[string]string{"app":"api-domain"},
				Traffic:  Traffic{
					PodSelector: map[string]string{"version":"1.2.3"},
				},
			},
			"port not in range 1024 - 65535",
		},
	}

	for _, tt := range cases {
		err := dr.Validate(tt.shift)
		assert.EqualError(t, err, tt.want)
	}
}

func TestDestinationRule_Clear(t *testing.T) {
	fakeIstioClient = &fake.Clientset{}
	dr := DestinationRule{
		TrackingId: "unit-testing-tracking-id",
		Istio:      fakeIstioClient,
	}

	shift := Shift{}

	err := dr.Clear(shift)
	assert.NoError(t, err)
}
