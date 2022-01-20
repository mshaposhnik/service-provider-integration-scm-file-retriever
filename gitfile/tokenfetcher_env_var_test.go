// Copyright (c) 2021 - 2022 Red Hat, Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gitfile

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateHeaderStructFromEnv(t *testing.T) {
	//t.Setenv("TOKEN", "abcd_foo") // only since 1.17
	os.Setenv("TOKEN", "abcd_foo")
	defer os.Unsetenv("TOKEN")
	headerStruct := new(EnvVarTokenFetcher).BuildHeader("https://github.com/any/test.git")
	assert.Equal(t, "Bearer abcd_foo", headerStruct.Authorization, "Authorization header value mismatch")
}
