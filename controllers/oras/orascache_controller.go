/*
Copyright 2023 Lawrence Livermore National Security, LLC

(c.f. AUTHORS, NOTICE.LLNS, COPYING)
SPDX-License-Identifier: MIT
*/

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	api "github.com/converged-computing/oras-operator/api/v1alpha1"
	"github.com/go-logr/logr"
)

// OrasCacheReconciler reconciles a OrasCache object
type OrasCacheReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Log        logr.Logger
	RESTClient rest.Interface
	RESTConfig *rest.Config
}

//+kubebuilder:rbac:groups=cache.converged-computing.github.io,resources=orascaches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.converged-computing.github.io,resources=orascaches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.converged-computing.github.io,resources=orascaches/finalizers,verbs=update

//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/log,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/exec,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=batch,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=core,resources=networks,verbs=create;patch
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources="ingresses",verbs=get;list;watch;create;update;patch;delete

//+kubebuilder:rbac:groups="",resources=events,verbs=create;watch;update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete;exec
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;list;watch;create;update;patch;delete;exec

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the OrasCache object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *OrasCacheReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var spec api.OrasCache

	// Keep developer informed what is going on.
	r.Log.Info("📦️ Event received by OrasCache controller!")
	r.Log.Info("Request: ", "req", req)

	r.Log.Info("Spec: ", "spec", spec)
	err := r.Get(ctx, req.NamespacedName, &spec)
	if err != nil {
		r.Log.Info("🟥️ Failed to get OrasCache. Re-running reconcile.")
		return ctrl.Result{Requeue: true}, err
	}

	// Show parameters provided and validate one flux runner
	if !spec.Validate() {
		r.Log.Info("🟥️ Your OrasCache config did not validate.")
		return ctrl.Result{}, nil
	}

	// Ensure the oras cache is deployed
	r.Log.Info("Spec: ", "spec", spec)
	result, err := r.ensureOrasCache(ctx, &spec)
	if err != nil {
		r.Log.Error(err, "🟥️ Issue ensuring OrasCache")
		return result, err
	}

	// By the time we get here we have a Job + pods + config maps!
	// What else do we want to do?
	r.Log.Info("📦️ OrasCache is Ready!")

	return ctrl.Result{}, nil
}

// ensureMetricsSet creates a JobSet and associated configs
func (r *OrasCacheReconciler) ensureOrasCache(
	ctx context.Context,
	spec *api.OrasCache,
) (ctrl.Result, error) {

	// Create headless service for the API to use
	// This must be created before the stateful set
	selector := map[string]string{"oras-name": spec.Name}
	result, err := r.exposeServices(ctx, spec, selector)
	if err != nil {
		return result, err
	}

	// The service running the oras registry is a stateful set
	_, result, _, err = r.getStatefulSet(ctx, spec)
	if err != nil {
		return result, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OrasCacheReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.OrasCache{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Pod{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
