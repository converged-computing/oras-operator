/*
Copyright 2023 Lawrence Livermore National Security, LLC

(c.f. AUTHORS, NOTICE.LLNS, COPYING)
SPDX-License-Identifier: MIT
*/

package oras

import (
	orasSettings "github.com/converged-computing/oras-operator/pkg/settings"
	corev1 "k8s.io/api/core/v1"
)

// AddSidecar adds a side car container
func AddSidecar(pod *corev1.Pod, settings *orasSettings.OrasCacheSettings) error {

	// Design the sidecar container
	sidecar := corev1.Container{
		Image:   settings.Get("oras-container"),
		Name:    "oras",
		Command: []string{"sh", "-c", "sleep infinity"},
	}
	pod.Spec.Containers = append(pod.Spec.Containers, sidecar)
	return nil

}
