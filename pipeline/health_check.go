package pipeline

import (
	"context"
	"errors"
	"fmt"
	"github.com/pismo/istiops/utils"
	v1core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func K8sHealthCheck(cid string, timeout time.Duration, api utils.ApiStruct, ctx context.Context) error {
	utils.Info("Starting kubernetes' healthcheck based in 'rollout' with a 180 seconds of timeout...", cid)

	pods, err := kubernetesClient.CoreV1().Pods(api.Namespace).List(v1.ListOptions{
		LabelSelector: "release=" + api.ApiFullname,
	})
	if err != nil {
		return err
	}

	watch, err := kubernetesClient.CoreV1().Pods(api.Namespace).Watch(v1.ListOptions{
		LabelSelector: "release=" + api.ApiFullname,
	})
	if err != nil {
		return err
	}


	c1 := make(chan bool, 1)
	i := 0
	podsSize := len(pods.Items)
	go func() {

		for event := range watch.ResultChan() {

			p, ok := event.Object.(*v1core.Pod)
			if !ok {
				utils.Fatal("unexpected type", cid)
			}

			utils.Info(fmt.Sprintf("Pod %s at status %s", p.Name, p.Status.Phase), cid)
			numberOfContainers := len(p.Status.ContainerStatuses)

			y := 0
			for _, containerStatus := range p.Status.ContainerStatuses {

				if containerStatus.Ready {
					y++
				}

				if y == numberOfContainers {
					i++
				}

			}

			if i == podsSize {
				c1 <- true
			}
		}

	}()


	select {
	case res := <-c1:
		fmt.Println(fmt.Sprintf("All containers are running! %v", res))
	case <-time.After(timeout * time.Second):
		return errors.New("TIMEOUT")
	}


	utils.Info("Application is running successfuly in pod!", cid)
	return nil
}
