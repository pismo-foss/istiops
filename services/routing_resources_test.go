package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateHttpRouteResource(t *testing.T) {
	apiStruct := ApiStruct{
		Name:      "api-pipetest",
		Namespace: "default",
		Version:   "bluegreeneb",
		Build:     "2210",
		HttpPort:  8080,
	}

	err := CreateRouteResource(apiStruct, "test-http-happy", context.Background())
	assert.Nil(t, err)
}

func TestCreateGrpcRouteResource(t *testing.T) {
	apiStruct := ApiStruct{
		Name:      "api-pipetest",
		Namespace: "default",
		Version:   "bluegreeneb",
		Build:     "2210",
		GrpcPort:  8777,
		HttpPort:  8005,
	}

	err := CreateRouteResource(apiStruct, "test-grpc-happy", context.Background())
	assert.Nil(t, err)
}
