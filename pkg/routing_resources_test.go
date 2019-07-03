package pkg

import (
	"context"
	"github.com/pismo/istiops/pipeline"
	"github.com/pismo/istiops/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateHttpRouteResource(t *testing.T) {
	apiStruct := utils.BuildApiStruct("api-pipelinetest", "default", "1.0.0", "2210")
	apiStruct.HttpPort = 8080
	apiStruct.ApiHostName = apiStruct.Name + "-" + apiStruct.Namespace + pipeline.PismoDomains[apiStruct.Name]

	err := CreateRouteResource(apiStruct, "test-http-happy", context.Background())
	assert.Nil(t, err)
}

func TestCreateGrpcRouteResource(t *testing.T) {
	apiStruct := utils.BuildApiStruct("api-pipelinetest", "default", "1.0.0", "2210")
	apiStruct.HttpPort = 8080
	apiStruct.GrpcPort = 8777
	apiStruct.ApiHostName = apiStruct.Name + "-" + apiStruct.Namespace + pipeline.PismoDomains[apiStruct.Name]

	err := CreateRouteResource(apiStruct, "test-grpc-happy", context.Background())
	assert.Nil(t, err)
}
