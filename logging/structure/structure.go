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

package structure

type Channel string

const (
	// Log time (time.Time)
	TimeKey = "time"
	// Call site for log invocation (go-stack.Call)
	CallerKey = "caller"
	// Log message (string)
	MessageKey = "message"
	// Top-level component (choose one) name
	ComponentKey = "component"
	// Vector-valued scope
	ScopeKey = "scope"
	// Globally unique identifier persisting while a single instance (root process)
	// of this program/service is running
	RunId = "run_id"
	// Stack trace for log call
	StackTraceKey = "trace"
	// Logging channel, Info or Trace
	ChannelKey = "channel"
	// Go error value
	ErrorKey = "error"

	// Channels (can be used for semantic filtering through LoggingConfig)
	InfoChannel  Channel = "info"
	TraceChannel Channel = "trace"
)

// Stateful index that tracks the location of a possible vector value
type vectorValueindex struct {
	// Location of the value belonging to a key in output slice
	valueIndex int
	// Whether or not the value is currently a vector
	vector bool
}

// 'Vectorises' values associated with repeated string keys member by collapsing many values into a single vector value.
// The result is a copy of keyvals where the first occurrence of each matching key and its first value are replaced by
// that key and all of its values in a single slice.
func Vectorise(keyvals []interface{}, vectorKeys ...string) []interface{} {
	// We rely on working against a single backing array, so we use a capacity that is the maximum possible size of the
	// slice after vectorising (in the case there are no duplicate keys and this is a no-op)
	outputKeyvals := make([]interface{}, 0, len(keyvals))
	// Track the location and vector status of the values in the output
	valueIndices := make(map[string]*vectorValueindex, len(vectorKeys))
	elided := 0
	for i := 0; i < 2*(len(keyvals)/2); i += 2 {
		key := keyvals[i]
		val := keyvals[i+1]

		// Only attempt to vectorise string keys
		if k, ok := key.(string); ok {
			if valueIndices[k] == nil {
				// Record that this key has been seen once
				valueIndices[k] = &vectorValueindex{
					valueIndex: i + 1 - elided,
				}
				// Copy the key-value to output with the single value
				outputKeyvals = append(outputKeyvals, key, val)
			} else {
				// We have seen this key before
				vi := valueIndices[k]
				if !vi.vector {
					// This must be the only second occurrence of the key so now vectorise the value
					outputKeyvals[vi.valueIndex] = []interface{}{outputKeyvals[vi.valueIndex]}
					vi.vector = true
				}
				// Grow the vector value
				outputKeyvals[vi.valueIndex] = append(outputKeyvals[vi.valueIndex].([]interface{}), val)
				// We are now running two more elements behind the input keyvals because we have absorbed this key-value pair
				elided += 2
			}
		} else {
			// Just copy the key-value to the output for non-string keys
			outputKeyvals = append(outputKeyvals, key, val)
		}
	}
	return outputKeyvals
}

// Return a single value corresponding to key in keyvals
func Value(keyvals []interface{}, key interface{}) interface{} {
	for i := 0; i < 2*(len(keyvals)/2); i += 2 {
		if keyvals[i] == key {
			return keyvals[i+1]
		}
	}
	return nil
}

// Maps key values pairs with a function (key, value) -> (new key, new value)
func MapKeyValues(keyvals []interface{}, fn func(interface{}, interface{}) (interface{}, interface{})) []interface{} {
	mappedKeyvals := make([]interface{}, len(keyvals))
	for i := 0; i < 2*(len(keyvals)/2); i += 2 {
		key := keyvals[i]
		val := keyvals[i+1]
		mappedKeyvals[i], mappedKeyvals[i+1] = fn(key, val)
	}
	return mappedKeyvals
}
