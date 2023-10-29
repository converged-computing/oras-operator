/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//	corev1 "k8s.io/api/core/v1"

	api "github.com/converged-computing/oras-operator/api/v1alpha1"
	"github.com/converged-computing/oras-operator/pkg/defaults"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

// getExistingJob gets an existing job that matches our CRD
func (r *OrasCacheReconciler) getExistingStatefulSet(
	ctx context.Context,
	set *api.OrasCache,
) (*appsv1.StatefulSet, error) {

	existing := &appsv1.StatefulSet{}
	err := r.Client.Get(
		ctx,
		types.NamespacedName{
			Name:      set.Name,
			Namespace: set.Namespace,
		},
		existing,
	)
	return existing, err
}

// getStatefulSet retrieves the stateful set (or creates a new one)
func (r *OrasCacheReconciler) getStatefulSet(
	ctx context.Context,
	spec *api.OrasCache,
) (*appsv1.StatefulSet, ctrl.Result, bool, error) {

	// Look for an existing job
	set, err := r.getExistingStatefulSet(ctx, spec)

	// Create a new job if it does not exist
	if err != nil {
		r.Log.Info(
			"✨ Creating a new Oras Cache ✨",
			"Namespace:", spec.Namespace,
			"Name:", spec.Name,
		)

		// Get one JobSet and container specs to create config maps
		set, err = r.createStatefulSet(ctx, spec)

		// We don't create it here, we need configmaps first
		return set, ctrl.Result{}, false, err

	}
	r.Log.Info(
		"🎉 Found existing Oras Cache 🎉",
		"Namespace:", set.Namespace,
		"Name:", set.Name,
	)
	return set, ctrl.Result{}, true, err
}

// createStatefulSet creates the set (after we know it does not exist)
func (r *OrasCacheReconciler) createStatefulSet(
	ctx context.Context,
	spec *api.OrasCache,
) (*appsv1.StatefulSet, error) {
	r.Log.Info(
		"🎉 Creating Oras Cache 🎉",
		"Namespace:", spec.Namespace,
		"Name:", spec.Name,
	)

	// start with one registry for now
	var replicas int32 = 1
	labels := map[string]string{
		defaults.OrasSelectorKey: spec.Namespace,
	}

	set := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.Name,
			Namespace: spec.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{

					// TODO likely we want support for volumes here (for more registry space)
					Volumes: []v1.Volume{},
					Containers: []v1.Container{{
						Name:  "oras",
						Image: spec.Spec.Image,
					}},
					// RestartPolicy defaults to Always
				},
			},
			ServiceName: spec.Name,
			// Default UpdateStrategy is RollingUpdate
		},
	}

	// Controller reference always needs to be set before creation
	ctrl.SetControllerReference(spec, set, r.Scheme)
	err := r.Client.Create(ctx, set)
	if err != nil {
		r.Log.Error(
			err,
			"Failed to create new Oras Cache",
			"Namespace:", set.Namespace,
			"Name:", set.Name,
		)
		return set, err
	}
	return set, nil
}
