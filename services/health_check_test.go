package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestK8sHealthCheck(t *testing.T) {
	apiStruct := ApiStruct{
		Name:      "api-statements",
		Namespace: "ext",
		Version:   "bluegreeneb",
		Build:     "2210"}

	err := K8sHealthCheck("cid-happy-test", 5, apiStruct, context.Background())
	assert.Nil(t, err)
}
