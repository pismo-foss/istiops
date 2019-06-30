package pipeline

import (
	"context"
	"github.com/pismo/istiops/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeployHelm(t *testing.T) {
	apiStruct := pkg.ApiStruct{
		Name:      "api-pipelinetest",
		Namespace: "default",
		Version:   "bluegreeneb",
		Build:     "2210"}

	err := DeployHelm(apiStruct, "cid-happy-yest", context.Background())
	assert.Nil(t, err)
}
