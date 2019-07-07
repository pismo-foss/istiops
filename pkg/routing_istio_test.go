package pkg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShaFromMap(t *testing.T) {
	mapToHash := map[string]string{
		"key": "value",
		"app": "api-xpto",
	}

	values, err := GenerateShaFromMap(mapToHash)

	assert.NoError(t, err)
	assert.Equal(
		t,
		fmt.Sprintf("%s", values[0]),
		"563f0357118d05ef145d6bddf2966cc23e86ca8f2f013f915e565afdf09f7a23",
	)
	assert.Equal(
		t,
		fmt.Sprintf("%s", values[1]),
		"dc705ce67cc0d8a71c4449a2933fb0ac0404e111325ee7a08c27f2a17fe4a9e2",
	)
}

func TestCompareMapsKeyPairsHash(t *testing.T) {
	mapMocked := map[string]string{
		"key": "value",
		"app": "api-xpto",
	}

	mapEmpty := map[string]string{}

	expectedTrue := CompareMapsKeyPairsHash(mapMocked, mapMocked)
	expectedFalse := CompareMapsKeyPairsHash(mapMocked, mapEmpty)

	assert.True(t, expectedTrue)
	assert.True(t, expectedFalse)
}
