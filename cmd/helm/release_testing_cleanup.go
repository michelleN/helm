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
	"io"
	"time"

	"github.com/spf13/cobra"

	"helm.sh/helm/cmd/helm/require"
	"helm.sh/helm/pkg/action"
)

const releaseTestCleanupHelp = `
The cleanup command cleans up all test artifacts for the given release.

The argument this command takes is the name of a deployed release.
The artifacts that will be cleaned up are Kubernetes resources that
contain the following labels:
  helm.sh/release: <release-name>
  helm.sh/test-artifact: <timestamp>
`

func newReleaseTestCleanupCmd(cfg *action.Configuration, out io.Writer) *cobra.Command {
	client := action.NewReleaseTesting(cfg)

	cmd := &cobra.Command{
		Use:   "cleanup [RELEASE]",
		Short: "cleanup test artifacts for a release",
		Long:  releaseTestCleanupHelp,
		Args:  require.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			//TODO: ensure pretty/consistent error output
			return client.Cleanup(args[0])

		},
	}

	f := cmd.Flags()
	f.DurationVar(&client.Timeout, "timeout", 300*time.Second, "time to wait for any individual Kubernetes operation (like Jobs for hooks)")

	return cmd
}
