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
	v1 "k8s.io/api/core/v1"
)

// Exit early if we don't have the launcher
func hasLauncher(containers []corev1.Container, launcher string) bool {

	// If launcher not defined, we target all containers in the pod
	if launcher == "" {
		return true
	}
	found := false
	for _, container := range containers {
		if container.Name == launcher {
			found = true
		}
	}
	return found
}

// AddSidecar adds a side car container toa  pod with oras. Example command;
// oras push {{ registry }}/{{ container }} .
// oras push 10.244.0.17:5000/hello-world/repo:tag --plain-http .
// oras push orascache-sample-0.orascache-sample.default.svc.cluster.local:5000/hello-world/repo:tag --plain-http .
func AddSidecar(
	spec *corev1.PodSpec,
	namespace string,
	settings *orasSettings.OrasCacheSettings) error {

	// Oras entrypoint will take as arguments
	cacheName := settings.Get("oras-cache")
	orasEntrypoint := settings.GetOrasEntrypoint(namespace)

	// Add the emptyDir that will have the new entrypoint to each launcher
	launcher := settings.Get("container")

	// Get the volumeMount
	volumeMount := getEmptyDirVolumeMount()

	// Design the sidecar container
	sidecar := corev1.Container{
		Image:        settings.Get("oras-container"),
		Name:         "oras",
		Command:      []string{"sh", "-c", orasEntrypoint},
		VolumeMounts: []corev1.VolumeMount{volumeMount},
		WorkingDir:   defaults.OrasMountPath,
	}
	spec.Subdomain = cacheName

	// If we have a secret to add with credentials, do it here
	// Are we also adding environment variables from a secret?
	orasEnv := settings.Get("oras-env")
	if orasEnv != "" {
		sidecar.EnvFrom = []v1.EnvFromSource{
			{
				SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: orasEnv,
					},
				},
			},
		}
	}

	// Add volume with emptyDir to the pod
	if spec.Volumes == nil {
		spec.Volumes = []corev1.Volume{}
	}
	spec.Volumes = append(spec.Volumes, getEmptyDirVolume())

	// If we have more than one container, launcher is required
	if len(spec.Containers) > 1 && launcher == "" {
		return fmt.Errorf("A launcher container name %s/container is required for >1 container.", defaults.OrasCachePrefix)
	}

	// We want to add the sidecar logic (and emptyVolume) to containers that are targeted
	// This means oras is added to the pods with any matching launcher
	updatedContainers := getUpdatedContainers(spec.Containers, launcher, volumeMount, settings)

	// Add the sidecar at the end ONLY if the targeted container is in the pod
	// We skip adding sidecar to pods entirey that don't have the launcher
	if hasLauncher(spec.Containers, launcher) {
		updatedContainers = append(updatedContainers, sidecar)
	}
	spec.Containers = updatedContainers
	return nil
}

// getUpdatedContainers for the object
func getUpdatedContainers(
	containers []corev1.Container,
	launcher string,
	volumeMount corev1.VolumeMount,
	settings *orasSettings.OrasCacheSettings,
) []corev1.Container {

	updatedContainers := []corev1.Container{}
	for _, container := range containers {

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

		// artifactInput, artifactOutput, original command that is wrapped
		entrypoint := settings.GetApplicationEntrypoint(cmd)

		// We should only be adding this to one container
		container.Command = []string{"sh", "-c", entrypoint}
		container.Args = []string{}
		updatedContainers = append(updatedContainers, container)
	}
	return updatedContainers
}
