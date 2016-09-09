/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

	"github.com/spf13/cobra"

	"k8s.io/helm/pkg/helm"
)

const testDesc = `
run the tests for a release by doing
$ helm test [RELEASE]
`

type testCmd struct {
	name   string
	out    io.Writer
	client helm.Interface
}

func newTestCmd(c helm.Interface, out io.Writer) *cobra.Command {
	t := &testCmd{
		out:    out,
		client: c,
	}

	cmd := &cobra.Command{
		Use:               "test [RELEASE]",
		Short:             "run tests for a release",
		Long:              testDesc,
		PersistentPreRunE: setupConnection,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkArgsLength(1, len(args), "release name"); err != nil {
				return err
			}
			t.name = args[0]
			t.client = ensureHelmClient(t.client)
			return t.run()
		},
	}

	return cmd
}

func (i *testCmd) run() error {

	_, err := i.client.TestRelease(i.name)
	if err != nil {
		return prettyError(err)
	}

	return nil
}
