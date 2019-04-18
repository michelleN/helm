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

package main

import (
	"testing"
	"time"

	"helm.sh/helm/pkg/release"
)

func TestReleaseTestingRun(t *testing.T) {
	// create a test with a test environment interface
	//  so I can control what output I am displaying

	timestamp := time.Unix(1452902400, 0).UTC()

	tests := []cmdTestCase{{
		name: "successful test",
		cmd:  "test run test-success",
		rels: []*release.Release{release.Mock(&release.MockReleaseOptions{
			Name: "test-success",
			TestSuiteResults: []*release.TestRun{
				{
					Name:        "test-success",
					Status:      release.TestRunSuccess,
					StartedAt:   timestamp,
					CompletedAt: timestamp,
				},
			},
		})},
		golden: "output/test-run-success.txt",
	}, {
		name:      "test error in release name",
		cmd:       "test run test-invalid-release-name-error-please-because-this-name-is-too-long",
		wantError: true,
		golden:    "output/test-run-fail-invalid-release-name.txt",
	}}
	runTestCmd(t, tests)
}
