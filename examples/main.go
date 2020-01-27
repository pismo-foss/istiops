package examples

import (
	"fmt"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"github.com/pismo/istiops/pkg/logger"
	"github.com/pismo/istiops/pkg/operator"
	"github.com/pismo/istiops/pkg/router"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {

	var istioClient router.IstioClientInterface

	kubeConfigPath := homedir.HomeDir() + "/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	istioClient, err = versioned.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	var build uint32
	var trackingId string
	var metadataName string
	var metadataNamespace string

	trackingId = "54ec4fd3-879b-404f-9812-c6b97f663b8d"
	metadataName = "api-xpto"
	metadataNamespace = "default"
	build = 1

	var dr operator.Router
	dr = &router.DestinationRule{
		TrackingId: trackingId,
		Name:       metadataName,
		Namespace:  metadataNamespace,
		Build:      build,
		Istio:      istioClient,
	}

	var vs operator.Router
	vs = &router.VirtualService{
		TrackingId: trackingId,
		Name:       metadataName,
		Namespace:  metadataNamespace,
		Build:      build,
		Istio:      istioClient,
	}

	shift := router.Shift{
		Port:     5000,
		Hostname: "api.domain.io",
		Selector: map[string]string{
			"environment": "pipeline-go",
		},
		Traffic: router.Traffic{
			PodSelector: map[string]string{
				"app":     "api",
				"version": "1.3.2",
				"build":   "24",
			},
			RequestHeaders: map[string]string{
				"x-version":    "PR-141",
				"x-account-id": "233",
			},
			Weight: 0,
		},
	}

	var op operator.Operator
	op = &operator.Istiops{
		DrRouter: dr,
		VsRouter: vs,
	}

	// clear all routes + subsets
	err = op.Clear(shift, "soft")
	if err != nil {
		logger.Fatal(fmt.Sprintf("%s", err), trackingId)
	}

	// Update a route
	err = op.Update(shift)
	if err != nil {
		logger.Fatal(fmt.Sprintf("%s", err), trackingId)
	}

}
