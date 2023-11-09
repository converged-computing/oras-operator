# Hello world Multiple Example

This examples includes two faux steps to:

- [pods-1.yaml](pods-1.yaml): deploys two pods that create two faux workflow artifacts
- [pod-2.yaml](pod-2.yaml): demonstrates annotations that will retrieve greater than one input.

Install!

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.1/cert-manager.yaml
kubectl apply -f https://raw.githubusercontent.com/converged-computing/oras-operator/main/examples/dist/oras-operator.yaml
```

Let's create our registry in the default namespace:

```bash
kubectl apply -f oras.yaml
```
You should see it running as a pod and a service

```bash
kubectl  get pods,svc | grep oras
```

## Step 1 Workflow

The basic logic of the operator is to use annotations to determine when to add a local storage cache.
For example, if you have installed the cert-manager and operator and create the first pod, you'll see the following
in the operator logs:

```bash
kubectl apply -f pods-1.yaml
```
```console
pod/hello-world-1 created
pod/hello-world-2 created
```

You can look at each respective log to see artifacts pushed to `dinosaur/hello-world:one` and `dinosaur/hello-world:two`

```bash
kubectl logs hello-world-1 -c oras -f
kubectl logs hello-world-2 -c oras -f
```

And the main "applications" running:

```
kubectl logs hello-world-1
kubectl logs hello-world-2
```

Delete when that is done.

```bash
kubectl delete -f pods-1.yaml
```

## Step 2 Workflow

Now apply the next workflow step that is going to extract two artifacts (and list the results)

```bash
kubectl apply -f pod-2.yaml
```

You can then look in the pod logs to see that both of the inputs are retrieved.

```console

2023-11-09 00:14:53 (15.0 MB/s) - 'oras-run-application.sh' saved [1534/1534]

Expecting: <artifact-input> <artifact-output> <command>...
Full provided set of arguments are NA NA NA ls inputs
Command is ls inputs
Pipe to is NA
Artifact input is NA
Artifact output is NA
🟧️  wait-fs: 2023/11/09 00:14:53 wait-fs.go:40: /mnt/oras/oras-operator-init.txt
🟧️  wait-fs: 2023/11/09 00:14:53 wait-fs.go:49: Found existing path /mnt/oras/oras-operator-init.txt
hello-world-1.txt
hello-world-2.txt
```

There is no specification to save output, so that's it. These are pods so they will erroneously restart, but you can
imagine a Job completing and cleaning up without issue!

