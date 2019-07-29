//// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
////
//// Licensed under the Apache License, Version 2.0 (the "License");
//// you may not use this file except in compliance with the License.
//// You may obtain a copy of the License at
////
////     http://www.apache.org/licenses/LICENSE-2.0
////
//// Unless required by applicable law or agreed to in writing, software
//// distributed under the License is distributed on an "AS IS" BASIS,
//// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//// See the License for the specific language governing permissions and
//// limitations under the License.
//
package cmd

//
//import (
//	"bufio"
//	"fmt"
//	"io"
//	"os"
//	"strings"
//	"time"
//
//	"github.com/spf13/cobra"
//
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	_ "k8s.io/client-go/plugin/pkg/client/auth"
//
//	"istio.io/api/networking/v1alpha3"
//	"istio.io/istio/pilot/pkg/config/kube/crd"
//	"istio.io/istio/pilot/pkg/model"
//)
//
//var (
//	kubeconfig    string
//	configContext string
//)
//
//// getCmd represents the get command
//var getCmd = &cobra.Command{
//	Use:   "get",
//	Short: "A brief description of your command",
//	Long: `A longer description that spans multiple lines and likely contains examples
//and usage of using your command. For example:
//
//Cobra is a CLI library for Go that empowers applications.
//This application is a tool to generate the needed files
//to quickly create a Cobra application.`,
//	Run: func(cmd *cobra.Command, args []string) {
//		fmt.Println("get called")
//
//		stdout := bufio.NewWriter(os.Stdout)
//		defer stdout.Flush()
//
//		configClient, err := newClient()
//		if err != nil {
//			panic(err)
//		}
//
//		typ, err := protoSchema(configClient, "virtualservice")
//		if err != nil {
//			panic(err)
//		}
//
//		var typs []model.ProtoSchema
//		typs = []model.ProtoSchema{typ}
//
//		//var ns string
//		//ns = v1.NamespaceAll
//
//		config, exists := configClient.Get(typs[0].Type,
//			"api-statements-virtualservice",
//			"qa")
//
//		if exists {
//			printShortVirtualService(config, stdout)
//		}
//	},
//}
//
//func init() {
//	rootCmd.AddCommand(getCmd)
//
//	// Here you will define your flags and configuration settings.
//
//	// Cobra supports Persistent Flags which will work for this command
//	// and all subcommands, e.g.:
//	// getCmd.PersistentFlags().String("foo", "", "A help for foo")
//
//	// Cobra supports local flags which will only run when this command
//	// is called directly, e.g.:
//	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
//}
//
//func newClient() (model.ConfigStore, error) {
//	return crd.NewClient(kubeconfig, configContext, model.IstioConfigTypes, "")
//}
//
//// The protoSchema is based on the kind (for example "virtualservice" or "destinationrule")
//func protoSchema(configClient model.ConfigStore, typ string) (model.ProtoSchema, error) {
//	for _, desc := range configClient.ConfigDescriptor() {
//		switch strings.ToLower(typ) {
//		case crd.ResourceName(desc.Type), crd.ResourceName(desc.Plural):
//			return desc, nil
//		case desc.Type, desc.Plural: // legacy hyphenated resources names
//			return model.ProtoSchema{}, fmt.Errorf("%q not recognized. Please use non-hyphenated resource name %q",
//				typ, crd.ResourceName(typ))
//		}
//	}
//	return model.ProtoSchema{}, fmt.Errorf("configuration type %s not found, the types are %v",
//		typ, strings.Join(supportedTypes(configClient), ", "))
//}
//
//func supportedTypes(configClient model.ConfigStore) []string {
//	types := configClient.ConfigDescriptor().Types()
//	for i := range types {
//		types[i] = crd.ResourceName(types[i])
//	}
//	return types
//}
//
//func printShortVirtualService(config *model.Config, w io.Writer) {
//	virtualService, ok := config.Spec.(*v1alpha3.VirtualService)
//	if !ok {
//		fmt.Fprintf(w, "Not a virtualservice: %v", config)
//		return
//	}
//	fmt.Fprintf(w, "%s\t%s\t%s\t%5d\t%4d\t%s\t%s\n",
//		config.Name,
//		strings.Join(virtualService.Gateways, ","),
//		strings.Join(virtualService.Hosts, ","),
//		len(virtualService.Http),
//		len(virtualService.Tcp),
//		config.Namespace,
//		renderTimestamp(config.CreationTimestamp))
//}
//
////renderTimestamp creates a human-readable age similar to docker and kubectl CLI output
//func renderTimestamp(ts metav1.Time) string {
//	if ts.IsZero() {
//		return "<unknown>"
//	}
//
//	seconds := int(time.Since(ts.Time).Seconds())
//	if seconds < -2 {
//		return fmt.Sprintf("<invalid>")
//	} else if seconds < 0 {
//		return fmt.Sprintf("0s")
//	} else if seconds < 60 {
//		return fmt.Sprintf("%ds", seconds)
//	}
//
//	minutes := int(time.Since(ts.Time).Minutes())
//	if minutes < 60 {
//		return fmt.Sprintf("%dm", minutes)
//	}
//
//	hours := int(time.Since(ts.Time).Hours())
//	if hours < 24 {
//		return fmt.Sprintf("%dh", hours)
//	} else if hours < 365*24 {
//		return fmt.Sprintf("%dd", hours/24)
//	}
//	return fmt.Sprintf("%dy", int((hours/24)/365))
//}
