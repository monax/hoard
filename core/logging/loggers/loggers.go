package loggers

import "github.com/go-kit/kit/log"

// Apply a pipeline of log handlers, the handlers will are applied with the from
// right to left so the earlier handlers are the outermost ones, so for an
// incoming message 'msg' and handlers 'h1', 'h2', 'h3':
// msg -> Compose(h1,h2,h3)(baseLogger) = msg -> h1 -> h2 -> h3 -> baseLogger
func Compose(handlers ...func(log.Logger) log.Logger) func(log.Logger) log.Logger {
	return func(logger log.Logger) log.Logger {
		for i := len(handlers) - 1; i >= 0; i-- {
			logger = handlers[i](logger)
		}

		return logger
	}
}
