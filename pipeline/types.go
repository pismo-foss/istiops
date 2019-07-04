package pipeline

import (
	"os"

	versionedclient "github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	WARN_NO_REGISTRY_FOUND                               = "Unable to get base docker registry. Check arguments."
	WARN_NO_PORT_SPECIFIED                               = "No grpc or http port specified!"
	WARN_NO_NECESSARY_NAMES_SPECIFIED                    = "No Name/namespace/version/build specified, check arguments!"
	WARN_NO_HEALTHCHECK_OR_READINESS_ENDPOINT_CONFIGURED = "Http port or Liveness or Readiness endpoint not informed."
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
