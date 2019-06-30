package pipeline

import (
	"context"
	"github.com/pismo/istiops/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeployHelm(t *testing.T) {
	apiStruct := utils.BuildApiStruct("api-pipelinetest", "default", "1.0.0", "2210")
	apiStruct.HttpPort = 8080
	err := DeployApi(apiStruct, "cid-happy-yest", context.Background())
	assert.Nil(t, err)
}
