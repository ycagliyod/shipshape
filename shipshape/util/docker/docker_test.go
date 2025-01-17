/*
 * Copyright 2015 Google Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package docker

import (
	"os/exec"
	"testing"
	"time"
)

func TestHasDocker(t *testing.T) {
	if got, want := HasDocker(), true; got != want {
		t.Errorf("Unexpected error for HasDocker test: got %v, expected %v",
			got, want)
	}
}

// NB! This test is using a set of predefined container names that might clash with already
// existing container names (docker_test_container*). The tests clean up after themselves so
// matching containers will be cleaned up in the test tear down. Rerunning a test should
// take care of the clash from the test perspective. However, a non-test related container with
// a colliding name will not be restored.
func TestContainerExists(t *testing.T) {
	tests := []struct {
		desc      string
		container string
		setup     *exec.Cmd
		teardown  *exec.Cmd
		exists    bool
	}{
		{
			desc:      "Detect matching running container",
			container: "docker_test_container1",
			setup:     exec.Command("docker", "run", "--name=docker_test_container1", "ubuntu:14.04", "sleep 30"),
			teardown:  exec.Command("docker", "rm", "docker_test_container1"),
			exists:    true,
		},
		{
			desc:      "Don't detect non-matching running container",
			container: "someother_container",
			setup:     exec.Command("docker", "run", "--name=docker_test_container2", "ubuntu:14.04", "sleep 30"),
			teardown:  exec.Command("docker", "rm", "docker_test_container2"),
			exists:    false,
		},
		{
			desc:      "Detect matching non-running container",
			container: "docker_test_container3",
			setup:     exec.Command("docker", "run", "--name=docker_test_container3", "ubuntu:14.04", "sleep 0"),
			teardown:  exec.Command("docker", "rm", "docker_test_container3"),
			exists:    true,
		},
		{
			desc:      "Don't detect non-matching non-running container",
			container: "someother_container",
			setup:     exec.Command("docker", "run", "--name=docker_test_container4", "ubuntu:14.04", "sleep 0"),
			teardown:  exec.Command("docker", "rm", "docker_test_container4"),
			exists:    false,
		},
		{
			desc:      "Don't detect container name substrings",
			container: "ocker_test_containe",
			setup:     exec.Command("docker", "run", "--name=docker_test_container4", "ubuntu:14.04", "sleep 0"),
			teardown:  exec.Command("docker", "rm", "docker_test_container4"),
			exists:    false,
		},
		{
			desc:      "Don't detect container name superstrings",
			container: "2docker_test_container2",
			setup:     exec.Command("docker", "run", "--name=docker_test_container4", "ubuntu:14.04", "sleep 0"),
			teardown:  exec.Command("docker", "rm", "docker_test_container4"),
			exists:    false,
		},
	}

	for _, test := range tests {
		test.setup.Run()
		time.Sleep(100 * time.Millisecond)
		got, err := ContainerExists(test.container)
		want := test.exists
		if err != nil {
			t.Errorf("%v: Error running test; %v", test.desc, err)
		} else if got != want {
			t.Errorf("%v: Wrong result; got %v, expected %v",
				test.desc, test.container, got, want)
		}
		test.teardown.Run()
	}
}
