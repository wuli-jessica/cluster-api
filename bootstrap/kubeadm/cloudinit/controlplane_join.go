/*
Copyright 2019 The Kubernetes Authors.

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

package cloudinit

import (
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api/bootstrap/kubeadm/internal/cluster"
)

const (
	controlPlaneJoinCloudInit = `{{.Header}}
{{template "files" .WriteFiles}}
-   path: /tmp/kubeadm-controlplane-join-config.yaml
    owner: root:root
    permissions: '0640'
    content: |
{{.JoinConfiguration | Indent 6}}
runcmd:
{{- template "commands" .PreKubeadmCommands }}
  - 'kubeadm join --config /tmp/kubeadm-controlplane-join-config.yaml'
{{- template "commands" .PostKubeadmCommands }}
{{- template "ntp" .NTP }}
{{- template "users" .Users }}
`
)

// ControlPlaneJoinInput defines context to generate controlplane instance user data for control plane node join.
type ControlPlaneJoinInput struct {
	BaseUserData
	cluster.Certificates

	BootstrapToken    string
	JoinConfiguration string
}

// NewJoinControlPlane returns the user data string to be used on a new control plane instance.
func NewJoinControlPlane(input *ControlPlaneJoinInput) ([]byte, error) {
	input.Header = cloudConfigHeader
	// TODO: Consider validating that the correct certificates exist. It is different for external/stacked etcd
	input.WriteFiles = input.Certificates.AsFiles()
	input.WriteFiles = append(input.WriteFiles, input.AdditionalFiles...)
	userData, err := generate("JoinControlplane", controlPlaneJoinCloudInit, input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate user data for machine joining control plane")
	}

	return userData, err
}
