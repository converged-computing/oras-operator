# oras-operator

Deploy an ORAS registry (cache for workflow or experiment artifacts) as a service.

## Usage

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

### Annotations

## Overview

These are early design notes while the operator is under development.

### Use Cases

- I am running experiments that save a ton of small files and I want a place to save them to get at the end.
- I am running a workflow that starts with data from a large storage, and I want to persist intermediate workflow files (but not clutter up the original source)
- I am saving data and I don't feel like messing with storage or a local host mount (yuck)

Arguably if this works, it should also work to push to an actual (non cluster-based) OCI registry. This would be for use cases when you want whatever you are doing to be persisted longer (and maybe shared with collaborators or something like that). This case would be non-ephemeral and require credentials, and arguably you could still use oras to move the artifact between the local temporary service and the final registry.

### Design

When running workflows, it is common to have intermediate artifacts that need to transition between
steps. The typical strategy is to mount some peristent storage (e.g., cloud object storage)
and then share the space. However, this approach has the drawback of requiring a different solution
per cloud, and requiring a different setup locally. The solution proposed here would work across clouds
(or local) and provide a registry (OCI Registry cache) that can serve this purpose. The benefits of ORAS are:

- An effective / efficient protocol for pushing/pulling artifacts from a namespaced registry
- Support for authentication / permissions (if needed, can be scoped to a namespace)
- Expecting artifacts that range from small to large
- With recent libraries, an ability to live patch an artifact (without retrieving the entire thing)

This means with an OCI registry (and ORAS as the client to handle interaction with it) running alongside
a workflow or set of experiments, we can test the following approach / take the following steps:

1. Create an operator (this one) that creates the ORAS registry to exist in a namespace and be provided via a service.
2. Watch for labels (coming from any Kubernetes workload abstraction that uses pods to run applications) that provide metadata about storage paths and needs.
3. Given one or more pods are detected, the controller will inject a sidecar using a [mutating admission webhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#mutatingadmissionwebhook) that will manage adding the storage. This is a simple interaction that would come down to:
  - A sidecar container with the ORAS client installed and credentials
  - A known identifier for the artifact to pull, and to what path
  - Pulling the artifact to the desired path via a shared empty volume
4. Some logic would need to be added to ensure that the application does not start until this is done (need to think about this, possibly something with readiness).
5. The workflow step would then proceed to run as expected, but with access to the needed assets from previous step(s)
6. When the main application is done, the sidecar would need to (if specified by the metadata) upload the result to an artifact for the next step. 
7. If desired, the user (at the end of the workflow run) can pull any artifacts that are needed to persist or save, and then cleanup.

The sidecar will need to have a way to determine when the main application is done (not sure about that yet). If we are somehow wrapping the main entrypoint while waiting, it could be grabbing the pid of that run and watching via a shared process namespace, but that might be too messy. Note that the pod will not clean up until both are done.


### Development Plan

I am going to start with a simple case of creating one Job with labels (metadata) to indicate creating artifacts for whatever my workflow does, and then retrieving some first step file, running something, and saving it and pulling the result to my local machine.

After that I'll bring in an actual DAG and workflow tool (e.g., Snakemake) and allow the tool to specify the metadata. I have this need already for Kueue and Snakemake (the executor plugin) so will work on that.

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

Then try one of the examples below.

## TODO:

- test out setup with scripts, merge when basics are working
- create docs and automated builds for containers
- test with a simple dag (maybe snakemake kueue executor)

### Hello World

For this hello world example we will create the 
This shows creating (and interacting) with a simple ORAS registry.

Try creating your oras cache:

```bash
$ kubectl  apply -f examples/tests/registry/oras.yaml 
```
```console
orascache.cache.converged-computing.github.io/orascache-sample created
```

This creates the registry. Note that I'll be working on the watch functionality and sidecar injection next. I don't actually know if I can do it, I suspect we will find out!

## License

HPCIC DevTools is distributed under the terms of the MIT license.
All new contributions must be made under this license.

See [LICENSE](https://github.com/converged-computing/cloud-select/blob/main/LICENSE),
[COPYRIGHT](https://github.com/converged-computing/cloud-select/blob/main/COPYRIGHT), and
[NOTICE](https://github.com/converged-computing/cloud-select/blob/main/NOTICE) for details.

SPDX-License-Identifier: (MIT)

LLNL-CODE- 842614
