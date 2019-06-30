package main

import (
	"context"
	"github.com/pismo/istiops/pipeline"
	"github.com/pismo/istiops/pkg"
	"fmt"
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
	"strings"
)

func main() {
	apiStruct := pkg.ApiStruct{
		Name:      "api-statements",
		Namespace: "default",
		Version:   "bluegreeneb",
		Build:     "2210"}

	// some k8s resources does not allow special and uppercase characters
	replacer := strings.NewReplacer(
		".", "",
		"-", "",
		"/", "")
	simplifiedVersion := replacer.Replace(apiStruct.Version)
	simplifiedVersion = strings.ToLower(simplifiedVersion)


	// services.CreateRouteResource(apiStruct, "cid-random", context.Background())
	pipeline.DeployHelm(apiStruct, "cid-random", context.Background())
	// services.K8sHealthCheck("cid-random", 5, apiStruct, context.Background())
}
