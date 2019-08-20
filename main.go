package main

import (
	"fmt"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"github.com/pismo/istiops/pkg/operator"
	"github.com/pismo/istiops/pkg/router"
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

	DrM := router.DrMetadata{
		TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
		Name:       "api-xpto",
		Namespace:  "default",
		Build:      29,
	}

	VsM := router.VsMetadata{
		TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
		Name:       "api-xpto",
		Namespace:  "default",
		Build:      29,
	}

	dr := &router.DestinationRule{
		Metadata: DrM,
		Istio:    istioClient,
	}

	vs := &router.VirtualService{
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
		Shift:    shift,
		DrRouter: dr,
		VsRouter: vs,
	}

	// Update a route
	err = op.Update(shift)

	if err != nil {
		fmt.Printf("")
	}

}
