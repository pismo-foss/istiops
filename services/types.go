package services

import (
	versionedclient "github.com/pismo/istiops/pkg/client/clientset/versioned"
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

type ApiStruct struct {
	Name        string     `json:"name"`
	ApiFullname string     `json:"api_fullname"`
	Namespace   string     `json:"namespace"`
	Version     string     `json:"version"`
	Build       string     `json:"build"`
	HttpPort    uint32     `json:"http_port"`
	GrpcPort    uint32     `json:"grpc_port"`
	ApiValues   *ApiValues `json:"api_values"`
}

type ApiValues struct {
	Deployment   Deployment               `yaml:"deployment"`
	Resources    map[string]interface{}   `yaml:"resources"`
	NodeSelector map[string]interface{}   `yaml:"nodeSelector"`
	Tolerations  []map[string]interface{} `yaml:"tolerations"`
	Affinity     map[string]interface{}   `yaml:"affinity"`
}

type Deployment struct {
	Role     string           `yaml:"role"`
	Replicas map[string]int64 `yaml:"replicas"`
	Image    Image            `yaml:"image"`
}

type Image struct {
	HealthCheck    map[string]string `yaml:"healthCheck"`
	Ports          map[string]int64  `yaml:"ports"`
	DockerRegistry string            `yaml:"dockerRegistry"`
	PullPolicy     string            `yaml:"pullPolicy"`
}

var (
	kubernetesClient *kubernetes.Clientset
	istioClient      *versionedclient.Clientset
	PismoDomains     = map[string]string{"default": ".pismolabs.io"}
)
