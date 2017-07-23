package interfaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrUnknownFrontendExists(t *testing.T) {
	assert.NotNil(t, ErrUnknownFrontend)
}
