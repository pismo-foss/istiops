package pipeline

import (
	"context"
	"github.com/pismo/istiops/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestK8sHealthCheck(t *testing.T) {
	apiStruct := pkg.ApiStruct{
		Name:      "api-pipelinetest",
		Namespace: "default",
		Version:   "bluegreeneb",
		Build:     "2210"}

	err := K8sHealthCheck("cid-happy-test", 5, apiStruct, context.Background())
	assert.Nil(t, err)
}
