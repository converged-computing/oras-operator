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
	corev1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// IMPORTANT: the builder will derive this name automatically from the gvk (kind, version, etc. so find the actual created path in the logs)
//  kubectl describe mutatingwebhookconfigurations.admissionregistration.k8s.io
//+kubebuilder:webhook:path=/mutate--v1-pod,mutating=true,failurePolicy=fail,sideEffects=None,groups="",resources=pods,verbs=create;update,versions=v1,name=morascache.kb.io,admissionReviewVersions=v1

var orascachelog = logf.Log.WithName("orascache-resource")

type PodInjector struct{}

func (r *OrasCache) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&corev1.Pod{}).
		WithDefaulter(&PodInjector{}).
		Complete()
}

var _ webhook.CustomDefaulter = &PodInjector{}

// Default is the expected entrypoint for a webhook
func (a *PodInjector) Default(ctx context.Context, obj runtime.Object) error {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("expected a Pod but got a %T", obj)
	}
	return oras.InjectPod(ctx, pod)
}
