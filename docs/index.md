# The ORAS Operator

Welcome to the ORAS (OCI Registry as Storage) Operator Documentation!

The ORAS Operator is a Kubernetes Cluster [Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
that you can install to your cluster to easily cache artifacts. Use cases might be (but are not limited to):

 - Workflow steps that need to retrieve outputs from previous steps
 - Metrics or other collectors that generate files
 - Single or one-off pods that you want to save assets from

The operator (at a high level) works as follows:

1. The custom resource definition creates an ORAS registry specific to a namespace
2. A mutating admission webhook watches for newly created Pod (this might change to Job)
3. The Pod is updated to have a sidecar container with the ORAS client
4. Annotations from the Pod dicate the behavior for the retrieval/save of artifacts
5. The registry is on a headless service that is used as the cache

The ORAS Operator is currently 🚧️ Under Construction! 🚧️
This is a *converged computing* project that aims
to unite the worlds and technologies typical of cloud computing and
high performance computing.

To get started, check out the links below!
Would you like to request a feature or contribute?
[Open an issue](https://github.com/converged-computing/oras-operator/issues).

```{toctree}
:caption: Getting Started
:maxdepth: 2
getting_started/index.md
development/index.md
```

```{toctree}
:caption: About
:maxdepth: 2
about/index.md
```

<script>
// This is a small hack to populate empty sidebar with an image!
document.addEventListener('DOMContentLoaded', function () {
    var currentNode = document.querySelector('.md-sidebar__scrollwrap');
    currentNode.outerHTML =
	'<div class="md-sidebar__scrollwrap">' +
		'<img style="width:100%" src="_static/images/the-oras-operator.png"/>' +

	'</div>';
}, false);

</script>
