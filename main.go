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
	clientSet, err := client.Set(homedir + "/.kube/config")
	if err != nil {
		utils.Fatal("Could not get clients", "cid")
	}

	resources := &operator.IstioResources{
		DestinationRuleName: "api-xpto-destinationrules",
		VirtualServiceName:  "api-xpto-virtualservices",
	}

	ips := operator.IstioOperator{
		TrackingId: "54ec4fd3-879b-404f-9812-c6b97f663b8d",
		Namespace:  "default",
		Client:     clientSet,
	}

	ips.Delete(resources)
	labels := operator.LabelSelector{Labels: map[string]string{"environment": "pipeline-go"}}
	err = ips.Clear(labels)
	if err != nil {
		fmt.Printf("")
	}

}
