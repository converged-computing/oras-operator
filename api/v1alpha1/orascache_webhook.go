/*
Copyright 2023 Lawrence Livermore National Security, LLC

(c.f. AUTHORS, NOTICE.LLNS, COPYING)
SPDX-License-Identifier: MIT
*/

// This file is not used, but maintained as the original addition of an OrasCache webhook

package v1alpha1

import (
	"context"
	"fmt"

	"github.com/converged-computing/oras-operator/pkg/oras"
	orasSettings "github.com/converged-computing/oras-operator/pkg/settings"
	corev1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// IMPORTANT: the builder will derive this name automatically from the gvk (kind, version, etc. so find the actual created path in the logs)
//  kubectl describe mutatingwebhookconfigurations.admissionregistration.k8s.io
//+kubebuilder:webhook:path=/mutate--v1-pod,mutating=true,failurePolicy=fail,sideEffects=None,groups="",resources=pods,verbs=create;update,versions=v1,name=morascache.kb.io,admissionReviewVersions=v1

type PodInjector struct {
	Cache *OrasCache
}

func (r *OrasCache) SetupWebhookWithManager(mgr ctrl.Manager) error {

	// Add the oras cache to the PodInjector
	injector := &PodInjector{Cache: r}

	return ctrl.NewWebhookManagedBy(mgr).
		For(&corev1.Pod{}).
		WithDefaulter(injector).
		Complete()
}

var _ webhook.CustomDefaulter = &PodInjector{}

// Default is the expected entrypoint for a webhook
func (a *PodInjector) Default(ctx context.Context, obj runtime.Object) error {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("expected a Pod but got a %T", obj)
	}
	return a.InjectPod(pod)
}

// Default hits the default mutating webhook endpoint
func (a *PodInjector) InjectPod(pod *corev1.Pod) error {

	// Cut out early if we have no labels
	if pod.Annotations == nil {
		logger.Info(fmt.Sprintf("Pod %s is not marked for oras storage.", pod.Name))
		return nil
	}

	// Parse oras known labels into settings
	settings := orasSettings.NewOrasCacheSettings(pod)

	// Cut out early if no oras identifiers!
	if !settings.MarkedForOras {
		logger.Warnf("Pod %s is not marked for oras storage.", pod.Name)
		return nil
	}

	// Validate, return error if no good here.
	if !settings.Validate() {
		logger.Warnf("Pod %s oras storage did not validate.", pod.Name)
		return fmt.Errorf("oras storage was requested but is not valid")
	}

	// Add the sidecar to the pod
	oras.AddSidecar(pod, settings)
	logger.Info(fmt.Sprintf("Pod %s is marked for oras storage.", pod.Name))
	return nil
}
