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

// K8sHealthCheck checks if all the containers inside a pod of a given release are
// on a ready state with the given timeout.
func K8sHealthCheck(cid string, timeout time.Duration, releaseName string, namespace string, ctx context.Context) error {
	utils.Info("Starting kubernetes' healthcheck based in 'rollout' with a 180 seconds of timeout...", cid)

	pods, err := kubernetesClient.CoreV1().Pods(namespace).List(v1.ListOptions{
		LabelSelector: "release=" + releaseName,
	})
	if err != nil {
		return err
	}

	watch, err := kubernetesClient.CoreV1().Pods(namespace).Watch(v1.ListOptions{
		LabelSelector: "release=" + releaseName,
	})
	if err != nil {
		return err
	}


	c1 := make(chan bool, 1)
	podsSize := len(pods.Items)
	podsReady := map[string]bool{}

	for _, v := range pods.Items {
		utils.Info(fmt.Sprintf("Waiting pod %s to be ready...", v.ObjectMeta.Name), cid)
		podsReady[v.ObjectMeta.Name] = false
	}

	go func() {

		for event := range watch.ResultChan() {

			p, ok := event.Object.(*v1core.Pod)
			if !ok {
				utils.Fatal("unexpected type", cid)
			}

			//Checking if all the containers inside the pod are running
			numberOfContainers := len(p.Status.ContainerStatuses)
			y := 0
			for _, containerStatus := range p.Status.ContainerStatuses {

				if containerStatus.Ready {
					utils.Info(fmt.Sprintf("Container %s ready", containerStatus.Name), cid, utils.Fields{"pod": p.Name})
					y++
				}

				if y == numberOfContainers {
					utils.Info(fmt.Sprintf("All containers running for pod %s", p.ObjectMeta.Name), cid, utils.Fields{"pod": p.Name})
					podsReady[p.ObjectMeta.Name] = true
				}

			}

			//Check for number of pods already active
			j := 0
			for _, v := range podsReady {
				if v == true {
					j++
				}
			}

			//If number of pods active equals number of pods, then all pods are active! Leave loop.
			if j ==  podsSize {
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
