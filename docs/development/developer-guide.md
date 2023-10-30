# Developer Guide

This developer guide includes complete instructions for setting up a developer
environment.

## Setup

To work on this operator you should:

 - Have a recent version of Go installed (1.20)
 - Have minikube or kind installed

**Important** For minikube, make sure to enable [DNS plugins](https://minikube.sigs.k8s.io/docs/handbook/addons/ingress-dns/).

```bash
$ minikube addons enable ingress
$ minikube addons enable ingress-dns
```

You'll then also want to clone the repository.

```bash
# Clone the source code
$ git clone https://github.com/converged-computing/oras-operator
$ cd metrics-operator
```

## Local Development

Create your cluster.

### 1. Quick Start

Here is a quick start for doing that, making the namespace, and installing the operator.

```console
# Start a minikube cluster
$ minikube start

# OR kind
$ kind create cluster
```

Here is how to build and install the operator - we recommend you build and load into MiniKube with this command:

```bash
$ make deploy-local
$ minikube image load ghcr.io/converged-computing/oras-operator:test
$ kubectl apply -f examples/dist/oras-operator-local.yaml
```

But you can also try the manual steps:

```bash
# Build the operator
$ make

# How to make your manifests
$ make manifests

# And install. This places an executable "bin/kustomize"
$ make install
```

At this point you can apply any of the examples (under "examples") and continue
testing / applying as you see fit!

## Build Images

If you want to build the "production" images - here is how to do that!

```bash
$ make docker-build
$ make docker-push
```

And helm charts.

```bash
$ make helm
```

Note that these are done in CI so you shouldn't need to do anything from the command line.

## Other Developer Commands

### Build Operator Yaml

To generate the CRD to install to a cluster, we've added a `make build-config` command:

```bash
$ make build-config
```

That will generate a yaml to install the operator (with default container image) to a
cluster in `examples/dist`. This file being updated is tested in the PR, so you
should do it before opening.

## Pre-push

I run this before I push to a GitHub branch.

```bash
$ make pre-push
```

We also use pre-commit for Python formatting:

```bash
pip install -r .github/dev-requirements.txt
pre-commit run --all-files
```

## Writing Metric Containers

This section will include instructions for how to write a metrics container.

### General Instructions

Metrics largely have functionality that comes from shared interfaces, such as a `LauncherWorker`
design that has a main node launcher tasks, and some number of worker nodes, and basic interfaces
for storage and applications. The best thing to do is explore the current metrics, find one that
is similar to what you want to do, and use it as a template. As long as you put it in a known group
directory, e.g., these:

```bash
pkg/metrics/
├── app
├── io
├── network
└── perf
```

It will be discovered and registered and available for use.

## Documentation

The documentation is provided in the `docs` folder of the repository, and generally most content that you might want to add is under `getting_started`. For ease of contribution, files that are likely to be updated by contributors (e.g., mostly everything but the module generated files)
are written in markdown. If you need to use [toctree](https://www.sphinx-doc.org/en/master/usage/restructuredtext/directives.html#table-of-contents) you should not use extra newlines or spaces (see index.md files for examples). The documentation is also provided in Markdown (instead of rst or restructured syntax) to make contribution easier for the community.

### Install Dependencies and Build

The documentation is built using sphinx, and generally you can
create a virtual environment:

```bash
$ cd docs
$ python -m venv env
$ source env/bin/activate
```
And then install dependencies:

```console
$ pip install -r requirements.txt

# Build the docs into _build/html
$ make html
```

### Preview Documentation

After `make html` you can enter into `_build/html` and start a local web
server to preview:

```console
$ python -m http.server 9999
```

And open your browser to `localhost:9999`
