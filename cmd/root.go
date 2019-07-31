package cmd

import (
	"fmt"
	"github.com/pismo/istiops/pkg"
	"github.com/spf13/cobra"
	"os"
)

var (
	istiops pkg.IstioOperationsInterface = pkg.IstioValues{"default"}
	VERSION string
	CID     string
)

func init() {
	var namespace string

	RootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes' cluster namespace")
	RootCmd.AddCommand(trafficCmd)
}

var RootCmd = &cobra.Command{
	Use:   "istiops",
	Short: "Main",
	Long:  `Istiops is a CLI library for Go that manages istio's traffic shifting easily.`,
}

func Execute(cid string, version string) {
	VERSION = version
	CID = cid

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
