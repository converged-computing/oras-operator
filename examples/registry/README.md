# Remote Registry Example

This example will walk through deploying the ORAS operator without a registry service. This means that your pods or jobs will
need to pull /push always from a remote registry, which might be desired if you are doing experiments that are bringing up and
down clusters frequently.

First, create a cluster and install the oras operator.

```bash
kind create cluster
```
```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.1/cert-manager.yaml
kubectl apply -f https://raw.githubusercontent.com/converged-computing/oras-operator/main/examples/dist/oras-operator.yaml
```

Next we want to create a secret that will give access to push / pull to a registry (that requires credentials).
For this I would recommend creating a GitHub personal access token with ONLY GitHub packages permission to read and write.
Then, export the variables to the envrionment:

```bash
export ORAS_USER=github-user
export ORAS_PUSH_PASS=xxxxxxxxxxxx
export ORAS_PULL_PASS=xxxxxxxxxxxx
```
Note that the push and pull password would likely be the same. We provide them separately to trigger adding to each of pull or push, respectively.
There are cases when you need a password, for example, to push, but not to pull (and adding one when it isn't needed may lead to error)!
Then we can use kubectl to create the secret directly.

```bash
kubectl create secret generic oras-env --from-literal="ORAS_USER=${ORAS_USER}" --from-literal="ORAS_PUSH_PASS=${ORAS_PASS}" --from-literal="ORAS_PULL_PASS=${ORAS_PASS}"
```

Let's create our registry in the default namespace:

```bash
kubectl apply -f oras.yaml
```

Check the registry. You should NOT see a pod running, but you should see the service.

```bash
kubectl  get pods,svc | grep oras
```

## Pod with ORAS

Next let's create our pod. Note that there is an annotation for a registry to direct to ghcr.io, which is where we have a repository
URI that our oras user and pass has permission to push to.

```bash
kubectl apply -f pod.yaml
```

Get the oras logs to see what is going on, first for the application container, and then for oras:

```bash
kubectl logs turkey-pod -f
kubectl logs turkey-pod -c oras -f
```

For the oras sidecar, you should see a message to indicate your credentials were found, and the artifact pushed!

```console
...
Registry user and password are set for pulling
Registry user and password are set for pushing
🟧️  wait-fs: 2023/11/30 21:49:22 wait-fs.go:40: /mnt/oras/oras-operator-done.txt
🟧️  wait-fs: 2023/11/30 21:49:22 wait-fs.go:53: Path /mnt/oras/oras-operator-done.txt does not exist yet, sleeping 5
🟧️  wait-fs: 2023/11/30 21:49:27 wait-fs.go:53: Path /mnt/oras/oras-operator-done.txt does not exist yet, sleeping 5
🟧️  wait-fs: 2023/11/30 21:49:32 wait-fs.go:49: Found existing path /mnt/oras/oras-operator-done.txt
WARNING! Using --password via the CLI is insecure. Use --password-stdin.
Uploading 18864c9923a4 .
Uploaded  18864c9923a4 .
Pushed [registry] ghcr.io/manbat/metrics-operator-results:test
Digest: sha256:f6a71c2f4b4e4e9e6b4007f4c8c8c8b43e8a892cf47ccee489d02a06ae9d8ea2
```

When you are done, clean up

```bash
kind delete cluster
```