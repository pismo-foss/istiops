package cmd

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/pismo/istiops/pkg/client"
	"github.com/pismo/istiops/pkg/logger"
	istiOperator "github.com/pismo/istiops/pkg/operator"
	"github.com/pismo/istiops/pkg/router"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
)

var (
	trackingId string
	clients    *client.Set
)

func init() {
	kubeConfigDefaultPath := homedir.HomeDir() + "/.kube/config"
	rootCmd.PersistentFlags().String("context", "", "kube context (optional)")
	rootCmd.PersistentFlags().String("kubeconfig", kubeConfigDefaultPath, "config path (optional)")

	rootCmd.AddCommand(trafficCmd)
	rootCmd.AddCommand(versionCmd)
}

func clientSetup(kubeContext string, kubeConfigPath string) {
	var err error

	clients, err = client.New(kubeContext, kubeConfigPath)
	if err != nil {
		logger.Fatal(fmt.Sprintf("%s", err), "cmd")
	}
	logger.Debug(fmt.Sprintf("Initialized client from context '%s' and kubeConfig '%s'", kubeContext, kubeConfigPath), "cmd")

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

// Execute will execute the cmd client
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(fmt.Sprintf("%s", err), "cmd")
		os.Exit(-1)
	}
}
