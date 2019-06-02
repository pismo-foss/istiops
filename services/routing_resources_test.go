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
		ApiValues: &ApiValues{Deployment: Deployment{Image: Image{Ports: map[string]uint32{"http": 8080}}}},
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
		ApiValues: &ApiValues{Deployment: Deployment{Image: Image{Ports: map[string]uint32{"http": 8005, "grpc": 8777}}}},
	}

	err := CreateRouteResource(apiStruct, "test-grpc-happy", context.Background())
	assert.Nil(t, err)
}
