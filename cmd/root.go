package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init(){
	var namespace string

	RootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes' cluster namespace")
	RootCmd.AddCommand(trafficCmd)

}

var RootCmd = &cobra.Command{
	Use: "istiops",
	Short: "Main",
	Long:  `Using Kubernetes' CRDs 'traffic' will manage it's rules`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Traffic shifting...")
	},
}
