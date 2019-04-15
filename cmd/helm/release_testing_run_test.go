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
	}}
	runTestCmd(t, tests)
}
