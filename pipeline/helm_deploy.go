package pipeline

import (
	"context"
	"errors"
	"fmt"
	"github.com/pismo/istiops/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	v1apps "k8s.io/api/apps/v1"
	v1core "k8s.io/api/core/v1"
	metaErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"os"
	"time"
)

func DeployApi(api utils.ApiStruct, cid string, ctx context.Context) error {
	utils.Info(fmt.Sprintf("Creating deploy in %s environment...", api.Namespace), cid)

	apiValues, err := getApiValues(api, cid, ctx)
	if err != nil {
		utils.Error(fmt.Sprintf(err.Error()), cid)
		return err
	}
	api.ApiValues = apiValues

	if err := createService(api, cid, ctx); err != nil {
		return err
	}

	if err := createConfig(api, cid, ctx); err != nil {
		return err
	}

	if err := createDeployment(api, cid, ctx); err != nil {
		return err
	}

	time.Sleep(2  * time.Second)
	if err := K8sHealthCheck(cid, 180, api, ctx); err != nil {
		return err
	}

	return nil
}

func createConfig(api utils.ApiStruct, cid string, ctx context.Context) error {

	configmapClient := kubernetesClient.CoreV1().ConfigMaps(api.Namespace)
	utils.Info(fmt.Sprintf("Creating configmap in %s environment...", api.Namespace), cid)

	cm, err := getConfigmapValues(api, cid, ctx)
	if err != nil {
		return err
	}

	_, err = configmapClient.Get(cm.Name, metav1.GetOptions{})
	if err != nil {
		customErr, ok := err.(*metaErrors.StatusError)
		if !ok {
			return err
		}

		if customErr.Status().Code != 404 {
			return err
		}


		utils.Info(fmt.Sprintf("Applying configmap: %s", cm.Name), cid)
		_, err = configmapClient.Create(&cm)
		if err != nil {
			return err
		}

		return nil
	}

	utils.Info(fmt.Sprintf("Configmap %s already exists, patching it.", cm.Name), cid)
	_, err = configmapClient.Update(&cm)
	if err != nil {
		return err
	}
	return nil
}

func createService(api utils.ApiStruct, cid string, ctx context.Context) error {

	k8sClientService := kubernetesClient.CoreV1().Services(api.Namespace)
	utils.Info(fmt.Sprintf("Creating services in %s environment...", api.Namespace), cid)


	service, err := k8sClientService.Get(api.Name, metav1.GetOptions{})
	if err != nil {
		customErr, ok := err.(*metaErrors.StatusError)
		if !ok {
			return err
		}

		if customErr.Status().Code != 404 {
			return err
		}

		service.Name = api.Name
		service.Namespace = api.Namespace
		service.Spec.Type = v1core.ServiceTypeClusterIP
		service.Spec.Selector = map[string]string{"app": api.Name}
		service.Spec.Ports = make([]v1core.ServicePort, 0)

		port := v1core.ServicePort{}
		port.Name = "http-" + api.Name
		port.Port = int32(api.HttpPort)
		port.Protocol = "TCP"
		port.TargetPort = intstr.IntOrString{Type:intstr.String, IntVal: int32(api.HttpPort)}

		service.Spec.Ports = append(service.Spec.Ports, port)

		utils.Info(fmt.Sprintf("Applying configmap: %s", service.Name), cid)
		_, err = k8sClientService.Create(service)
		if err != nil {
			return err
		}

		return nil
	}

	utils.Info(fmt.Sprintf("Service %s already exists, patching it.", service.Name), cid)
	_, err = k8sClientService.Update(service)
	if err != nil {
		return err
	}
	return nil
}

func createDeployment(api utils.ApiStruct, cid string, ctx context.Context) error {

	apiValues := api.ApiValues
	if err := validateCreateDeploymentArgs(api); err != nil {
		return err
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
			Name: api.ApiFullname,
		},
		Spec: v1apps.DeploymentSpec{
			Replicas: int32Ptr(int32(apiValues.Deployment.Replicas[api.Namespace])),
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
					Annotations: map[string]string{
						"iam.amazonaws.com/role": apiValues.Deployment.Role,
					},
					Labels: map[string]string{
						"app":     api.Name,
						"build":   api.Build,
						"release": api.ApiFullname,
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
								InitialDelaySeconds: 10,
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
	utils.Info(fmt.Sprintf("Creating new deployment %s...", api.ApiFullname), cid)
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		return err
	}
	utils.Info(fmt.Sprintf("Created deployment %q!", result.GetObjectMeta().GetName()), cid)

	return nil
}


func getConfigmapValues(api utils.ApiStruct, cid string, ctx context.Context) (v1core.ConfigMap, error) {
	cm := v1core.ConfigMap{}
	pwd, err := os.Getwd()
	if err != nil {
		utils.Error(fmt.Sprintf(err.Error()), cid)
		return cm, err
	}

	cmBytes, err := ioutil.ReadFile(pwd + "/" + api.Name + "/kubernetes/" + api.Namespace + "/" + api.Name + "-config.yaml")
	if err != nil {
		utils.Fatal(fmt.Sprintf(err.Error()), cid)
		return cm, err
	}
	if err := yaml.Unmarshal(cmBytes, &cm); err != nil {
		utils.Fatal(fmt.Sprintf(err.Error()), cid)
		return cm, err
	}

	cm.Name = api.Name + "-config"
	cm.Namespace = api.Namespace
	utils.Info(fmt.Sprintf("Configmap extracted: %s", cm.Name), cid)
	return cm, nil
}

func getApiValues(api utils.ApiStruct, cid string, ctx context.Context) (utils.ApiValues, error) {
	apiValues := utils.ApiValues{}

	pwd, err := os.Getwd()
	if err != nil {
		utils.Error(fmt.Sprintf(err.Error()), cid)
		return apiValues, err
	}

	valuesBytes, err := ioutil.ReadFile(pwd + "/" + api.Name + "/values.yaml")
	if err != nil {
		utils.Error(fmt.Sprintf(err.Error()), cid)
		return apiValues, err
	}

	if err := yaml.Unmarshal(valuesBytes, &apiValues); err != nil {
		utils.Error(fmt.Sprintf(err.Error()), cid)
		return apiValues, err
	}

	if apiValues.Deployment.Image.HealthCheck.Enabled && (apiValues.Deployment.Image.HealthCheck.LivenessProbeEndpoint == "" || apiValues.Deployment.Image.HealthCheck.ReadinessProbeEndpoint  == "") {
		apiValues.Deployment.Image.HealthCheck.ReadinessProbeEndpoint = "/health"
		apiValues.Deployment.Image.HealthCheck.LivenessProbeEndpoint = "/health"
		apiValues.Deployment.Image.HealthCheck.HealthPort = apiValues.Deployment.Image.Ports["http"]
	}

	utils.Info(fmt.Sprintf("Values extracted: %v", apiValues), cid)
	return apiValues, nil
}


func validateCreateDeploymentArgs(apiStruct utils.ApiStruct) error {

	if apiStruct.HttpPort <= 0 && apiStruct.GrpcPort <= 0 {
		return errors.New(WARN_NO_PORT_SPECIFIED)
	}

	if apiStruct.Name == "" || apiStruct.Namespace == "" || apiStruct.Version == "" || apiStruct.Build == "" {
		return errors.New(WARN_NO_NECESSARY_NAMES_SPECIFIED)
	}

	if apiStruct.ApiValues.Deployment.Image.DockerRegistry == "" {
		return errors.New(WARN_NO_REGISTRY_FOUND)
	}

	deployment := apiStruct.ApiValues.Deployment
	if deployment.Image.HealthCheck.Enabled && (deployment.Image.HealthCheck.HealthPort <= 0 || deployment.Image.HealthCheck.LivenessProbeEndpoint == "" || deployment.Image.HealthCheck.ReadinessProbeEndpoint == "") {
		return errors.New(WARN_NO_HEALTHCHECK_OR_READINESS_ENDPOINT_CONFIGURED)
	}


	return nil
}

func int32Ptr(i int32) *int32 { return &i }
