/*
Copyright 2023 Lawrence Livermore National Security, LLC

(c.f. AUTHORS, NOTICE.LLNS, COPYING)
SPDX-License-Identifier: MIT
*/

// This file is not used, but maintained as the original addition of an OrasCache webhook

package v1alpha1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/converged-computing/oras-operator/pkg/defaults"
	"github.com/converged-computing/oras-operator/pkg/oras"
	orasSettings "github.com/converged-computing/oras-operator/pkg/settings"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// IMPORTANT: if you use the controller-runtime builder, it will derive this name automatically from the gvk (kind, version, etc. so find the actual created path in the logs)
// kubectl describe mutatingwebhookconfigurations.admissionregistration.k8s.io
// Disabled so we manually manage it: config/webhook/manifests.yaml
////kubebuilder:webhook:path=/mutate-v1-sidecar,mutating=true,failurePolicy=fail,sideEffects=None,groups=core;batch,resources=pods;jobs,verbs=create,versions=v1,name=morascache.kb.io,admissionReviewVersions=v1

// +kubebuilder:object:generate=true
type SidecarInjector struct {
	Decoder *admission.Decoder
	Cache   *OrasCache
}

func (a *SidecarInjector) Handle(ctx context.Context, req admission.Request) admission.Response {

	// First try for job
	job := &batchv1.Job{}
	err := a.Decoder.Decode(req, job)
	if err != nil {

		// Try for a pod next
		pod := &corev1.Pod{}
		err := a.Decoder.Decode(req, pod)
		if err != nil {
			logger.Error("Admission error.", err)
			return admission.Errored(http.StatusBadRequest, err)
		}

		// If we get here, we decoded a pod
		err = a.InjectPod(pod)
		if err != nil {
			logger.Error("Inject pod error.", err)
			return admission.Errored(http.StatusBadRequest, err)
		}

		// mutate the fields in pod
		marshalledPod, err := json.Marshal(pod)
		if err != nil {
			logger.Error("Marshalling object error.", err)
			return admission.Errored(http.StatusInternalServerError, err)
		}
		logger.Info("Admission pod success.")
		return admission.PatchResponseFromRaw(req.Object.Raw, marshalledPod)
	}

	// If we get here, we found a job
	// If we get here, we decoded a pod
	err = a.InjectJob(job)
	if err != nil {
		logger.Error("Inject job error.", err)
		return admission.Errored(http.StatusBadRequest, err)
	}

	// mutate the fields in pod
	marshalledJob, err := json.Marshal(job)
	if err != nil {
		logger.Error("Marshalling object error.", err)
		return admission.Errored(http.StatusInternalServerError, err)
	}
	logger.Info("Admission job success.")
	return admission.PatchResponseFromRaw(req.Object.Raw, marshalledJob)
}

// Default is the expected entrypoint for a webhook
func (a *SidecarInjector) Default(ctx context.Context, obj runtime.Object) error {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		job, ok := obj.(*batchv1.Job)
		if !ok {
			return fmt.Errorf("expected a Pod or Job but got a %T", obj)
		}
		logger.Info(fmt.Sprintf("Job %s is marked for ORAS registry cache.", job.Name))
		return a.InjectJob(job)
	}
	logger.Info(fmt.Sprintf("Pod %s is marked for ORAS registry cache.", pod.Name))
	return a.InjectPod(pod)
}

// Default hits the default mutating webhook endpoint
func (a *SidecarInjector) InjectPod(pod *corev1.Pod) error {

	// Cut out early if we have no labels
	if pod.Annotations == nil {
		logger.Info(fmt.Sprintf("Pod %s is not marked for oras storage.", pod.Name))
		return nil
	}

	// Parse oras known labels into settings
	settings := orasSettings.NewOrasCacheSettings(pod.Annotations)

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

	// The selector for the namespaced registry is the namespace
	if pod.Labels == nil {
		pod.Labels = map[string]string{}
	}

	// Even pods without say, the launcher, that are marked should have the network added
	pod.Labels[defaults.OrasSelectorKey] = pod.ObjectMeta.Namespace
	oras.AddSidecar(&pod.Spec, pod.ObjectMeta.Namespace, settings)
	logger.Info(fmt.Sprintf("Pod %s is marked for oras storage.", pod.Name))
	return nil
}

// Default hits the default mutating webhook endpoint
func (a *SidecarInjector) InjectJob(job *batchv1.Job) error {

	// Cut out early if we have no labels
	if job.Annotations == nil {
		logger.Info(fmt.Sprintf("Job %s is not marked for oras storage.", job.Name))
		return nil
	}

	// Parse oras known labels into settings
	settings := orasSettings.NewOrasCacheSettings(job.Annotations)

	// Cut out early if no oras identifiers!
	if !settings.MarkedForOras {
		logger.Warnf("Job %s is not marked for oras storage.", job.Name)
		return nil
	}

	// Validate, return error if no good here.
	if !settings.Validate() {
		logger.Warnf("Job %s oras storage did not validate.", job.Name)
		return fmt.Errorf("oras storage was requested but is not valid")
	}

	// Add the sidecar to the podspec of the job
	if job.Spec.Template.Labels == nil {
		job.Spec.Template.Labels = map[string]string{}
	}

	// Add network to spec template so all pods are targeted
	job.Spec.Template.Labels[defaults.OrasSelectorKey] = job.ObjectMeta.Namespace
	oras.AddSidecar(&job.Spec.Template.Spec, job.ObjectMeta.Namespace, settings)
	logger.Info(fmt.Sprintf("Job %s is marked for oras storage.", job.Name))
	return nil
}
