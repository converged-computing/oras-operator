/*
Copyright 2023 Lawrence Livermore National Security, LLC

(c.f. AUTHORS, NOTICE.LLNS, COPYING)
SPDX-License-Identifier: MIT
*/

package oras

import (
	"github.com/converged-computing/oras-operator/pkg/defaults"
	orasSettings "github.com/converged-computing/oras-operator/pkg/settings"
	corev1 "k8s.io/api/core/v1"
)

// AddSidecar adds a side car container with oras. Example command;
// oras push {{ registry }}/{{ container }} .
// oras push 10.244.0.17:5000/hello-world/repo:tag --plain-http .
// oras push orascache-sample-0.orascache-sample.default.svc.cluster.local:5000/hello-world/repo:tag --plain-http .
func AddSidecar(pod *corev1.Pod, settings *orasSettings.OrasCacheSettings) error {

	// Design the sidecar container
	sidecar := corev1.Container{
		Image:   settings.Get("oras-container"),
		Name:    "oras",
		Command: []string{"sh", "-c", "sleep infinity"},
	}
	pod.Spec.Containers = append(pod.Spec.Containers, sidecar)

	// The selector for the namespaced registry is the namespace
	// We don't technically need this if we are in the same pod
	if pod.Labels == nil {
		pod.Labels = map[string]string{}
	}
	pod.Labels[defaults.OrasSelectorKey] = pod.ObjectMeta.Namespace
	pod.Spec.Subdomain = settings.Get("oras-cache")
	return nil

}
