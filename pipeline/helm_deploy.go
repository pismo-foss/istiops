package pipeline

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/pismo/istiops/utils"
	"gopkg.in/yaml.v2"
	v1apps "k8s.io/api/apps/v1"
	v1core "k8s.io/api/core/v1"
	metaErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// DeployapiValues deploys an apiValues for a given apiValues struct.
// this function will create a configmap, a service, and a deployment and then health check the application.
func DeployApi(apiValues utils.ApiValues, cid string, ctx context.Context) error {
	utils.Info(fmt.Sprintf("Creating deploy in %s environment...", apiValues.Namespace), cid)

	apiValues, err := getApiValues(apiValues, cid, ctx)
	if err != nil {
		utils.Error(fmt.Sprintf(err.Error()), cid)
		return err
	}

	if err := createService(apiValues, cid, ctx); err != nil {
		return err
	}

	if err := createConfig(apiValues, cid, ctx); err != nil {
		return err
	}

	fmt.Println("test")
	if err := createDeployment(apiValues, cid, ctx); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)
	if err := K8sHealthCheck(cid, 180, apiValues.ApiFullname, apiValues.Namespace, ctx); err != nil {
		return err
	}

	return nil
}

// createConfig deploys an configmap with the given apiValues struct.
// this function will create if it doesn't exists or patch the existing resource if exists.
func createConfig(apiValues utils.ApiValues, cid string, ctx context.Context) error {

	configmapClient := kubernetesClient.CoreV1().ConfigMaps(apiValues.Namespace)
	utils.Info(fmt.Sprintf("Creating configmap in %s environment...", apiValues.Namespace), cid)

	cm, err := getConfigmapValues(apiValues, cid, ctx)
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

// createService deploys an service with the given apiValues struct.
// this function will create if it doesn't exists or patch the existing resource if exists.
func createService(apiValues utils.ApiValues, cid string, ctx context.Context) error {

	k8sClientService := kubernetesClient.CoreV1().Services(apiValues.Namespace)
	utils.Info(fmt.Sprintf("Creating services in %s environment...", apiValues.Namespace), cid)

	service, err := k8sClientService.Get(apiValues.Name, metav1.GetOptions{})
	if err != nil {
		customErr, ok := err.(*metaErrors.StatusError)
		if !ok {
			return err
		}

		if customErr.Status().Code != 404 {
			return err
		}

		service.Name = apiValues.Name
		service.Namespace = apiValues.Namespace
		service.Spec.Type = v1core.ServiceTypeClusterIP
		service.Spec.Selector = map[string]string{"app": apiValues.Name}
		service.Spec.Ports = make([]v1core.ServicePort, 0)

		port := v1core.ServicePort{}
		port.Name = "http-" + apiValues.Name
		port.Port = int32(apiValues.HttpPort)
		port.Protocol = "TCP"
		port.TargetPort = intstr.IntOrString{Type: intstr.String, IntVal: int32(apiValues.HttpPort)}

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

// createDeployment deploys an deployment with the given apiValues struct.
// this function will create the deployment if it doesn't exists. Otherwise it will throw an error.
func createDeployment(apiValues utils.ApiValues, cid string, ctx context.Context) error {

	if err := validateCreateDeploymentArgs(apiValues); err != nil {
		utils.Error(fmt.Sprintf("Could not complete deployment due to '%s'...", err), cid)
		return err
	}

	cm, err := getConfigmapValues(apiValues, cid, ctx)
	if err != nil {
		return err
	}

	deploymentsClient := kubernetesClient.AppsV1().Deployments(apiValues.Namespace)
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
			Name: apiValues.ApiFullname,
		},
		Spec: v1apps.DeploymentSpec{
			Replicas: int32Ptr(int32(apiValues.Deployment.Replicas[apiValues.Namespace])),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":     apiValues.Name,
					"build":   apiValues.Build,
					"release": apiValues.ApiFullname,
					"version": apiValues.Version,
				},
			},
			Template: v1core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"iam.amazonaws.com/role": apiValues.Deployment.Role,
					},
					Labels: map[string]string{
						"app":     apiValues.Name,
						"build":   apiValues.Build,
						"release": apiValues.ApiFullname,
						"version": apiValues.Version,
					},
				},
				Spec: v1core.PodSpec{
					RestartPolicy: "Always",
					Containers: []v1core.Container{
						{
							Name:            apiValues.Name,
							Image:           apiValues.Deployment.Image.DockerRegistry + apiValues.Name + ":" + apiValues.Version,
							Ports:           containerPorts,
							ImagePullPolicy: "Always",
							EnvFrom: []v1core.EnvFromSource{
								v1core.EnvFromSource{
									ConfigMapRef: &v1core.ConfigMapEnvSource{
										LocalObjectReference: v1core.LocalObjectReference{
											Name: cm.Name,
										},
									},
								},
							},
							Resources: v1core.ResourceRequirements{
								Limits: v1core.ResourceList{
									"cpu":    resource.MustParse(apiValues.Resources.Limits.Cpu),
									"memory": resource.MustParse(apiValues.Resources.Limits.Memory),
								},
								Requests: v1core.ResourceList{
									"cpu":    resource.MustParse(apiValues.Resources.Requests.Cpu),
									"memory": resource.MustParse(apiValues.Resources.Requests.Memory),
								},
							},
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
	utils.Info(fmt.Sprintf("Creating new deployment %s...", apiValues.ApiFullname), cid)
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		return err
	}
	utils.Info(fmt.Sprintf("Created deployment %q!", result.GetObjectMeta().GetName()), cid)

	return nil
}

// getConfigmapValues retrieves an configmap value inside the project folder /kubernetes/{namespace}/{api-name}-config.yaml
func getConfigmapValues(apiValues utils.ApiValues, cid string, ctx context.Context) (v1core.ConfigMap, error) {
	cm := v1core.ConfigMap{}
	pwd, err := os.Getwd()
	if err != nil {
		utils.Error(fmt.Sprintf(err.Error()), cid)
		return cm, err
	}

	cmBytes, err := ioutil.ReadFile(pwd + "/" + apiValues.Name + "/kubernetes/" + apiValues.Namespace + "/" + apiValues.Name + "-config.yaml")
	if err != nil {
		utils.Fatal(fmt.Sprintf(err.Error()), cid)
		return cm, err
	}
	if err := yaml.Unmarshal(cmBytes, &cm); err != nil {
		utils.Fatal(fmt.Sprintf(err.Error()), cid)
		return cm, err
	}

	cm.Name = apiValues.ApiFullname + "-config"
	cm.Namespace = apiValues.Namespace
	utils.Info(fmt.Sprintf("Configmap extracted: %s", cm.Name), cid)
	return cm, nil
}

// getApiValues retrieves the values.yaml inside the project folder /rootProjectdir/values.yaml
func getApiValues(apiValues utils.ApiValues, cid string, ctx context.Context) (utils.ApiValues, error) {

	pwd, err := os.Getwd()
	if err != nil {
		utils.Error(fmt.Sprintf(err.Error()), cid)
		return apiValues, err
	}

	valuesBytes, err := ioutil.ReadFile(pwd + "/" + apiValues.Name + "/values.yaml")
	if err != nil {
		utils.Error(fmt.Sprintf(err.Error()), cid)
		return apiValues, err
	}

	if err := yaml.Unmarshal(valuesBytes, &apiValues); err != nil {
		utils.Error(fmt.Sprintf(err.Error()), cid)
		return apiValues, err
	}

	if apiValues.Deployment.Image.HealthCheck.Enabled && (apiValues.Deployment.Image.HealthCheck.LivenessProbeEndpoint == "" || apiValues.Deployment.Image.HealthCheck.ReadinessProbeEndpoint == "") {
		apiValues.Deployment.Image.HealthCheck.ReadinessProbeEndpoint = "/health"
		apiValues.Deployment.Image.HealthCheck.LivenessProbeEndpoint = "/health"
		apiValues.Deployment.Image.HealthCheck.HealthPort = apiValues.Deployment.Image.Ports["http"]
	}

	// set default Limits & Requests values in case of an empty ones in `values.yaml`
	if apiValues.Resources.Limits.Cpu == "" {
		apiValues.Resources.Limits.Cpu = "1"
	}

	if apiValues.Resources.Limits.Memory == "" {
		apiValues.Resources.Limits.Memory = "1Gi"
	}

	if apiValues.Resources.Requests.Cpu == "" {
		apiValues.Resources.Requests.Cpu = "0.1"
	}

	if apiValues.Resources.Requests.Memory == "" {
		apiValues.Resources.Requests.Memory = "256Mi"
	}

	utils.Info(fmt.Sprintf("Values extracted: %v", apiValues), cid)
	return apiValues, nil
}

// validateCreateDeploymentArgs this function validates the given apiValues fields necessary to build a deployment.
func validateCreateDeploymentArgs(apiValues utils.ApiValues) error {
	// if apiValues.HttpPort <= 0 && apiValues.GrpcPort <= 0 {
	// 	return errors.New(WARN_NO_PORT_SPECIFIED)
	// }

	if apiValues.Name == "" || apiValues.Namespace == "" || apiValues.Version == "" || apiValues.Build == "" {
		return errors.New(WARN_NO_NECESSARY_NAMES_SPECIFIED)
	}

	if apiValues.Deployment.Image.DockerRegistry == "" {
		return errors.New(WARN_NO_REGISTRY_FOUND)
	}

	deployment := apiValues.Deployment
	if deployment.Image.HealthCheck.Enabled && (deployment.Image.HealthCheck.HealthPort <= 0 || deployment.Image.HealthCheck.LivenessProbeEndpoint == "" || deployment.Image.HealthCheck.ReadinessProbeEndpoint == "") {
		return errors.New(WARN_NO_HEALTHCHECK_OR_READINESS_ENDPOINT_CONFIGURED)
	}

	return nil
}

func int32Ptr(i int32) *int32 { return &i }
