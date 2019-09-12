package router

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var fakeIstioClient IstioClientInterface

func TestStringify_Unit(t *testing.T) {
	mapSelector := map[string]string{}
	mapSelector["app"] = "api-domain"

	stringified, err := Stringify("integration-tests-uuid", mapSelector)
	assert.NoError(t, err)
	assert.Equal(t, "app=api-domain", stringified)
}

func TestStringify_Unit_EmptyLabelSelector(t *testing.T) {
	labelSelector := map[string]string{}

	_, err := Stringify("", labelSelector)
	assert.EqualError(t, err, "got an empty labelSelector")
}

func TestMapify_Unit(t *testing.T) {
	cases := []struct {
		selector string
		want     map[string]string
	}{
		{
			"app=api-domain,role=aws/my-role",
			map[string]string{
				"app":  "api-domain",
				"role": "aws/my-role",
			},
		},
		{
			"role=aws/my-role,app=api-domain",
			map[string]string{
				"app":  "api-domain",
				"role": "aws/my-role",
			},
		},
		{
			"app=api-domain,version=2.1.3',role=aws/my-role",
			map[string]string{
				"app":     "api-domain",
				"version": "2.1.3'",
				"role":    "aws/my-role",
			},
		},
	}

	for _, tt := range cases {
		mapified, err := Mapify("", tt.selector)
		assert.Equal(t, tt.want, mapified)
		assert.NoError(t, err)
	}
}

func TestMapify_Unit_EmptyLabelSelector(t *testing.T) {
	_, err := Mapify("", "")
	assert.EqualError(t, err, "got an empty labelSelector string")
}

func TestMapify_Unit_MalformedLabelSelector(t *testing.T) {
	_, err := Mapify("", "app:domain")
	assert.EqualError(t, err, "missing '=' operator for labelSelector")
}
