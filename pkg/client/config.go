package client

import (
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"github.com/pismo/istiops/pkg/router"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Set struct {
	Kubernetes *kubernetes.Clientset
	Istio      router.IstioClientInterface
}

func New(kubeConfigPath string) (*Set, error) {
	var istioClient router.IstioClientInterface
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)

	// create the clientset
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	istioClient, err = versioned.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	client := &Set{
		Kubernetes: kubeClient,
		Istio:      istioClient,
	}

	return client, nil
}
