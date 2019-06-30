package pipeline

import (
	"context"
	"github.com/pismo/istiops/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestK8sHealthCheck(t *testing.T) {
	apiStruct := utils.ApiStruct{
		Name:      "api-pipelinetest",
		Namespace: "default",
		Version:   "bluegreeneb",
		Build:     "2210"}

	err := K8sHealthCheck("cid-happy-test", 5, apiStruct, context.Background())
	assert.Nil(t, err)
}
