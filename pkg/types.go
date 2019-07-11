package pkg

import (
	"github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	versionedclient "github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
)

func init() {
	//var err error
	homedir := homedir.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", homedir+"/.kube/config")

	// create the clientset
	kubernetesClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	istioClient, err = versionedclient.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Set global environment variables
	os.Setenv("SYSTEM", "Jenkins")
	os.Setenv("ENV", "dev")
}

type IstioValues struct {
	Name      string
	Version   string
	Build     int32
	Namespace string
}

type IstioResources struct {
	DestinationRule IstioMatchedDestinationRule
	VirtualService  IstioMatchedVirtualService
}

type IstioMatchedDestinationRule struct {
	Name string
	Item v1alpha3.DestinationRule
}

type IstioMatchedVirtualService struct {
	Subset string
	Item   v1alpha3.VirtualService
}

var (
	kubernetesClient *kubernetes.Clientset
	istioClient      versionedclient.Interface
	err              error
)
