package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pismo/istiops/pkg/client"
	"github.com/pismo/istiops/pkg/logger"
	istiOperator "github.com/pismo/istiops/pkg/operator"
	"github.com/pismo/istiops/pkg/router"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
	"os"
)

var (
	trackingId string
	clients    *client.Set
)

func init() {
	setup()
	rootCmd.AddCommand(trafficCmd)
	rootCmd.AddCommand(versionCmd)
}

func setup() {
	kubeConfigPath := homedir.HomeDir() + "/.kube/config"
	var err error
	clients, err = client.New(kubeConfigPath)
	if err != nil {
		logger.Fatal(fmt.Sprintf("%s", err), "cmd")
	}

	// generate random uuid
	tracking, err := uuid.NewUUID()
	if err != nil {
		logger.Fatal(fmt.Sprintf("%s", err), "cmd")
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
	Long: `
 _     _   _
(_)___| |_(_) ___  _ __  ___
| / __| __| |/ _ \| '_ \/ __|
| \__ \ |_| | (_) | |_) \__ \
|_|___/\__|_|\___/| .__/|___/
                  |_|

Istiops is a CLI library for Go that manages istio's traffic shifting easily.
	`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(fmt.Sprintf("%s", err), "cmd")
		os.Exit(-1)
	}
}
