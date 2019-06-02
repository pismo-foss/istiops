package main

import (
	"context"
	"github.com/pismo/istiops/services"
	_ "github.com/pkg/errors"
	_ "github.com/sirupsen/logrus"
	_ "github.com/snowzach/rotatefilehook"
	_ "gopkg.in/yaml.v2"
	_ "istio.io/api/networking/v1alpha3"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/apimachinery/pkg/runtime/schema"
	_ "k8s.io/client-go/dynamic"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	_ "k8s.io/client-go/tools/clientcmd"
)

func main() {
	apiStruct := services.ApiStruct{
		Name:      "api-pipelinetest",
		Namespace: "ext",
		Version:   "pipethree",
		Build:     "80"}

	services.DeployHelm(apiStruct, "cid-random", context.Background())
	services.K8sHealthCheck("cid-random", 5, apiStruct, context.Background())
}
