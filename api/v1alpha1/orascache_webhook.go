/*
Copyright 2023 Lawrence Livermore National Security, LLC

(c.f. AUTHORS, NOTICE.LLNS, COPYING)
SPDX-License-Identifier: MIT
*/

// This file is not used, but maintained as the original addition of an OrasCache webhook

package v1alpha1

import (
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var orascachelog = logf.Log.WithName("orascache-resource")

func (r *OrasCache) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-cache-converged-computing-github-io-v1alpha1-orascache,mutating=true,failurePolicy=fail,sideEffects=None,groups="",resources=pods,verbs=create;update,versions=v1,name=morascache.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &OrasCache{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *OrasCache) Default() {
	orascachelog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}
