package interfaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrUnknownServiceExists(t *testing.T) {
	assert.NotNil(t, ErrUnknownService)
}
