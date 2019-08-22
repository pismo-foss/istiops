package operator

import (
	"fmt"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"io/ioutil"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"testing"

	"github.com/pismo/istiops/pkg/router"
)

func TearUp() {

}

func TestMain(m *testing.M) {
	// discard stdout logs if not being run with '-v' flag
	log.SetOutput(ioutil.Discard)
	TearUp()
	result := m.Run()
	os.Exit(result)
}

func TestCreate(t *testing.T) {
	kubeConfigPath := homedir.HomeDir() + "/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	istioClient, err := versioned.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	mockedTrackingId := "54ec4fd3-879b-404f-9812-c6b97f663b8d"
	mockedMetadataName := "api-xpto"
	mockedMetadataNamespace := "default"
	mockedBuild := uint32(35)

	DrM := router.DrMetadata{
		TrackingId: mockedTrackingId,
		Name:       mockedMetadataName,
		Namespace:  mockedMetadataNamespace,
		Build:      mockedBuild,
	}

	VsM := router.VsMetadata{
		TrackingId: mockedTrackingId,
		Name:       mockedMetadataName,
		Namespace:  mockedMetadataNamespace,
		Build:      mockedBuild,
	}

	mockedDr := &router.DestinationRule{
		Metadata: DrM,
		Istio:    istioClient,
	}

	mockedVs := &router.VirtualService{
		Metadata: VsM,
		Istio:    istioClient,
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
