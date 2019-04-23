/*
Copyright The Helm Authors.

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

package action

import (
	"testing"
)

func releaseTestRunAction(t *testing.T, testEnv *reltesting.Environment) *ReleaseTest {
	config := actionConfigFixture(t)
	return NewReleaseTesting(config)
}

func TestReleaseTesting(t *testing.T) {
	name := "mockRelease"
	rt := NewReleaseTesting(cfg)
	ch, errc := rt.Run(name)
}

type successRelTestingEnv struct {
}

func (s *succesfulRelTestingEnv) CreateTestPod(test *reltesting.Test) error { return nil }
func (s *succesfulRelTestingEnv) GetTestPodStatus(test *reltesting.Test) (v1.PodPhase, error) {
	return v1.PodSuccess, nil
}
func (s *succesfulRelTestingEnv) DeleteTestPods(testManifests []string) {}
func (s *succesfulRelTestingEnv) StreamError(info string) error         { return nil }
func (s *succesfulRelTestingEnv) StreamFailed(info string) error        { return nil }
func (s *succesfulRelTestingEnv) StreamMessage(msg string, status release.TestStatus) error {
	return nil
}
