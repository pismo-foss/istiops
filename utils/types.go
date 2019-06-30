package utils

type ApiStruct struct {
	Name      string `json:"name"`
	ApiHostName string `json:"api_host_name"`
	Namespace string `json:"namespace"`
	Version   string `json:"version"`
	Build     string `json:"build"`
	HttpPort  uint32 `json:"http_port"`
	GrpcPort  uint32 `json:"grpc_port"`
}

type ApiValues struct {
	Deployment   Deployment               `yaml:"deployment"`
	Resources    map[string]interface{}   `yaml:"resources"`
	NodeSelector map[string]interface{}   `yaml:"nodeSelector"`
	Tolerations  []map[string]interface{} `yaml:"tolerations"`
	Affinity     map[string]interface{}   `yaml:"affinity"`
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

