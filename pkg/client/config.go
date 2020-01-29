package client

import (
	"fmt"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"github.com/pismo/istiops/pkg/logger"
	"github.com/pismo/istiops/pkg/router"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Set struct {
	Kubernetes kubernetes.Interface
	Istio      router.IstioClientInterface
}

func New(kubeConfigPath string) (*Set, error) {
	var istioClient router.IstioClientInterface
	var config *rest.Config
	var err error
	config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while loading ./kube/config file %s", err.Error()), "bootstart")
		config, err = rest.InClusterConfig()
		if err != nil {
			logger.Error(fmt.Sprintf("Error while loading kubernetes conf (inside pod) %s", err.Error()), "bootstart")
		}
	}

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
