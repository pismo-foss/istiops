package pipeline

import (
	"context"
	"fmt"
	"github.com/pismo/istiops/pkg"
	"github.com/pismo/istiops/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"time"
)

func K8sHealthCheck(cid string, timeout int, api pkg.ApiStruct, ctx context.Context) error {
	api_fullname := fmt.Sprintf("%s-%s-%s-%s", api.Name, api.Namespace, api.Version, api.Build)
	utils.Info("Starting kubernetes' healthcheck based in 'rollout' with a 180 seconds of timeout...", cid)

	dpl, err := kubernetesClient.AppsV1().Deployments(api.Namespace).List(v1.ListOptions{})
	if err != nil {
		utils.Fatal(err.Error(), cid)
		return err
	}

	for _, v := range dpl.Items {
		if strings.Contains(v.Name, api_fullname) {
			// -> todo : timeout implementation
			pods, err := kubernetesClient.CoreV1().Pods(api.Namespace).List(v1.ListOptions{})
			if err != nil {
				utils.Fatal(err.Error(), cid)
				return err
			}

			// validate pod statuses per container
			for _, p := range pods.Items {
				if strings.Contains(p.Name, api_fullname) {
					utils.Info(fmt.Sprintf("Waiting for a healthy status for pod: '%s'...", p.Name), cid)
					sum := 0
					for _, cst := range p.Status.ContainerStatuses {
						// waiting for container to be healthy
						for {
							sum++
							if cst.State.Running != nil {
								utils.Debug(fmt.Sprintf("container '%s' for  pod is healthy! %s...", cst.Name, p.Name), cid)
								break
							}
							time.Sleep(5 * time.Second)
							// if the container is not in status 'Running' after 5 minutes, terminate health check with exit 1
							if sum >= 60 {
								utils.Fatal(fmt.Sprintf("Container '%s' from pod '%s' had the validation time expired due to an unknown failure. Check the pod's logs for additional details", cst.Name, p.Name), cid)
							}
						}
					}
				}
			}
		}
	}
	utils.Info("Application is running successfuly in pod!", cid)
	return nil
}
