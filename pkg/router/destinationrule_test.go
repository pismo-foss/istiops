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
	ds := DestinationRule{
		TrackingId: "unit-testing-tracking-id",
		Name:       "unit-testing",
		Namespace:  "default",
		Build:      5,
		Istio:      fakeIstioClient,
	}

	mockedDestinationRule := &v1alpha32.DestinationRule{}
	mockedDestinationRule.Name = "mocked-destination-rule"
	mockedDestinationRule.Namespace = "default"

	err := UpdateDestinationRule(&ds, mockedDestinationRule)
	assert.Error(t, err)
}
