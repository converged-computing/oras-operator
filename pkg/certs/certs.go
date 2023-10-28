/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

 SPDX-License-Identifier: MIT
*/

package certs

import (
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/open-policy-agent/cert-controller/pkg/rotator"
)

const (
	webhookServiceName       = "oras-webhook-service"
	mutatingWebhookName      = "mutating-webhook-configuration"
	certificateAuthorityName = "converged-computiong-ca"
	certificateAuthorityOrg  = "LLNL"
	generationDir            = "/tmp/k8s-webhook-server/serving-certs"
)

//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;update
//+kubebuilder:rbac:groups="admissionregistration.k8s.io",resources=mutatingwebhookconfigurations,verbs=get;list;watch;update
//+kubebuilder:rbac:groups="admissionregistration.k8s.io",resources=validatingwebhookconfigurations,verbs=get;list;watch;update

// Create creates certificates for webhooks. We need to do this before creating the controller
func Create(mgr ctrl.Manager, setupFinished chan struct{}) error {

	// DNSName is <service name>.<namespace>.svc
	// TODO get namespace from api
	var dnsName = fmt.Sprintf("%s.%s.svc", webhookServiceName, "default")

	return rotator.AddRotator(mgr, &rotator.CertRotator{
		SecretKey: types.NamespacedName{
			Namespace: "default",
			Name:      webhookServiceName,
		},
		CertDir:        generationDir,
		CAName:         certificateAuthorityName,
		CAOrganization: certificateAuthorityOrg,
		DNSName:        dnsName,
		IsReady:        setupFinished,
		Webhooks: []rotator.WebhookInfo{
			{
				Type: rotator.Mutating,
				Name: mutatingWebhookName,
			},
		},
		RequireLeaderElection: false,
	})
}
