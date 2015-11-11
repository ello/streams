package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateInt(t *testing.T) {
	testVal, err := ValidateInt("15", 10)

	assert.Nil(t, err, "Error should be nil on correct value")
	assert.Equal(t, 15, testVal, "validation properly parses correct value")

	testVal, err = ValidateInt("1a5", 10)

	assert.NotNil(t, err, "Error should be non nil on incorrect value")
	assert.Equal(t, 10, testVal, "validation properly returns default value")

	testVal, err = ValidateInt("", 10)

	assert.Nil(t, err, "Error should be nil on no value")
	assert.Equal(t, 10, testVal, "validation properly returns default value")
}
