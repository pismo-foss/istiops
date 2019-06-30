package pipeline

import (
	versionedclient "github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
)

func init() {
	var err error
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
	println("SYSTEM:", os.Getenv("SYSTEM"))

}

var (
	kubernetesClient *kubernetes.Clientset
	istioClient      *versionedclient.Clientset
	PismoDomains     = map[string]string{"ext": ".pismolabs.io", "prod": ".pismo.io", "itau": ".pismo.cloud", "default": ".pismolabs.io"}
)
