package utils

import (
	"fmt"
	"strings"
)

func BuildApiStruct(name string, namespace string, version string, build string) ApiStruct {
	apiStruct := ApiStruct{
		Name:      name,
		Namespace: namespace,
		Version:   version,
		Build:     build,
	}
	replacer := strings.NewReplacer(".", "", "-", "", "/", "")
	simplifiedVersion := replacer.Replace(apiStruct.Version)
	simplifiedVersion = strings.ToLower(simplifiedVersion)

	apiStruct.ApiFullname = fmt.Sprintf("%s-%s-%s-%s",
		apiStruct.Name,
		apiStruct.Namespace,
		simplifiedVersion,
		apiStruct.Build)

	return apiStruct
}

type ApiStruct struct {
	Name        string    `json:"name"`
	ApiHostName string    `json:"api_host_name"`
	ApiFullname string    `json:"api_fullname"`
	Namespace   string    `json:"namespace"`
	Version     string    `json:"version"`
	Build       string    `json:"build"`
	HttpPort    uint32    `json:"http_port"`
	GrpcPort    uint32    `json:"grpc_port"`
	ApiValues   ApiValues `json:"api_values"`
}

type ApiValues struct {
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
