# Hello world example

This examples includes two pods that you can apply to a cluster that has the oras-operator installed.

- [pod-with-storage.yaml](pod-with-storage.yaml): has a label that indicates it wants storage setup.
- [pod.yaml](pod.yaml) does not

And a job that will accomplish the same:

- [job.yaml](job.yaml)

Let's create our registry in the default namespace:

```bash
kubectl apply -f oras.yaml
```
You should see it running as a pod and a service

```bash
kubectl  get pods,svc | grep oras
```

## Pod without ORAS

The basic logic of the operator is to use annotations to determine when to add a local storage cache.
For example, if you have installed the cert-manager and operator and create the first pod, you'll see the following
in the operator logs:

```bash
kubectl apply -f pod.yaml
```
```console
{"level":"warn","ts":1698514977.5395513,"caller":"oras/oras.go:31","msg":"Pod pumpkin-pod is not marked for oras storage."}
```


## Pod with ORAS

But for the other pod:

```bash
kubectl delete -f pod.yaml
kubectl apply -f pod-with-storage.yaml
```
```console
{"level":"info","ts":1698515712.4857638,"caller":"oras/settings.go:46","msg":"map[identifier:{true true hello-world} output-path:{false true /workflow/hello-world.txt}]"}
{"level":"info","ts":1698515712.4857695,"caller":"oras/oras.go:42","msg":"Pod pumpkin-pod is marked for oras storage."}
```

You can then look at the logs of each of the containers to see the artifact generating, being saved, and pushed.

## Job with ORAS

Finally, create a job to run Pi. This job shows piping the command into an output file.

```
kubectl apply -f job.yaml
```

## Pull Output

And then create a port forward on your local machine and pull the final thing with oras!


```bash
$ kubectl port-forward orascache-sample-0 5000:5000
Forwarding from 127.0.0.1:5000 -> 5000
Forwarding from [::1]:5000 -> 5000
Handling connection for 5000
Handling connection for 5000
```

```bash
oras pull localhost:5000/dinosaur/hello-world:latest --insecure
oras pull localhost:5000/dinosaur/hello-world:pancakes --insecure
Downloading d2164606501f .
Downloaded  d2164606501f .
Pulled [registry] localhost:5000/dinosaur/hello-world:latest
Digest: sha256:9efa0709ca99b09f68f2ed90a43aaf5feebe69d7158d40fc2025785811f166cb
```