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

package structure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVectorise(t *testing.T) {
	kvs := []interface{}{
		"scope", "lawnmower",
		"hub", "budub",
		"occupation", "fish brewer",
		"scope", "hose pipe",
		"flub", "dub",
		"scope", "rake",
		"flub", "brub",
	}

	kvsVector := Vectorise(kvs, "occupation", "scope")
	// Vectorise scope
	assert.Equal(t, []interface{}{
		"scope", []interface{}{"lawnmower", "hose pipe", "rake"},
		"hub", "budub",
		"occupation", "fish brewer",
		"flub", []interface{}{"dub", "brub"},
	},
		kvsVector)
}
