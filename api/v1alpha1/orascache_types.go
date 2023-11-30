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

	// Names of secrets for the operator
	// +optional
	Secrets Secrets `json:"secrets"`

	// Skip deploying the registry (stateful set) implying all references
	// are for a remote (existing) registry
	// +kubebuilder:default=true
	// +default=true
	// +optional
	Deploy bool `json:"deploy"`
}

type Secrets struct {

	// Secrets for the environment for the ORAS operator sidecar pod to push
	// e.g., oras pull -u username -p password myregistry.io/myimage:latest
	// This should have ORAS_USER and ORAS_PASS
	// +optional
	OrasEnv string `json:"orasEnv"`

	// Secret for the registry REGISTRY_HTTP_SECRET
	// +optional
	RegistryHttp string `json:"registryHttp"`
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
