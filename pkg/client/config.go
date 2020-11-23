package client

import (
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"github.com/pismo/istiops/pkg/router"
	"k8s.io/client-go/kubernetes"

	// in order to solve a gcp bug when trying to get the kubeconfig
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Set will define both kubernetes and istio interfaces
type Set struct {
	Kubernetes kubernetes.Interface
	Istio      router.IstioClientInterface
}

// ToRawKubeConfigLoader returns a ClientConfig with overrided attributes such as 'context'
func ToRawKubeConfigLoader(kubeContext string, kubeConfigPath string) clientcmd.ClientConfig {

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	// if you want to change the loading rules (which files in which order), you can do so here
	loadingRules.ExplicitPath = kubeConfigPath

	// if you want to change override values or bind them to flags, there are methods to help you
	configOverrides := &clientcmd.ConfigOverrides{
		ClusterDefaults: clientcmd.ClusterDefaults,
	}

	if kubeContext != "" {
		configOverrides.CurrentContext = kubeContext
	}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	return kubeConfig
}

// New will return a clientset with both kubernetes and istio ones
func New(kubeContext string, kubeConfigPath string) (*Set, error) {
	var istioClient router.IstioClientInterface
	var config *rest.Config
	var err error

	config, err = ToRawKubeConfigLoader(kubeContext, kubeConfigPath).ClientConfig()
	if err != nil {
		return &Set{}, err
	}

	// create both clientset
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &Set{}, err
	}

	istioClient, err = versioned.NewForConfig(config)
	if err != nil {
		return &Set{}, err
	}

	client := &Set{
		Kubernetes: kubeClient,
		Istio:      istioClient,
	}

	return client, nil
}
