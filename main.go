package main

import (
	"fmt"

	"github.com/heptio/contour/apis/generated/clientset/versioned"
	"github.com/pismo/istiops/pkg/operator"
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

	dr = &DestinationRuleMock{
		Istio: nil,
	}

	vs = &VirtualService{
		Istio: istioClient,
	}

	var op Operator
	op = &Istiops{
		TrackingId:      "54ec4fd3-879b-404f-9812-c6b97f663b8d",
		Name:            "api-xpto",
		Namespace:       "default",
		Build:           26,
		DestinationRule: dr,
		VirtualService:  vs,
	}

	route := &operator.Route{
		Port:     5000,
		Hostname: "api-xpto.domain.io",
		Selector: operator.Selector{
			Labels: map[string]string{"environment": "pipeline-go"},
		},
		Headers: map[string]string{
			"x-version": "PR-141",
			"x-cid":     "blau",
		},
		Weight: 0,
	}

	// Update a route
	err = op.Update(route)
	if err != nil {
		fmt.Printf("")
	}

}
