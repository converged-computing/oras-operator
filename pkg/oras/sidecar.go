/*
Copyright 2023 Lawrence Livermore National Security, LLC

(c.f. AUTHORS, NOTICE.LLNS, COPYING)
SPDX-License-Identifier: MIT
*/

package oras

import (
	"fmt"
	"strings"

	"github.com/converged-computing/oras-operator/pkg/defaults"
	orasSettings "github.com/converged-computing/oras-operator/pkg/settings"
	corev1 "k8s.io/api/core/v1"
)

// AddSidecar adds a side car container with oras. Example command;
// oras push {{ registry }}/{{ container }} .
// oras push 10.244.0.17:5000/hello-world/repo:tag --plain-http .
// oras push orascache-sample-0.orascache-sample.default.svc.cluster.local:5000/hello-world/repo:tag --plain-http .
func AddSidecar(pod *corev1.Pod, settings *orasSettings.OrasCacheSettings) error {

	// Oras entrypoint will take as arguments
	cacheName := settings.Get("oras-cache")
	orasEntrypoint := settings.GetOrasEntrypoint(pod)
	logger.Info(orasEntrypoint)

	// Get the volumeMount
	volumeMount := getEmptyDirVolumeMount()

	// Design the sidecar container
	sidecar := corev1.Container{
		Image: settings.Get("oras-container"),
		Name:  "oras",

		// TODO this is sleep, but we will interactively test
		Command:      []string{"sh", "-c", "sleep infinity"},
		VolumeMounts: []corev1.VolumeMount{volumeMount},
		WorkingDir:   defaults.OrasMountPath,
	}
	// The selector for the namespaced registry is the namespace
	if pod.Labels == nil {
		pod.Labels = map[string]string{}
	}
	pod.Labels[defaults.OrasSelectorKey] = pod.ObjectMeta.Namespace
	pod.Spec.Subdomain = cacheName

	// Add volume with emptyDir to the pod
	if pod.Spec.Volumes == nil {
		pod.Spec.Volumes = []corev1.Volume{}
	}
	pod.Spec.Volumes = append(pod.Spec.Volumes, getEmptyDirVolume())

	// Add the emptyDir that will have the new entrypoint to each launcher
	launcher := settings.Get("container")

	// If we have more than one container, launcher is required
	if len(pod.Spec.Containers) > 1 && launcher == "" {
		return fmt.Errorf("A launcher container name %s/container is required for >1 container.", defaults.OrasCachePrefix)
	}

	updatedContainers := []corev1.Container{}
	logger.Infof("Starting container loop")
	for _, container := range pod.Spec.Containers {

		// If launcher defined and this isn't it, skip
		if launcher != "" && container.Name != launcher {
			logger.Infof("Launcher is defined as %s and container name %s does not match, skipping", launcher, container.Name)
			updatedContainers = append(updatedContainers, container)
			continue
		}

		// Add the emptyDir volume
		if container.VolumeMounts == nil {
			container.VolumeMounts = []corev1.VolumeMount{}
		}
		container.VolumeMounts = append(container.VolumeMounts, volumeMount)

		// Assemble the old entrypoint command
		cmd := strings.Join(append(container.Command, container.Args...), " ")
		logger.Infof("Old command is %s", cmd)

		// artifactInput, artifactOutput, original command that is wrapped
		entrypoint := settings.GetApplicationEntrypoint(cmd)

		// We should only be adding this to one container
		container.Command = []string{"sh", "-c", entrypoint}
		container.Args = []string{}
		logger.Infof("Updating container %s", container)
		updatedContainers = append(updatedContainers, container)
	}

	// Add the sidecar at the end
	updatedContainers = append(updatedContainers, sidecar)
	pod.Spec.Containers = updatedContainers
	logger.Infof("Updated containers %s", pod.Spec.Containers)
	return nil
}
