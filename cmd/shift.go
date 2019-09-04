package cmd

import (
	"fmt"
	"github.com/pismo/istiops/pkg/router"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

func init() {
	shiftCmd.PersistentFlags().StringP("namespace", "n", "default", "kubernetes' cluster namespace")
	shiftCmd.PersistentFlags().StringP("destination", "d", "", "* destination's hostname with port ('api.domain.io:8080' or 'k8s-service:8080')")
	shiftCmd.PersistentFlags().Uint32P("build", "b", 0, "* build")
	shiftCmd.PersistentFlags().StringP("label-selector", "l", "", "* labels selector to filter istio' resources")
	shiftCmd.PersistentFlags().StringP("headers", "H", "", "headers")
	shiftCmd.PersistentFlags().StringP("pod-selector", "p", "", "* pod")
	shiftCmd.PersistentFlags().Uint32P("weight", "w", 0, "* weight (percentage) of routing")

	_ = shiftCmd.MarkPersistentFlagRequired("destination")
	_ = shiftCmd.MarkPersistentFlagRequired("pod-selector")
	_ = shiftCmd.MarkPersistentFlagRequired("port")
	_ = shiftCmd.MarkPersistentFlagRequired("build")
}

var shiftCmd = &cobra.Command{
	Use:   "shift",
	Short: "Shift istio's traffic",
	Run: func(cmd *cobra.Command, args []string) {

		namespace := cmd.Flag("namespace").Value.String()
		if namespace == "" {
			namespace = "default"
		} else {
			namespace = cmd.Flag("namespace").Value.String()
		}

		destination := cmd.Flag("destination").Value.String()
		destinationSplitted := strings.Split(destination, ":")
		if len(destinationSplitted) != 2 {
			panic(fmt.Sprintf("destination '%s' does not follow the format 'destination:port'", destination))
		}

		var portUint uint64
		portUint, err := strconv.ParseUint(destinationSplitted[1], 10, 32)
		if err != nil {
			panic(err)
		}

		mappedLabelSelector, err := router.Mapify(trackingId, cmd.Flag("label-selector").Value.String())
		if err != nil {
			panic(err)
		}

		mappedPodSelector, err := router.Mapify(trackingId, cmd.Flag("pod-selector").Value.String())
		if err != nil {
			panic(err)
		}

		var headers map[string]string
		if cmd.Flag("headers").Value.String() == "" {
			headers = nil
		} else {
			headers, err = router.Mapify(trackingId, cmd.Flag("headers").Value.String())
			if err != nil {
				panic(err)
			}
		}

		var buildInt uint64
		if cmd.Flag("build").Value.String() != "" {
			buildInt, err = strconv.ParseUint(cmd.Flag("build").Value.String(), 10, 32)
			if err != nil {
				panic(err)
			}
		}

		var weightInt int64
		if cmd.Flag("weight").Value.String() == "" {
			weightInt = 0
		} else {
			weightInt, err = strconv.ParseInt(cmd.Flag("weight").Value.String(), 10, 32)
			if err != nil {
				panic(err)
			}
		}

		drR := router.DestinationRule{
			TrackingId: trackingId,
			Name:       destinationSplitted[0],
			Namespace:  namespace,
			Build:      uint32(buildInt),
			Istio:      client,
		}

		vsR := router.VirtualService{
			TrackingId: trackingId,
			Name:       destinationSplitted[0],
			Namespace:  namespace,
			Build:      uint32(buildInt),
			Istio:      client,
		}

		shift := router.Shift{
			Selector: mappedLabelSelector,
			Hostname: destinationSplitted[0],
			Port:     uint32(portUint),
			Traffic: router.Traffic{
				PodSelector:    mappedPodSelector,
				RequestHeaders: headers,
				Weight:         int32(weightInt),
			},
		}

		op := operator(&drR, &vsR)
		err = op.Update(shift)
		if err != nil {
			fmt.Println(err)
		}
	},
}
