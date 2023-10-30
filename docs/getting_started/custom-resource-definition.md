# OrasCache

> The CRD "Custom Resource Definition" defines an ORAS Cache (registry)

A CRD is a yaml file that you can apply to your cluster (with the ORAS Operator
installed) to ask for an OrasCache to be deployed. Kubernetes has these [custom resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
to make it easy to automate tasks, and in fact this is the goal of an operator!
A Kubernetes operator is conceptually like a human operator that takes your CRD,
looks at the cluster state, and does whatever is necessary to get your cluster state
to match your request. In the case of the ORAS Operator, this means deploying an OCI Registry
as storage (ORAS) in a particular namespace to cache workflow or other artifacts. 
This document describes the spec of our custom resource definition.

## Custom Resource Definition

The operator custom resource defintion (CRD) currently takes just a few fields!

### header

The header should specify the name (you'll need this for your workflow pods) and the namespace (the operator will run in this namespace and can interact with these pods)

```yaml
apiVersion: cache.converged-computing.github.io/v1alpha1
kind: OrasCache
metadata:
  name: orascache-sample
spec:
...
```

In the above, we generate it in the "default" namespace and name it "orascache-sample." Note that all fields described below go under "spec."

### spec

#### image

The image to use for the registry, which defaults to the one deployed by oras `ghcr.io/oras-project/registry:latest`
Here is what that might look like (reproducing the default):

```yaml
spec:
  # We can use all the defaults here (this is a default)
  image: ghcr.io/oras-project/registry:latest
```

#### secret

In the case of deploying more than one registry, the push secret can be customized.

```yaml
spec:
  # We can use all the defaults here (this is a default)
  secret: mysecret
```

Note that this is not supported yet - likely we would want to add custom volumes (for shared storage) between more than one pod in the stateful set. 
For now we just need to save small amounts of data and will add this functionality when needed (and the secret will then be relevant).
