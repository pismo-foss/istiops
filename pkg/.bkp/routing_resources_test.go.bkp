package pkg

import (
	"context"
	"github.com/pismo/istiops/pipeline"
	"github.com/pismo/istiops/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateHttpRouteResource(t *testing.T) {
	apiValues := utils.BuildApiValues("api-pipelinetest", "default", "1.0.0", "2210")
	apiValues.HttpPort = 8080
	apiValues.ApiHostName = apiValues.Name + "-" + apiValues.Namespace + pipeline.PismoDomains[apiValues.Name]

	err := CreateRouteResource(apiValues, "test-http-happy", context.Background())
	assert.Nil(t, err)
}

func TestCreateGrpcRouteResource(t *testing.T) {
	apiValues := utils.BuildApiValues("api-pipelinetest", "default", "1.0.0", "2210")
	apiValues.HttpPort = 8080
	apiValues.GrpcPort = 8777
	apiValues.ApiHostName = apiValues.Name + "-" + apiValues.Namespace + pipeline.PismoDomains[apiValues.Name]

	err := CreateRouteResource(apiValues, "test-grpc-happy", context.Background())
	assert.Nil(t, err)
}
