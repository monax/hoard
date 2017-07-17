package logging

import (
	"testing"

	"bytes"

	"github.com/monax/hoard/core/logging/structure"
	"github.com/stretchr/testify/assert"
)

func TestTerminalLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	logger, err := NewTerminalLogger(Json, buf)
	assert.NoError(t, err)
	logger.Log(structure.ChannelKey, "foo")
	assert.Equal(t, "{\"channel\":\"foo\"}\n", buf.String())
}

func TestLoggerFromLoggingConfig(t *testing.T) {
	buf := new(bytes.Buffer)
	logger, err := LoggerFromLoggingConfig(DefaultConfig, buf)
	assert.NoError(t, err)
	logger.Log(structure.ChannelKey, structure.TraceChannel, "foo", "bar")
	assert.Equal(t, "channel=trace foo=bar\n", buf.String())
}
