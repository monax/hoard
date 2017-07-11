package logging

import (
	"io"

	"fmt"

	"os"

	"code.monax.io/platform/hoard/core/logging"
	"code.monax.io/platform/hoard/core/logging/loggers"
	"code.monax.io/platform/hoard/core/logging/structure"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/term"
)

type LoggingType string

const (
	Logfmt LoggingType = "logfmt"
	Json   LoggingType = "json"
)

type LoggingConfig struct {
	// Channels to listen on, messages not on these channels will be filtered
	// leaving empty disables logging
	LoggingType LoggingType
	Channels    []structure.Channel
}

func NewLoggingConfig(loggingType LoggingType,
	channels ...structure.Channel) *LoggingConfig {

	return &LoggingConfig{
		LoggingType: loggingType,
		Channels:    channels,
	}
}

func LoggerFromLoggingConfig(loggingConfig *LoggingConfig) (log.Logger, error) {
	terminalLogger, err := TerminalLogger(loggingConfig.LoggingType, os.Stderr)
	if err != nil {
		return nil, err
	}
	return loggers.NewFilterLogger(terminalLogger,
		excludeChannelsNotIn(loggingConfig.Channels)), nil
}

func TerminalLogger(loggingType LoggingType, writer io.Writer) (log.Logger, error) {
	loggerMaker, err := OutputLoggerMaker(loggingType)
	if err != nil {
		return nil, err
	}
	return term.NewLogger(writer, loggerMaker, logging.Colors), nil
}

func OutputLoggerMaker(loggingType LoggingType) (func(writer io.Writer) log.Logger, error) {
	switch loggingType {
	case Logfmt:
		return func(writer io.Writer) log.Logger {
			return log.NewLogfmtLogger(writer)
		}, nil
	case Json:
		return func(writer io.Writer) log.Logger {
			return log.NewJSONLogger(writer)
		}, nil
	default:
		return nil, fmt.Errorf("Could not create logger with logging "+
			"type '%s'.", loggingType)
	}
}

func excludeChannelsNotIn(includeChannels []structure.Channel) func(keyvals []interface{}) bool {
	return func(keyvals []interface{}) bool {
		channel := structure.Value(keyvals, structure.ChannelKey)

		if channel != nil {
			for includeChannel := range includeChannels {
				if channel == includeChannel {
					// Do NOT filter
					return false
				}
			}
		}
		// Filter
		return true
	}
}
