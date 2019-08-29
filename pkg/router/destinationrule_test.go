package router

import (
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var mockedRouter Router
var fClient Client

type MockedDestinationRule struct{
	TrackingId string
	Name string
	Build string
	Istio Client
}

func TestMain(m *testing.M) {
	// discard stdout logs if not being run with '-v' flag
	log.SetOutput(ioutil.Discard)
	TearUp()
	result := m.Run()
	os.Exit(result)

}

func TearUp() {
	fClient = Client{
		Fake: &fake.Clientset{},
	}
}

func TestCreate(t *testing.T) {
	successCases := []struct {
		dr     *DestinationRule
		shift  Shift
		want map[string]string
	}{
		{
			&DestinationRule{ Istio: fClient },
			Shift{
				Port:     5000,
				Hostname: "api.domain.io",
				Selector: &Selector{
					Labels: map[string]string{"environment": "pipeline-go"},
				},
				Traffic: &Traffic{
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
			},
			map[string]string{
				"app":     "api",
				"version": "1.3.2",
				"build":   "24",
			},
		},
	}

	for _, tt := range successCases {
		rule, _ := tt.dr.Create(&tt.shift)
		assert.Equal(t, rule.Subset.Labels, tt.want)
	}
}