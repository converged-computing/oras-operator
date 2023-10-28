# Hello world example

This examples includes two pods that you can apply to a cluster that has the oras-operator installed.

- [pod-with-storage.yaml](pod-with-storage.yaml): has a label that indicates it wants storage setup.
- [pod.yaml](pod.yaml) does not

The basic logic of the operator is to use annotations to determine when to add a local storage cache.
For example, if you have installed the cert-manager and operator and create the first pod, you'll see the following
in the operator logs:

```bash
kubectl apply -f pod.yaml
```
```console
{"level":"warn","ts":1698514977.5395513,"caller":"oras/oras.go:31","msg":"Pod pumpkin-pod is not marked for oras storage."}
```

But for the other pod:

```bash
kubectl delete -f pod.yaml
kubectl apply -f pod-with-storage.yaml
```
```console
{"level":"info","ts":1698515712.4857638,"caller":"oras/settings.go:46","msg":"map[identifier:{true true hello-world} output-path:{false true /workflow/hello-world.txt}]"}
{"level":"info","ts":1698515712.4857695,"caller":"oras/oras.go:42","msg":"Pod pumpkin-pod is marked for oras storage."}
```

We will next be adding the actual functionality for oras, likely with a container that has it installed first.