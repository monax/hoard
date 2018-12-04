package logging

import (
	"io"

	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/term"
	"github.com/monax/hoard/logging/loggers"
	"github.com/monax/hoard/logging/structure"
)

type LoggingType string

var DefaultConfig = NewLoggingConfig(Json, structure.InfoChannel,
	structure.TraceChannel)

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

func LoggerFromLoggingConfig(loggingConfig *LoggingConfig,
	writer io.Writer) (log.Logger, error) {
	terminalLogger, err := NewTerminalLogger(loggingConfig.LoggingType, writer)
	if err != nil {
		return nil, err
	}
	return loggers.NewFilterLogger(terminalLogger,
		excludeChannelsNotIn(loggingConfig.Channels)), nil
}

func NewTerminalLogger(loggingType LoggingType, writer io.Writer) (log.Logger, error) {
	loggerMaker, err := OutputLoggerMaker(loggingType)
	if err != nil {
		return nil, err
	}
	return term.NewLogger(writer, loggerMaker, channelColours), nil
}

func OutputLoggerMaker(loggingType LoggingType) (func(writer io.Writer) log.Logger, error) {
	var logger func(io.Writer) log.Logger
	switch loggingType {
	case Logfmt:
		logger = log.NewLogfmtLogger
	case Json:
		logger = log.NewJSONLogger
	default:
		return nil, fmt.Errorf("could not create logger with logging type '%s'", loggingType)
	}

	return func(writer io.Writer) log.Logger {
		return logger(writer)
	}, nil
}

func excludeChannelsNotIn(includeChannels []structure.Channel) func(keyvals []interface{}) bool {
	return func(keyvals []interface{}) bool {
		channel := structure.Value(keyvals, structure.ChannelKey)

		if channel != nil {
			for _, includeChannel := range includeChannels {
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

func channelColours(keyvals ...interface{}) term.FgBgColor {
	channel := structure.Value(keyvals, structure.ChannelKey)
	switch channel {
	case structure.TraceChannel:
		return term.FgBgColor{Fg: term.DarkBlue}
	case structure.InfoChannel:
		return term.FgBgColor{Fg: term.DarkGreen}
	default:
		return term.FgBgColor{}
	}
}
