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
	"github.com/pismo/istiops/pkg/operator"
	"github.com/pismo/istiops/pkg/router"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
)

var (
	op                  operator.Operator
	build               uint32
	trackingId          string
	metadataName        string
	namespace           string
	labelSelector       string
	mappedLabelSelector map[string]string
	shift               router.Shift
)

func init() {
	rootCmd.AddCommand(trafficCmd)
	rootCmd.AddCommand(versionCmd)
	setup()
}

func setup() {
	kubeConfigPath := homedir.HomeDir() + "/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	istioClient, err := versioned.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// generate random uuid
	uuid, err := uuid.NewUUID()
	if err != nil {
		panic(err.Error())
	}

	trackingId = uuid.String()

	ic := router.Client{
		Versioned: istioClient,
	}

	var dr operator.Router
	dr = &router.DestinationRule{
		TrackingId: trackingId,
		Name:       metadataName,
		Namespace:  namespace,
		Build:      build,
		Istio:      ic,
	}

	var vs operator.Router

	vs = &router.VirtualService{
		TrackingId: trackingId,
		Name:       metadataName,
		Namespace:  namespace,
		Build:      build,
		Istio:      ic,
	}

	shift = router.Shift{
		Port:     5000,
		Hostname: "api.domain.io",
		Selector: router.Selector{
			Labels: mappedLabelSelector,
		},
		Traffic: router.Traffic{
			PodSelector: map[string]string{
				"app":     "api",
				"version": "1.3.3",
				"build":   "24",
			},
			RequestHeaders: map[string]string{
				"x-version":    "PR-142",
				"x-account-id": "233",
			},
			Weight: 0,
		},
	}

	op = &operator.Istiops{
		DrRouter: dr,
		VsRouter: vs,
	}
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
