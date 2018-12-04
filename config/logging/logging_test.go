package logging

import (
	"testing"

	"bytes"

	"github.com/monax/hoard/logging/structure"
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
	logger, err := LoggerFromLoggingConfig(NewLoggingConfig(Logfmt, structure.TraceChannel), buf)
	assert.NoError(t, err)
	logger.Log(structure.ChannelKey, structure.TraceChannel, "foo", "bar")
	assert.Equal(t, "channel=trace foo=bar\n", buf.String())
}
