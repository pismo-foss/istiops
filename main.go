package main

import (
	"github.com/pismo/istiops/cmd"
	"github.com/pismo/istiops/utils"
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
	err := cmd.RootCmd.Execute()
	if err != nil {
		utils.Fatal("error", "cid")
	}
}

//func main() {
//	apiValues := utils.BuildApiValues("api-pipelinetest", "default", "1.0.0", "2210")
//	// pkg.CreateRouteResource(apiValues, "cid-random", context.Background())
//	// pipeline.DeployApi(apiValues, "cid-random", context.Background())
//	pipeline.IstioRouting(apiValues, "cid-random", context.Background())
//	// pipeline.K8sHealthCheck("cid-random", 5, apiValues, context.Background())
//}
