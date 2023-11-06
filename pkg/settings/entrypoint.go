/*
Copyright 2023 Lawrence Livermore National Security, LLC

(c.f. AUTHORS, NOTICE.LLNS, COPYING)
SPDX-License-Identifier: MIT
*/

package settings

import (
	"fmt"
)

const (
	updates = "apt-get update && apt-get install -y wget bash || apk add wget bash || yum install -y wget bash &&"
)

// GetOrasEntrypoint will derive the entrypoint for the sidecar
func (s *OrasCacheSettings) GetOrasEntrypoint(namespace string) string {

	orasScript := s.Get("oras-entrypoint")
	cacheName := s.Get("oras-cache")

	// This is a stateful set so we assume always index 0. Assume same port for now
	registry := fmt.Sprintf("%s-0.%s.%s.svc.cluster.local:5000", cacheName, cacheName, namespace)
	pullFromURI := s.Get("input-uri")
	pushToURI := s.Get("output-uri")

	// Unique name for script
	n := "oras-run-cache.sh"

	// Assemble pull to and from
	var pullFrom, pushTo string = "NA", "NA"
	if pullFromURI != "NA" {
		pullFrom = fmt.Sprintf("%s/%s", registry, pullFromURI)
	}
	if pushToURI != "NA" {
		pushTo = fmt.Sprintf("%s/%s", registry, pushToURI)
	}

	// Ensure we have wget
	orasEntrypoint := fmt.Sprintf("%s wget --no-cache -O %s %s && chmod +x %s && ./%s %s %s", updates, n, orasScript, n, n, pullFrom, pushTo)
	logger.Infof("Oras entrypoint: %s\n", orasEntrypoint)
	return orasEntrypoint

}

func (s *OrasCacheSettings) GetApplicationEntrypoint(cmd string) string {
	script := s.Get("entrypoint") // Application entrypoint
	artifactInput := s.Get("input-path")
	artifactOutput := s.Get("output-path")
	artifactPipe := s.Get("output-pipe")

	// Try to go for a unique name that won't clobber something else
	n := "oras-run-application.sh"

	// wget the new script to run
	cmd = fmt.Sprintf("%s %s %s %s", artifactPipe, artifactInput, artifactOutput, cmd)
	return fmt.Sprintf("%s wget --no-cache -O %s %s && chmod +x %s && ./%s %s", updates, n, script, n, n, cmd)
}
