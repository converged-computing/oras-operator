/*
Copyright 2023 Lawrence Livermore National Security, LLC

(c.f. AUTHORS, NOTICE.LLNS, COPYING)
SPDX-License-Identifier: MIT
*/

package oras

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

// Default hits the default mutating webhook endpoint
func InjectPod(ctx context.Context, pod *corev1.Pod) error {

	// Cut out early if we have no labels
	if pod.Annotations == nil {
		logger.Info(fmt.Sprintf("Pod %s is not marked for oras storage.", pod.Name))
		return nil
	}

	// Parse oras known labels into settings
	settings := NewOrasCacheSettings(pod)

	// Cut out early if no oras identifiers!
	if !settings.MarkedForOras {
		logger.Warnf("Pod %s is not marked for oras storage.", pod.Name)
		return nil
	}

	// Validate, return error if no good here.
	if !settings.validate() {
		logger.Warnf("Pod %s oras storage did not validate.", pod.Name)
		return fmt.Errorf("oras storage was requested but is not valid")
	}

	// TODO edit pod here to add sidecar!
	logger.Info(fmt.Sprintf("Pod %s is marked for oras storage.", pod.Name))
	return nil
}
