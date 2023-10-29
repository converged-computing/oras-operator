/*
Copyright 2023 Lawrence Livermore National Security, LLC

(c.f. AUTHORS, NOTICE.LLNS, COPYING)
SPDX-License-Identifier: MIT
*/

package oras

import (
	"github.com/converged-computing/oras-operator/pkg/defaults"
	corev1 "k8s.io/api/core/v1"
)

func getEmptyDirVolume() corev1.Volume {
	return corev1.Volume{
		Name: defaults.OrasEmptyDirKey,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
}

func getEmptyDirVolumeMount() corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      defaults.OrasEmptyDirKey,
		ReadOnly:  false,
		MountPath: defaults.OrasMountPath,
	}
}
