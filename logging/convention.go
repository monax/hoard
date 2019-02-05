// Copyright 2017 Monax Industries Limited
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/logging/structure"
)

func InfoLogger(logger log.Logger) log.Logger {
	return log.WithPrefix(logger, structure.ChannelKey, structure.InfoChannel)
}
func TraceLogger(logger log.Logger) log.Logger {
	return log.WithPrefix(logger, structure.ChannelKey, structure.TraceChannel)
}

func InfoMsg(logger log.Logger, message string, keyvals ...interface{}) {
	Msg(InfoLogger(logger), message, keyvals...)
}

// Record structured Trace log line with a message
func TraceMsg(logger log.Logger, message string, keyvals ...interface{}) {
	Msg(TraceLogger(logger), message, keyvals...)
}

// Establish or extend the scope of this logger by appending scopeName to the Scope vector.
// Like With the logging scope is append only but can be used to provide parenthetical scopes by hanging on to the
// parent scope and using once the scope has been exited. The scope mechanism does is agnostic to the type of scope
// so can be used to identify certain segments of the call stack, a lexical scope, or any other nested scope.
func WithScope(logger log.Logger, scopeName string) log.Logger {
	// InfoTraceLogger will collapse successive (ScopeKey, scopeName) pairs into a vector in the order which they appear
	return log.With(logger, structure.ScopeKey, scopeName)
}

// Record a structured log line with a message
func Msg(logger log.Logger, message string, keyvals ...interface{}) error {
	prepended := make([]interface{}, 0, len(keyvals)+2)
	prepended = append(prepended, structure.MessageKey, message)
	prepended = append(prepended, keyvals...)
	return logger.Log(prepended...)
}

// Log an error with consistent error key and metadata and return the error passed in
// or an logging error enclosing the original error if there is a logging error
func Err(logger log.Logger, err error) error {
	if err != nil {
		errLogger := logger.Log(structure.MessageKey, "Failure", structure.ErrorKey, err)
		if errLogger != nil {
			return fmt.Errorf("failed to log error '%s': %s", errLogger, err)
		}
		return err
	}
	return nil
}
