// Copyright 2017 Monax Industries Limited
//
// Licensed under the Apache License, Type 2.0 (the "License");
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

package loggers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVectorValuedLogger(t *testing.T) {
	logger := newTestLogger()
	vvl := VectorValuedLogger(logger)
	vvl.Log("foo", "bar", "seen", 1, "seen", 3, "seen", 2)
	lls, err := logger.logLines(1)
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{"foo", "bar", "seen", []interface{}{1, 3, 2}},
		lls[0])
}
