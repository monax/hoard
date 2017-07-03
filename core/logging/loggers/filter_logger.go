package loggers

import "github.com/go-kit/kit/log"

// Filter logger allows us to filter lines logged to it before passing on to underlying
// output logger
type filterLogger struct {
	logger    log.Logger
	predicate func(keyvals []interface{}) bool
}

var _ log.Logger = (*filterLogger)(nil)

func (fl filterLogger) Log(keyvals ...interface{}) error {
	if !fl.predicate(keyvals) {
		return fl.logger.Log(keyvals...)
	}
	return nil
}

// Creates a logger that removes lines from output when the predicate evaluates true
func NewFilterLogger(outputLogger log.Logger,
	predicate func(keyvals []interface{}) bool) log.Logger {
	return &filterLogger{
		logger:    outputLogger,
		predicate: predicate,
	}
}
