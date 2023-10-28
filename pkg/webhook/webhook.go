/*
Copyright 2023 Lawrence Livermore National Security, LLC

(c.f. AUTHORS, NOTICE.LLNS, COPYING)
SPDX-License-Identifier: MIT
*/
package webhook

import (
	"context"
	"encoding/json"
	"net/http"

	logf "sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,sideEffects=None,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.kb.io,admissionReviewVersions=v1

// log is for logging in this package.
var log = logf.Log.WithName("orascache-resource")

type PodInjector struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (a *PodInjector) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}
	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// mutate the fields in pod
	log.Info("Mutating pod on admission")
	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}
