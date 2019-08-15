package main

import (
	"fmt"
	"github.com/pismo/istiops/pkg/router"

	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
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

	dr := &router.DestinationRule{
		Istio: istioClient,
	}

	vs := &router.VirtualService{
		Istio: istioClient,
	}

	var op operator.Operator
	op = &operator.Istiops{
		TrackingId:      "54ec4fd3-879b-404f-9812-c6b97f663b8d",
		Name:            "api-xpto",
		Namespace:       "default",
		Build:           26,
		DestinationRuleRouter: dr,
		VirtualServiceRouter:  vs,
	}

	route := &router.Route{
		Port:     5000,
		Hostname: "api.domain.io",
		Selector: map[string]string{"environment": "pipeline-go"},
		Headers: map[string]string{
			"x-version": "PR-141",
			"x-cid":     "12312-123121-1212-1231-12131",
		},
		Weight: 0,
	}

	// Update a route
	err = op.Update(route)
	if err != nil {
		fmt.Printf("")
	}

}
