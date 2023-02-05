package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultField(t *testing.T) {
	logger := NewLogger(true, true)

	assert.NoError(t, logger.SetDefaultField("default_field", "test"))
	assert.Error(t, logger.SetDefaultField("message", "test"))
	assert.Error(t, logger.SetDefaultField("component", "test"))
	assert.Error(t, logger.SetDefaultField("context", "test"))
	assert.Error(t, logger.SetDefaultField("level", "test"))
	assert.Error(t, logger.SetDefaultField("timestamp", "test"))
	assert.Error(t, logger.SetDefaultField("error", "test"))
	assert.Error(t, logger.SetDefaultField("request_id", "test"))
}
