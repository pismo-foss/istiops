package client

import (
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ClientSet struct {
	Kubernetes *kubernetes.Clientset
	Istio      *versioned.Clientset
}

func New(kubeConfigPath string) (*ClientSet, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)

	// create the clientset
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	istioClient, err := versioned.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	client := &ClientSet{
		Kubernetes: kubeClient,
		Istio:      istioClient,
	}

	return client, nil
}
