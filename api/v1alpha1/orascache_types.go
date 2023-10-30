/*
Copyright 2023 Lawrence Livermore National Security, LLC

(c.f. AUTHORS, NOTICE.LLNS, COPYING)
SPDX-License-Identifier: MIT
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OrasCacheSpec defines the desired state of OrasCache
type OrasCacheSpec struct {

	// Image is the oras registry to deploy
	// +kubebuilder:default="ghcr.io/oras-project/registry:latest"
	// +default="ghcr.io/oras-project/registry:latest"
	// +optional
	Image string `json:"image"`

	// Secret for the registry REGISTRY_HTTP_SECRET
	// +optional
	Secret string `json:"secret"`
}

// OrasCacheStatus defines the observed state of OrasCache
type OrasCacheStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// OrasCache is the Schema for the orascaches API
type OrasCache struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OrasCacheSpec   `json:"spec,omitempty"`
	Status OrasCacheStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OrasCacheList contains a list of OrasCache
type OrasCacheList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OrasCache `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OrasCache{}, &OrasCacheList{})
}

// Validate a requested metricset
func (o *OrasCache) Validate() bool {
	if o.Spec.Image == "" {
		o.Spec.Image = "ghcr.io/oras-project/registry:latest"
	}
	return true
}
