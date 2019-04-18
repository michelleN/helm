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

package releasetesting

import (
	"bytes"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"

	"helm.sh/helm/pkg/kube"
	"helm.sh/helm/pkg/release"
)

// EnvironmentInterfacee
type TestEnvironment interface {
	CreateTestPod(test *test) error
	GetTestPodStatus(test *test) (v1.PodPhase, error)
	DeleteTestPods(testManifests []string)

	StreamError(info string) error
	StreamFailed(name string) error
	StreamMessage(msg string, status release.TestRunStatus) error
	StreamResult(r *release.TestRun) error
	StreamRunning(name string) error
	StreamSuccess(name string) error
	StreamUnknown(name string, info string) error
}

// Environment encapsulates information about where test suite executes and returns results
type Environment struct {
	Namespace  string
	KubeClient kube.KubernetesClient
	Messages   chan *release.TestReleaseResponse
	Timeout    int64
}

func (env *Environment) CreateTestPod(test *test) error {
	b := bytes.NewBufferString(test.manifest)
	if err := env.KubeClient.Create(env.Namespace, b, env.Timeout, false); err != nil {
		test.result.Info = err.Error()
		test.result.Status = release.TestRunFailure
		return err
	}

	return nil
}

func (env *Environment) GetTestPodStatus(test *test) (v1.PodPhase, error) {
	status, err := env.KubeClient.WaitAndGetCompletedPodPhase(env.Namespace, test.name, env.Timeout)
	if err != nil {
		log.Printf("Error getting status for pod %s: %s", test.result.Name, err)
		test.result.Info = err.Error()
		test.result.Status = release.TestRunUnknown
		return status, err
	}

	return status, err
}

func (env *Environment) StreamResult(r *release.TestRun) error {
	switch r.Status {
	case release.TestRunSuccess:
		if err := env.StreamSuccess(r.Name); err != nil {
			return err
		}
	case release.TestRunFailure:
		if err := env.StreamFailed(r.Name); err != nil {
			return err
		}

	default:
		if err := env.StreamUnknown(r.Name, r.Info); err != nil {
			return err
		}
	}
	return nil
}

func (env *Environment) StreamRunning(name string) error {
	msg := "RUNNING: " + name
	return env.StreamMessage(msg, release.TestRunRunning)
}

func (env *Environment) StreamError(info string) error {
	msg := "ERROR: " + info
	return env.StreamMessage(msg, release.TestRunFailure)
}

func (env *Environment) StreamFailed(name string) error {
	msg := fmt.Sprintf("FAILED: %s, run `kubectl logs %s --namespace %s` for more info", name, name, env.Namespace)
	return env.StreamMessage(msg, release.TestRunFailure)
}

func (env *Environment) StreamSuccess(name string) error {
	msg := fmt.Sprintf("PASSED: %s", name)
	return env.StreamMessage(msg, release.TestRunSuccess)
}

func (env *Environment) StreamUnknown(name, info string) error {
	msg := fmt.Sprintf("UNKNOWN: %s: %s", name, info)
	return env.StreamMessage(msg, release.TestRunUnknown)
}

func (env *Environment) StreamMessage(msg string, status release.TestRunStatus) error {
	resp := &release.TestReleaseResponse{Msg: msg, Status: status}
	env.Messages <- resp
	return nil
}

// DeleteTestPods deletes resources given in testManifests
func (env *Environment) DeleteTestPods(testManifests []string) {
	for _, testManifest := range testManifests {
		err := env.KubeClient.Delete(env.Namespace, bytes.NewBufferString(testManifest))
		if err != nil {
			env.StreamError(err.Error())
		}
	}
}
