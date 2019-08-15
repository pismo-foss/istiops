package main

import (
	"fmt"
	"github.com/pismo/istiops/pkg/client"
	"github.com/pismo/istiops/pkg/operator"
	"github.com/pismo/istiops/utils"
	_ "github.com/pkg/errors"
	"k8s.io/client-go/util/homedir"
)

func main() {
	homedir := homedir.HomeDir()
	clientSet, err := client.Add(homedir + "/.kube/config")
	if err != nil {
		utils.Fatal("Could not get clients", "cid")
	}

	var ips operator.Istiops
	ips = &operator.IstioOperator{
		TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
		Name:       "api-xpto",
		Namespace:  "default",
		Build:      26,
		Client:     clientSet,
	}

	routeResource := &operator.IstioRoute{
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
	err = ips.Update(routeResource)
	if err != nil {
		fmt.Printf("")
	}

}
