package main

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
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

var (
	// VERSION is set during build
	VERSION = "0.0.1"
)

func main() {
	uuid, err := uuid.NewV4()
	if err != nil {
		utils.Fatal("Could not generate CID", "")
	}
	cid := fmt.Sprintf("%v", uuid)

	cmd.Execute(cid, VERSION)
}

//func main() {
//	apiValues := utils.BuildApiValues("api-pipelinetest", "default", "1.0.0", "2210")
//	// pkg.CreateRouteResource(apiValues, "cid-random", context.Background())
//	// pipeline.DeployApi(apiValues, "cid-random", context.Background())
//	pipeline.IstioRouting(apiValues, "cid-random", context.Background())
//	// pipeline.K8sHealthCheck("cid-random", 5, apiValues, context.Background())
//}
