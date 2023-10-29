/*
Copyright 2023 Lawrence Livermore National Security, LLC

(c.f. AUTHORS, NOTICE.LLNS, COPYING)
SPDX-License-Identifier: MIT
*/

package defaults

// Defaults shared across the library

const (
	OrasCachePrefix       = "oras.converged-computing.github.io"
	OrasBaseImage         = "ghcr.io/oras-project/oras:v1.1.0"
	OrasSelectorKey       = "oras-namespace"
	OrasEmptyDirKey       = "oras-share"
	OrasMountPath         = "/mnt/oras"
	OrasEntrypoint        = "https://raw.githubusercontent.com/converged-computing/oras-operator/main/hack/oras-entrypoint.sh"
	ApplicationEntrypoint = "https://raw.githubusercontent.com/converged-computing/oras-operator/main/hack/entrypoint.sh"
	DefaultMissing        = "NA"
)
