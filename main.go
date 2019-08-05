package main

import (
	versionedclient "github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"github.com/pismo/istiops/pkg/client"
	_ "github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	//var err error
	homedir := homedir.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", homedir+"/.kube/config")

	// create the clientset
	kubernetesClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	istioClient, err := versionedclient.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	client := client.IstioClient{
		kubernetesClient,
		istioClient,
	}
}
