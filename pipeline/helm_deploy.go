package pipeline

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

func getConfigmapValues(api utils.ApiStruct, cid string, ctx context.Context) (*utils.ApiValues, error) {
	apiValues := &utils.ApiValues{}

	pwd, err := os.Getwd()
	if err != nil {
		utils.Error(fmt.Sprintf(err.Error()), cid)
		return apiValues, err
	}

	cmBytes, err := ioutil.ReadFile(pwd + "/" + api.Name + "/kubernetes/" + api.Namespace + "/" + api.Name + "-config.yaml")
	if err != nil {
		utils.Fatal(fmt.Sprintf(err.Error()), cid)
		return apiValues, err
	}
	cm := &v1core.ConfigMap{}
	if err := yaml.Unmarshal(cmBytes, cm); err != nil {
		utils.Fatal(fmt.Sprintf(err.Error()), cid)
		return apiValues, err
	}

	cm.Name = api.Name + "-config"
	cm.Namespace = api.Namespace
	utils.Info(fmt.Sprintf("Configmap extracted: %s", cm.Name), cid)

	valuesBytes, err := ioutil.ReadFile(pwd + "/" + api.Name + "/values.yaml")
	if err != nil {
		utils.Error(fmt.Sprintf(err.Error()), cid)
		return apiValues, err
	}

	if err := yaml.Unmarshal(valuesBytes, apiValues); err != nil {
		utils.Error(fmt.Sprintf(err.Error()), cid)
		return apiValues, err
	}
	utils.Info(fmt.Sprintf("Values extracted: %v", apiValues), cid)
	return apiValues, nil
}

func DeployHelm(api utils.ApiStruct, cid string, ctx context.Context) error {
	utils.Info(fmt.Sprintf("Applying configmap to %s environment...", api.Namespace), cid)
	apiValues, err := getConfigmapValues(api, cid, ctx)
	if err != nil {
		utils.Error(fmt.Sprintf(err.Error()), cid)
		return err
	}

	createDeployment(api, *apiValues, "random-cid", ctx)
	return nil
}

func createDeployment(api utils.ApiStruct, apiValues utils.ApiValues, cid string, ctx context.Context) error {
	api_fullname := fmt.Sprintf("%s-%s-%s-%s", api.Name, api.Namespace, api.Version, api.Build)
	if apiValues.Deployment.Image.DockerRegistry == "" {
		apiValues.Deployment.Image.DockerRegistry = "270036487593.dkr.ecr.us-east-1.amazonaws.com/"
	}
	deploymentsClient := kubernetesClient.AppsV1().Deployments(api.Namespace)

	// Getting dynamic protocol & ports
	containerPorts := []v1core.ContainerPort{}

	for portName, portValue := range apiValues.Deployment.Image.Ports {
		containerPorts = append(containerPorts, v1core.ContainerPort{
			Name:          portName,
			Protocol:      v1core.ProtocolTCP,
			ContainerPort: int32(portValue),
		})
	}

	deployment := &v1apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: api_fullname,
		},
		Spec: v1apps.DeploymentSpec{
			Replicas: int32Ptr(int32(apiValues.Deployment.Replicas[api.Namespace])),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":     api.Name,
					"build":   api.Build,
					"release": api_fullname,
					"version": api.Version,
				},
			},
			Template: v1core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"iam.amazonaws.com/role": apiValues.Deployment.Role,
					},
					Labels: map[string]string{
						"app":     api.Name,
						"build":   api.Build,
						"release": api_fullname,
						"version": api.Version,
					},
				},
				Spec: v1core.PodSpec{
					RestartPolicy: "Always",
					Containers: []v1core.Container{
						{
							Name:            api.Name,
							Image:           apiValues.Deployment.Image.DockerRegistry + api.Name + ":" + api.Version,
							Ports:           containerPorts,
							ImagePullPolicy: "Always",
							LivenessProbe: &v1core.Probe{
								Handler: v1core.Handler{
									Exec: &v1core.ExecAction{
										Command: []string{
											"curl",
											"-fsS",
											fmt.Sprintf("http://localhost:%d%s",
												apiValues.Deployment.Image.HealthCheck.HealthPort,
												apiValues.Deployment.Image.HealthCheck.LivenessProbeEndpoint,
											),
										},
									},
								},
								InitialDelaySeconds: 15,
								FailureThreshold:    3,
								PeriodSeconds:       30,
								SuccessThreshold:    1,
								TimeoutSeconds:      1,
							},
							ReadinessProbe: &v1core.Probe{
								Handler: v1core.Handler{
									Exec: &v1core.ExecAction{
										Command: []string{
											"curl",
											"-fsS",
											fmt.Sprintf(
												"http://localhost:%d%s",
												apiValues.Deployment.Image.HealthCheck.HealthPort,
												apiValues.Deployment.Image.HealthCheck.ReadinessProbeEndpoint,
											),
										},
									},
								},
								InitialDelaySeconds: 15,
								FailureThreshold:    3,
								PeriodSeconds:       30,
								SuccessThreshold:    1,
								TimeoutSeconds:      1,
							},
						},
					},
					NodeSelector: map[string]string{
						"kops.k8s.io/instancegroup": "nodes",
					},
				},
			},
		},
	}

	// Create Deployment
	utils.Info(fmt.Sprintf("Creating new deployment %s...", api_fullname), cid)
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		utils.Fatal(fmt.Sprintf(err.Error()), cid)
	}
	utils.Info(fmt.Sprintf("Created deployment %q!", result.GetObjectMeta().GetName()), cid)

	return nil
}

func int32Ptr(i int32) *int32 { return &i }
