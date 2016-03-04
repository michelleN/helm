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

package kubectl

// Create uploads a chart to Kubernetes
func (r RealRunner) Create(stdin []byte) ([]byte, error) {
	args := []string{"create", "-f", "-"}

	cmd := command(args...)
	assignStdin(cmd, stdin)

	return cmd.CombinedOutput()
}

// Create returns the commands to kubectl
func (r PrintRunner) Create(stdin []byte) ([]byte, error) {
	args := []string{"create", "-f", "-"}

	cmd := command(args...)
	assignStdin(cmd, stdin)

	return []byte(cmd.String()), nil
}
