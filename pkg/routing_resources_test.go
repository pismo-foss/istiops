package pkg

import (
	"context"
	"github.com/pismo/istiops/pipeline"
	"github.com/pismo/istiops/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateHttpRouteResource(t *testing.T) {
	apiStruct := utils.ApiStruct{
		Name:      "api-pipetest",
		Namespace: "default",
		Version:   "bluegreeneb",
		Build:     "2210",
		HttpPort:  8080,
	}
	apiStruct.ApiHostName = apiStruct.Name + "-" + apiStruct.Namespace + pipeline.PismoDomains[apiStruct.Name]

	err := CreateRouteResource(apiStruct, "test-http-happy", context.Background())
	assert.Nil(t, err)
}

func TestCreateGrpcRouteResource(t *testing.T) {
	apiStruct := utils.ApiStruct{
		Name:      "api-pipetest",
		Namespace: "default",
		Version:   "bluegreeneb",
		Build:     "2210",
		GrpcPort:  8777,
		HttpPort:  8005,
	}
	apiStruct.ApiHostName = apiStruct.Name + "-" + apiStruct.Namespace + pipeline.PismoDomains[apiStruct.Name]

	err := CreateRouteResource(apiStruct, "test-grpc-happy", context.Background())
	assert.Nil(t, err)
}
