package pipeline

import (
	"context"
	"github.com/pismo/istiops/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestK8sHealthCheck(t *testing.T) {
	apiValues := utils.BuildApiValues("api-pipelinetest", "default", "1.0.0", "2210")
	err := K8sHealthCheck("cid-happy-test", 5, apiValues.ApiFullname, apiValues.Namespace, context.Background())
	assert.Nil(t, err)
}
