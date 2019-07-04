package utils

import (
	"fmt"
	"strings"
)

func BuildApiValues(name string, namespace string, version string, build string) ApiValues {
	apiValues := ApiValues{
		Name:      name,
		Namespace: namespace,
		Version:   version,
		Build:     build,
	}
	replacer := strings.NewReplacer(".", "", "-", "", "/", "")
	simplifiedVersion := replacer.Replace(apiValues.Version)
	simplifiedVersion = strings.ToLower(simplifiedVersion)

	apiValues.ApiFullname = fmt.Sprintf("%s-%s-%s-%s",
		apiValues.Name,
		apiValues.Namespace,
		simplifiedVersion,
		apiValues.Build)

	return apiValues
}

type ApiValues struct {
	Name         string                   `json:"name"`
	ApiHostName  string                   `json:"api_host_name"`
	ApiFullname  string                   `json:"api_fullname"`
	Namespace    string                   `json:"namespace"`
	Version      string                   `json:"version"`
	Build        string                   `json:"build"`
	HttpPort     uint32                   `json:"http_port"`
	GrpcPort     uint32                   `json:"grpc_port"`
	Deployment   Deployment               `yaml:"deployment"`
	Resources    Resources                `yaml:"resources"`
	NodeSelector map[string]interface{}   `yaml:"nodeSelector"`
	Tolerations  []map[string]interface{} `yaml:"tolerations"`
	Affinity     map[string]interface{}   `yaml:"affinity"`
}

type Resources struct {
	Limits   Limits   `yaml:"limits"`
	Requests Requests `yaml:"requests"`
}

type Limits struct {
	Cpu    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

type Requests struct {
	Cpu    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

type Deployment struct {
	Role     string           `yaml:"role"`
	Replicas map[string]int64 `yaml:"replicas"`
	Image    Image            `yaml:"image"`
}

type Image struct {
	HealthCheck    Probes           `yaml:"healthCheck"`
	Ports          map[string]int64 `yaml:"ports"`
	DockerRegistry string           `yaml:"dockerRegistry"`
	PullPolicy     string           `yaml:"pullPolicy"`
}

type Probes struct {
	HealthPort             int64  `yaml:"healthPort"`
	LivenessProbeEndpoint  string `yaml:"livenessProbeEndpoint"`
	ReadinessProbeEndpoint string `yaml:"readinessProbeEndpoint"`
	Enabled                bool   `yaml:"enabled"`
}
