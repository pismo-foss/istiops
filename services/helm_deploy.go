package services

import (
	"context"
	"fmt"
	"github.com/pismo/istiops/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	v1apps "k8s.io/api/apps/v1"
	v1core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

func DeployHelm(api ApiStruct, cid string, ctx context.Context) error {
	utils.Info(fmt.Sprintf("Applying configmap to %s environment...", api.Namespace), cid)

	pwd, err := os.Getwd()
	if err != nil {
		utils.Error(fmt.Sprintf(err.Error()), cid)
		return err
	}

	cmBytes, err := ioutil.ReadFile(pwd + "/" + api.Name + "/kubernetes/" + api.Namespace + "/" + api.Name + "-config.yaml")
	if err != nil {
		utils.Fatal(fmt.Sprintf(err.Error()), cid)
		return err
	}
	cm := &v1core.ConfigMap{}
	if err := yaml.Unmarshal(cmBytes, cm); err != nil {
		utils.Fatal(fmt.Sprintf(err.Error()), cid)
		return err
	}

	cm.Name = api.Name + "-config"
	cm.Namespace = api.Namespace
	utils.Info(fmt.Sprintf("Configmap extracted: %s", cm.Name), cid)

	valuesBytes, err := ioutil.ReadFile(pwd + "/" + api.Name + "/values.yaml")
	if err != nil {
		utils.Error(fmt.Sprintf(err.Error()), cid)
		return err
	}

	api.ApiValues = &ApiValues{}
	if err := yaml.Unmarshal(valuesBytes, api.ApiValues); err != nil {
		utils.Error(fmt.Sprintf(err.Error()), cid)
		return err
	}
	utils.Info(fmt.Sprintf("Values extracted: %v", api.ApiValues), cid)

	if err := createDeployment(api, "random-cid", ctx); err != nil {
		return err
	}
	return nil
}

func createDeployment(api ApiStruct, cid string, ctx context.Context) error {
	api.ApiValues.Deployment.Image.DockerRegistry = "270036487593.dkr.ecr.us-east-1.amazonaws.com/"
	deploymentsClient := kubernetesClient.AppsV1().Deployments(api.Namespace)

	// Getting dynamic protocol & ports
	containerPorts := []v1core.ContainerPort{}
	for portName, portValue := range api.ApiValues.Deployment.Image.Ports {
		containerPorts = append(containerPorts, v1core.ContainerPort{
			Name:          portName,
			Protocol:      v1core.ProtocolTCP,
			ContainerPort: int32(portValue),
		})
	}

	fmt.Println(api.ApiValues.Deployment.Image.DockerRegistry + api.Name + ":" + api.Version)
	deployment := &v1apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: api.ApiFullname,
		},
		Spec: v1apps.DeploymentSpec{
			Replicas: int32Ptr(int32(api.ApiValues.Deployment.Replicas[api.Namespace])),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":     api.Name,
					"build":   api.Build,
					"release": api.ApiFullname,
					"version": api.Version,
				},
			},
			Template: v1core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":     api.Name,
						"build":   api.Build,
						"release": api.ApiFullname,
						"version": api.Version,
					},
				},
				Spec: v1core.PodSpec{
					Containers: []v1core.Container{
						{
							Name:  api.Name,
							Image: api.ApiValues.Deployment.Image.DockerRegistry + api.Name + ":" + api.Version,
							Ports: containerPorts,
						},
					},
				},
			},
		},
	}

	// Create Deployment
	utils.Info(fmt.Sprintf("Creating new deployment %s...", api.ApiFullname), cid)
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		utils.Fatal(fmt.Sprintf(err.Error()), cid)
	}
	utils.Info(fmt.Sprintf("Created deployment %q!", result.GetObjectMeta().GetName()), cid)

	return nil
}

func int32Ptr(i int32) *int32 { return &i }
