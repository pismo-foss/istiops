package pipeline

import (
	"context"
	"github.com/pismo/istiops/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeployHelm(t *testing.T) {
	apiValues := utils.BuildApiValues("api-pipelinetest", "default", "1.0.0", "2210")
	apiValues.HttpPort = 8080
	err := DeployApi(apiValues, "cid-happy-yest", context.Background())
	assert.Nil(t, err)
}
