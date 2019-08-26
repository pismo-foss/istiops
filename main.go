package main

import (
	"fmt"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"github.com/pismo/istiops/pkg/operator"
	"github.com/pismo/istiops/pkg/router"
	"github.com/pismo/istiops/utils"
	_ "github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {

	kubeConfigPath := homedir.HomeDir() + "/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	istioClient, err := versioned.NewForConfig(config)
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
	build = 35

	DrM := router.DrMetadata{
		TrackingId: trackingId,
		Name:       metadataName,
		Namespace:  metadataNamespace,
		Build:      build,
	}

	VsM := router.VsMetadata{
		TrackingId: trackingId,
		Name:       metadataName,
		Namespace:  metadataNamespace,
		Build:      build,
	}

	var dr router.Router
	dr = &router.DestinationRule{
		Metadata: DrM,
		Istio:    istioClient,
	}

	var vs router.Router
	vs = &router.VirtualService{
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

	var op operator.Operator
	op = &operator.Istiops{
		DrRouter: dr,
		VsRouter: vs,
	}

	// clear all routes + subsets
	err = op.Clear(shift)
	if err != nil {
		utils.Fatal(fmt.Sprintf("%s", err), trackingId)
	}

	//// Update a route
	err = op.Update(shift)
	if err != nil {
		utils.Fatal(fmt.Sprintf("%s", err), trackingId)
	}

}
