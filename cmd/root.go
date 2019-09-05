/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"github.com/google/uuid"
	istiOperator "github.com/pismo/istiops/pkg/operator"
	"github.com/pismo/istiops/pkg/router"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
)

var (
	trackingId  string
	istioClient router.IstioClientInterface
)

func init() {
	setup()
	rootCmd.AddCommand(trafficCmd)
	rootCmd.AddCommand(versionCmd)
}

func setup() {
	kubeConfigPath := homedir.HomeDir() + "/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	istioClient, err = versioned.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// generate random uuid
	tracking, err := uuid.NewUUID()
	if err != nil {
		panic(err.Error())
	}

	trackingId = tracking.String()
}

func operator(dr *router.DestinationRule, vs *router.VirtualService) istiOperator.Operator {
	op := &istiOperator.Istiops{
		DrRouter: dr,
		VsRouter: vs,
	}

	return op
}

var rootCmd = &cobra.Command{
	Use:   "istiops",
	Short: "Main",
	Long:  `Istiops is a CLI library for Go that manages istio's traffic shifting easily.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
