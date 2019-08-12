package main

import (
	"github.com/pismo/istiops/pkg/client"
	"github.com/pismo/istiops/pkg/operator"
	"github.com/pismo/istiops/pkg/router"
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

	resources := &operator.IstioRoute{
		Port:     5000,
		Hostname: "api-xpto.domain.io",
		Selector: operator.Selector{
			Labels: map[string]string{
				"environment": "pipeline-go",
			},
		},
		Weight: &router.WeightShift{
			Headers: map[string]string{"x-version": "2.1.0"},
			Weight:  100,
		},
	}

	ips := operator.IstioOperator{
		TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
		Name:       "api-xpto",
		Namespace:  "default",
		Client:     clientSet,
	}

	// Create Resource
	//ips.Create(resources)

	// Update a route
	ips.Update(resources)

	// Delete resources example
	//ips.Delete(resources)

	// Clear rules example
	//labels := operator.Selector{
	//	Labels: map[string]string{"environment": "pipeline-go"},
	//}
	//err = ips.Clear(labels)
	//if err != nil {
	//	fmt.Printf("")
	//}

}
