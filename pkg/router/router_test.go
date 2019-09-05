package router

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringify_Unit(t *testing.T) {
	cases := []struct {
		mapSelector    map[string]string
		want  string
	}{
		{
			map[string]string{
				"app": "api-domain",
				"role": "aws/my-role",
			},
			"app=api-domain,role=aws/my-role",
		},
		{
			map[string]string{
				"app": "api-domain",
				"version": "2.1.3'",
				"role": "aws/my-role",
			},
			"app=api-domain,version=2.1.3',role=aws/my-role",
		},
	}

	for _, tt := range cases {
		stringified, err := Stringify("", tt.mapSelector)
		assert.Equal(t, tt.want, stringified)
		assert.NoError(t, err)
	}
}

func TestStringify_Unit_EmptyLabelSelector(t *testing.T) {
	labelSelector := map[string]string{}

	_, err := Stringify("", labelSelector)
	assert.EqualError(t, err, "got an empty labelSelector")
}

func TestMapify_Unit(t *testing.T) {
	cases := []struct {
		selector    string
		want  map[string]string
	}{
		{
			"app=api-domain,role=aws/my-role",
			map[string]string{
				"app": "api-domain",
				"role": "aws/my-role",
			},
		},
		{
			"role=aws/my-role,app=api-domain",
			map[string]string{
				"app": "api-domain",
				"role": "aws/my-role",
			},
		},
		{
			"app=api-domain,version=2.1.3',role=aws/my-role",
			map[string]string{
				"app": "api-domain",
				"version": "2.1.3'",
				"role": "aws/my-role",
			},
		},
	}

	for _, tt := range cases {
		mapified, err := Mapify("", tt.selector)
		assert.Equal(t, tt.want, mapified)
		t.Log(tt.want)
		assert.NoError(t, err)
	}
}

func TestMapify_Unit_EmptyLabelSelector(t *testing.T) {
	_, err := Mapify("","")
	assert.EqualError(t, err, "got an empty labelSelector string")
}

func TestMapify_Unit_MalformedLabelSelector(t *testing.T) {
	_, err := Mapify("","app:domain")
	assert.EqualError(t, err, "missing '=' operator for labelSelector")
}