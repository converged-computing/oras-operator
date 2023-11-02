# Current Design

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

## Development Plan

I am going to start with a simple case of creating one Job with labels (metadata) to indicate creating artifacts for whatever my workflow does, and then retrieving some first step file, running something, and saving it and pulling the result to my local machine.

After that I'll bring in an actual DAG and workflow tool (e.g., Snakemake) and allow the tool to specify the metadata. I have this need already for Kueue and Snakemake (the executor plugin) so will work on that.