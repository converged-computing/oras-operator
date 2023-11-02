# User Guide

Welcome to the ORAS Operator user guide! If you come here, we are assuming you have a cluster
with the ORAS Operator installed and are interested to submit your own [custom resource](custom-resource-definition.md) to create an OrasCache, 
or that someone has already done it for you. If you are a developer wanting to work on new functionality or features, see our [Developer Guides](../development/index.md) instead.

## Usage

### Overview

An OrasCache deploys an OCI Registry as Storage (ORAS) registry. If you are familiar with Docker and container registries, the idea is similar.
Instead of pushing and pulling layers that are assembled into Docker images, we can instead push and pull artifacts. Use cases for
the ORAS operator might include (but are not limited to):

- I am running experiments that save a ton of small files and I want a place to save them to get at the end.
- I am running a workflow that starts with data from a large storage, and I want to persist intermediate workflow files (but not clutter up the original source)
- I am saving data and I don't feel like messing with storage or a local host mount (yuck)

Arguably if this works, it should also work to push to an actual (non cluster-based) OCI registry. This would be for use cases when you want whatever you are doing to be persisted longer (and maybe shared with collaborators or something like that). This case would be non-ephemeral and require credentials, and arguably you could still use oras to move the artifact between the local temporary service and the final registry.

### Install

#### Quick Install

This works best for production Kubernetes clusters, and you can start with creating a Kind cluster:

```bash
kind create cluster
```

and then downloading the latest Metrics Operator yaml config, and applying it.

```bash
kubectl apply -f https://raw.githubusercontent.com/converged-computing/oras-operator/main/examples/dist/oras-operator.yaml
```

Note that from the repository, this config is generated with:

```bash
$ make build-config
```

and then saved to the main branch where you retrieve it from.

#### Helm Install

We optionally provide an install with helm, which you can do either from the charts in the repository:

```bash
$ git clone https://github.com/converged-computing/oras-operator
$ cd oras-operator
$ helm install ./chart
```

Or directly from GitHub packages (an OCI registry):

```
# helm prior to v3.8.0
$ export HELM_EXPERIMENTAL_OCI=1
$ helm pull oci://ghcr.io/converged-computing/oras-operator-helm/chart
```
```console
Pulled: ghcr.io/converged-computing/oras-operator-helm/chart:0.1.0
```

And install!

```bash
$ helm install chart-0.1.0.tgz
```
```console
NAME: oras-operator
LAST DEPLOYED: Fri Mar 24 18:36:18 2023
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

### Annotations

The Oras Operator works by way of deploying an ORAS (OCI Registry as Storage) Registry to a namespace, and then the workflow tool can add annotations to pods to control how artifacts are cached (retrieved and saved for subsequent steps). 
In that most workflow tools understand inputs and outputs and the DAG, this should be feasible to do. Annotations and their defaults include:

| Name | Description | Required | Default |
|------|-------------|----------|---------|
| input-path | The path in the container that any requested archive is expected to be extracted to | false | the working directory of the application container |
| output-path | The output path in the container to save files | false | the working directory of the application container |
| input-uri | The input unique resource identifier for the registry step, including repository, name, and tag | false | NA will be used if not defined, meaning the step has no inputs |
| output-uri | The output unique resource identifier for the registry step, including repository, name, and tag | false | NA will be used if not defined, meaning the step has no outputs |
| oras-cache | The name of the sidecar orchestrator | false | oras |
| oras-container | The container with oras to run for the service | false | ghcr.io/oras-project/oras:v1.1.0 |
| container | The name of the launcher container | false | assumes the first container found requires the launcher |
| entrypoint | The https address of the application entrypoint to wget | false | [entrypoint.sh](https://raw.githubusercontent.com/converged-computing/oras-operator/main/hack/entrypoint.sh) |
| oras-entrypoint | The https address of the oras cache sidecar entrypoint to wget | false | [oras-entrypoint.sh](https://raw.githubusercontent.com/converged-computing/oras-operator/main/hack/oras-entrypoint.sh) |
| debug | Print all discovered settings in the operator log | false | "false" |


There should not be a need to change the oras-cache (sidecar container) unless for some reason you have another container in the pod also called oras. It is exposed for this rare case.

Currently not supported (but will be soon / if needed):

- An ability to save specific (single) files or groups of files. It's much easier to target a directory so we are taking that approach to start.
- A target of the mutating admission webhook for job or jobset instead of pod. The pod target might not scale, but Job has a better chance.
- More than one launcher container in a pod

Note that while the above can be set manually, the expectation is that a workflow tool will do it. For each of the `input-path` and `output-path` we recommend providing
specific files or directories, and note that if one is not set we use the working directory, which (if this is the root of the container) will result in an error.



## Getting Started

You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows). Create the cluster:

```bash
kind create cluster
```

Install the operator (development here) and this include [cert-manager](https://github.com/cert-manager/cert-manager) for webhook certificates:

```sh
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.1/cert-manager.yaml
make test-deploy-recreate

# same as...
make test-deploy
kubectl apply -f examples/dist/oras-operator-dev.yaml
```

See logs:

```bash
kubectl logs -n oras-operator-system oras-operator-controller-manager-ff66845dd-5299h 
```

Then try one of the [examples](https://github.com/converged-computing/oras-operator/tree/main/examples) in the repository.