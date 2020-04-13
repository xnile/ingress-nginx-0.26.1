/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package process

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestIsRespawnIfRequired(t *testing.T) {
	cases := []struct {
		err   error
		isReq bool
	}{
		{nil, false},
		{fmt.Errorf("dummy"), false},
		{&exec.ExitError{
			ProcessState: &os.ProcessState{},
		}, true},
	}

	for _, tc := range cases {
		isReq := IsRespawnIfRequired(tc.err)
		if tc.isReq != isReq {
			t.Error("expected an error but none returned")
		}
	}
}
