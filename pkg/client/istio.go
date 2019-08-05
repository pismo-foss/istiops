package client

import (
	versionedclient "github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
)

type IstioClient struct {
	kubernetesClient *kubernetes.Interface
	istioClient      *versionedclient.Interface
}
