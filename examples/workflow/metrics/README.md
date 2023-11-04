# Experiments with Metrics

This setup will test running different experiments with the metrics operator, and specifically those that generate output files,
and using the ORAS operator to push and save them to a local registry.

## Cluster

First, create a test cluster.

```bash
kind create cluster
```

Install the metrics operator and oras operator:

```bash
# The Metrics Operator requires JobSet
VERSION=v0.2.1
kubectl apply --server-side -f https://github.com/kubernetes-sigs/jobset/releases/download/$VERSION/manifests.yaml
kubectl apply -f https://raw.githubusercontent.com/converged-computing/metrics-operator/main/examples/dist/metrics-operator.yaml

# The oras operator requires cert manager (wait a minute or so for this to be ready)
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.1/cert-manager.yaml

# Wait about a minute...
kubectl apply -f ./examples/dist/oras-operator.yaml
```

## Registry

Create the registry, which is called "oras" - our pods will need to have this annotation to use it.

```bash
kubectl apply -f oras.yaml
```

Now we have the metrics operator and oras operator running in the same namespace! Let's run an experiment that will save several artifacts for us. Since we get to decide the organization and naming of the artifacts, we can take a simple approach of having each named based on the metric (the repository of the registry). I made this manually to keep it simple:

## Metrics

### Creation

```bash
kubectl apply -f data/
```
```console
metricset.flux-framework.org/hwloc-0 created
metricset.flux-framework.org/lammps-0 created
```

Note that we are just running one lammps/hwloc at a time, as they go on the oras network (and would find one another). I had thought we could create more than one headless service for a pod, but I haven't gotten this working yet. 
If you are running experiments you'd likely be running one at a time. If you look at the pods, you should only see two containers for targeted LAMMPS or the main HWLOC one (e.g., the launcher) where we expect to generate output:

```
$ kubectl get pods
NAME                   READY   STATUS    RESTARTS   AGE
hwloc-0-m-0-0-zlf9d    2/2     Running   0          6s
lammps-0-l-0-0-k9ztb   2/2     Running   0          6s
lammps-0-w-0-0-hg49p   1/1     Running   0          6s
oras-0                 1/1     Running   0          83s
```

The others (e.g., the LAMMPS workers) have the network and other pod customization, but will not have the sidecar added.
You can now wait until they are finished:

```
kubectl  get pods
NAME                   READY   STATUS      RESTARTS   AGE
hwloc-0-m-0-0-528kq    0/2     Completed   0          109s
lammps-0-l-0-0-fdckq   0/2     Completed   0          109s
oras-0                 1/1     Running     0          16m
```

#### LAMMPS

Wait until they are completed. You should be able to see the lammps log and the push of the artifact in the logs of the launcher.

```bash
# Here is the launcher (running lammps)
$ kubectl logs lammps-0-l-0-0-7ppcg 
```
```console
...
Expecting: <artifact-input> <artifact-output> <command>...
Full provided set of arguments are NA /opt/lammps/examples/reaxff/HNS/lammps.out /bin/bash /metrics_operator/launcher.sh
Command is /bin/bash /metrics_operator/launcher.sh
Artifact input is NA
Artifact output is /opt/lammps/examples/reaxff/HNS/lammps.out
🟧️  wait-fs: 2023/11/04 16:42:55 wait-fs.go:40: /mnt/oras/oras-operator-init.txt
🟧️  wait-fs: 2023/11/04 16:42:55 wait-fs.go:49: Found existing path /mnt/oras/oras-operator-init.txt
METADATA START {"pods":2,"metricName":"app-lammps","metricDescription":"LAMMPS molecular dynamic simulation","metricOptions":{"command":"mpirun --hostfile ./hostlist.txt -np 4 -ppn 2 lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite 2\u003e\u00261 | tee -a lammps.out","soleTenancy":"false","workdir":"/opt/lammps/examples/reaxff/HNS"}}
METADATA END
Sleeping for 10 seconds waiting for network...
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT
LAMMPS (29 Sep 2021 - Update 2)
OMP_NUM_THREADS environment is not set. Defaulting to 1 thread. (src/comm.cpp:98)
  using 1 OpenMP thread(s) per MPI task
Reading data file ...
  triclinic box = (0.0000000 0.0000000 0.0000000) to (22.326000 11.141200 13.778966) with tilt (0.0000000 -5.0260300 0.0000000)
  2 by 1 by 2 MPI processor grid
  reading atoms ...
  304 atoms
  reading velocities ...
  304 velocities
  read_data CPU = 0.002 seconds
Replicating atoms ...
  triclinic box = (0.0000000 0.0000000 0.0000000) to (44.652000 22.282400 27.557932) with tilt (0.0000000 -10.052060 0.0000000)
  2 by 1 by 2 MPI processor grid
  bounding box image = (0 -1 -1) to (0 1 1)
  bounding box extra memory = 0.03 MB
  average # of replicas added to proc = 5.00 out of 8 (62.50%)
  2432 atoms
  replicate CPU = 0.001 seconds
Neighbor list info ...
  update every 20 steps, delay 0 steps, check no
  max neighbors/atom: 2000, page size: 100000
  master list distance cutoff = 11
  ghost atom cutoff = 11
  binsize = 5.5, bins = 10 5 6
  2 neighbor lists, perpetual/occasional/extra = 2 0 0
  (1) pair reax/c, perpetual
      attributes: half, newton off, ghost
      pair build: half/bin/newtoff/ghost
      stencil: full/ghost/bin/3d
      bin: standard
  (2) fix qeq/reax, perpetual, copy from (1)
      attributes: half, newton off, ghost
      pair build: copy
      stencil: none
      bin: none
Setting up Verlet run ...
  Unit style    : real
  Current step  : 0
  Time step     : 0.1
Per MPI rank memory allocation (min/avg/max) = 103.8 | 103.8 | 103.8 Mbytes
Step Temp PotEng Press E_vdwl E_coul Volume 
       0          300   -113.27833    437.52118   -111.57687   -1.7014647    27418.867 
      10    299.38517   -113.27631    1439.2449   -111.57492   -1.7013814    27418.867 
      20    300.27106   -113.27884    3764.3565   -111.57762   -1.7012246    27418.867 
      30    302.21063   -113.28428     7007.709   -111.58335   -1.7009363    27418.867 
      40    303.52265   -113.28799    9844.8297   -111.58747   -1.7005186    27418.867 
      50    301.87059   -113.28324    9663.0567   -111.58318   -1.7000523    27418.867 
      60    296.67806   -113.26777    7273.8146   -111.56815   -1.6996137    27418.867 
      70    292.19998   -113.25435    5533.6324   -111.55514   -1.6992157    27418.867 
      80    293.58677   -113.25831    5993.3848   -111.55946   -1.6988534    27418.867 
      90    300.62636   -113.27925    7202.8542   -111.58069   -1.6985592    27418.867 
     100    305.38275   -113.29357     10085.75   -111.59518   -1.6983875    27418.867 
Loop time of 12.5347 on 4 procs for 100 steps with 2432 atoms

Performance: 0.069 ns/day, 348.186 hours/ns, 7.978 timesteps/s
85.9% CPU use with 4 MPI tasks x 1 OpenMP threads

MPI task timing breakdown:
Section |  min time  |  avg time  |  max time  |%varavg| %total
---------------------------------------------------------------
Pair    | 4.8613     | 6.4139     | 9.1093     |  64.6 | 51.17
Neigh   | 0.13942    | 0.15179    | 0.18031    |   4.3 |  1.21
Comm    | 0.60957    | 3.3044     | 4.8564     |  89.9 | 26.36
Output  | 0.0017551  | 0.0018499  | 0.0019553  |   0.2 |  0.01
Modify  | 2.6326     | 2.6617     | 2.6746     |   1.1 | 21.23
Other   |            | 0.001108   |            |       |  0.01

Nlocal:        608.000 ave         612 max         604 min
Histogram: 1 0 0 0 0 2 0 0 0 1
Nghost:        5737.25 ave        5744 max        5732 min
Histogram: 1 0 1 0 0 1 0 0 0 1
Neighs:        231539.0 ave      233090 max      229970 min
Histogram: 1 0 0 0 1 1 0 0 0 1

Total # of neighbors = 926155
Ave neighs/atom = 380.82031
Neighbor list builds = 5
Dangerous builds not checked
Total wall time: 0:00:12
METRICS OPERATOR COLLECTION END
```

And here is the oras sidecar that is waiting for the run to finish:

```bash
$ kubectl logs lammps-0-l-0-0-7ppcg -c oras
```
```console
Expecting: <pull-from> <push-to>
Full provided set of arguments are NA oras-0.oras.default.svc.cluster.local:5000/metric/lammps:iter-0
Artifact URI to retrieve is: NA
Artifact URI to push to is: oras-0.oras.default.svc.cluster.local:5000/metric/lammps:iter-0
🟧️  wait-fs: 2023/11/04 16:42:50 wait-fs.go:40: /mnt/oras/oras-operator-done.txt
🟧️  wait-fs: 2023/11/04 16:42:50 wait-fs.go:53: Path /mnt/oras/oras-operator-done.txt does not exist yet, sleeping 5
🟧️  wait-fs: 2023/11/04 16:42:55 wait-fs.go:53: Path /mnt/oras/oras-operator-done.txt does not exist yet, sleeping 5
🟧️  wait-fs: 2023/11/04 16:43:00 wait-fs.go:53: Path /mnt/oras/oras-operator-done.txt does not exist yet, sleeping 5
🟧️  wait-fs: 2023/11/04 16:43:05 wait-fs.go:53: Path /mnt/oras/oras-operator-done.txt does not exist yet, sleeping 5
🟧️  wait-fs: 2023/11/04 16:43:10 wait-fs.go:53: Path /mnt/oras/oras-operator-done.txt does not exist yet, sleeping 5
🟧️  wait-fs: 2023/11/04 16:43:15 wait-fs.go:53: Path /mnt/oras/oras-operator-done.txt does not exist yet, sleeping 5
🟧️  wait-fs: 2023/11/04 16:43:20 wait-fs.go:49: Found existing path /mnt/oras/oras-operator-done.txt
Uploading fff26963dcb1 .
Uploaded  fff26963dcb1 .
Pushed [registry] oras-0.oras.default.svc.cluster.local:5000/metric/lammps:iter-0
Digest: sha256:d01ff185fdc0974ac7ea974f0e5279ead62d270cfb38b57774ad33d9ea25ed33
```


#### HWLOC

The same can be seen for HWLOC. Here is the main log (that generates the architecture xml, etc).

```bash
$ kubectl logs hwloc-0-m-0-0-nh66b 
```
```console
Expecting: <artifact-input> <artifact-output> <command>...
Full provided set of arguments are NA /tmp/analysis /bin/bash /metrics_operator/entrypoint-0.sh
Command is /bin/bash /metrics_operator/entrypoint-0.sh
Artifact input is NA
Artifact output is /tmp/analysis
🟧️  wait-fs: 2023/11/04 17:37:28 wait-fs.go:40: /mnt/oras/oras-operator-init.txt
🟧️  wait-fs: 2023/11/04 17:37:28 wait-fs.go:49: Found existing path /mnt/oras/oras-operator-init.txt
METADATA START {"pods":1,"metricName":"sys-hwloc","metricDescription":"install hwloc for inspecting hardware locality","metricListOptions":{"commands":["mkdir -p /tmp/analysis","lstopo /tmp/analysis/architecture.png","hwloc-ls /tmp/analysis/machine.xml"]}}
METADATA END
METRICS OPERATOR COLLECTION START
mkdir -p /tmp/analysis
METRICS OPERATOR TIMEPOINT
lstopo /tmp/analysis/architecture.png
METRICS OPERATOR TIMEPOINT
hwloc-ls /tmp/analysis/machine.xml
METRICS OPERATOR TIMEPOINT
METRICS OPERATOR COLLECTION END
bin            etc     lib32   metrics_operator         proc          run   tmp
boot           home    lib64   mnt                      product_name  sbin  usr
dev            inputs  libx32  opt                      product_uuid  srv   var
entrypoint.sh  lib     media   oras-run-application.sh  root          sys
```

And here is the ORAS sidecar:

```bash
$ kubectl logs hwloc-0-m-0-0-nh66b -c oras
```
```console
Expecting: <pull-from> <push-to>
Full provided set of arguments are NA oras-0.oras.default.svc.cluster.local:5000/metric/hwloc:iter-0
Artifact URI to retrieve is: NA
Artifact URI to push to is: oras-0.oras.default.svc.cluster.local:5000/metric/hwloc:iter-0
🟧️  wait-fs: 2023/11/04 17:37:22 wait-fs.go:40: /mnt/oras/oras-operator-done.txt
🟧️  wait-fs: 2023/11/04 17:37:22 wait-fs.go:53: Path /mnt/oras/oras-operator-done.txt does not exist yet, sleeping 5
🟧️  wait-fs: 2023/11/04 17:37:27 wait-fs.go:53: Path /mnt/oras/oras-operator-done.txt does not exist yet, sleeping 5
🟧️  wait-fs: 2023/11/04 17:37:32 wait-fs.go:49: Found existing path /mnt/oras/oras-operator-done.txt
Uploading 74bf636ebdde .
Uploaded  74bf636ebdde .
Pushed [registry] oras-0.oras.default.svc.cluster.local:5000/metric/hwloc:iter-0
Digest: sha256:5209373deb3ce18e01943cbee8eb0da2a9f4929e636c85f9e49e85074b441714
```

#### Cleanup

We can now delete our metrics pods - the data is safely cached in the registry!

```bash
$ kubectl delete -f data/
```

#### Download

While they are better ways to do this, we can easily create a port forward to interact with the registry.
In one terminal:

```bash
$ kubectl port-forward oras-0 5000:5000
```
```console
Forwarding from 127.0.0.1:5000 -> 5000
Forwarding from [::1]:5000 -> 5000
Handling connection for 5000
Handling connection for 5000
```

Note that you will need [ORAS Installed](https://oras.land) on your local machine.

```bash
$ oras repo ls localhost:5000/metric
hwloc
lammps
```

There they are! Now let's try listing tags under each. For this simple experiment, we had the tag correspond to the iteration, and we only had one (index 0) for each.
You can imnagine running more complex setups than that.

```
$ oras repo tags localhost:5000/metric/lammps
iter-0
(env) (base) vanessa@vanessa-ThinkPad-T14-Gen-4:~/Desktop/Code/oras-operator/examples/workflow/metrics$ oras repo tags localhost:5000/metric/hwloc
iter-0
```

And now the moment of truth! let's download the data. Note that if you are extracting multiple tags (with files of the same name) you likely want to do this programatically and
into organized directories. If you don't use the Go-based oras client (which is good imho) you can use the [Oras Python](https://github.com/oras-project/oras-py) SDK instead (I @vsoch maintain it).
Let's just dump these into our [data](data) directory:

```bash
cd data
oras pull localhost:5000/metric/lammps:iter-0 --insecure
oras pull localhost:5000/metric/hwloc:iter-0 --insecure
```

And there you have it! The single file for lammps (with output) and the `/tmp/analysis` directory with hwloc output (likely recommended approach to target a directory for >1 file)!

```bash
$ tree .
.
├── analysis
│   ├── architecture.png
│   └── machine.xml
├── hwloc-iter-0.yaml
├── lammps-iter-0.yaml
└── lammps.out

1 directory, 5 files
```

I am so excited about this I can't tell you - take the above and apply it to a workflow? We can run workflows (with different steps) in Kubernetes without needing to mount
some complex storage! This is what I will work on next. <3

## Clean Up

When you are done, clean up.

```bash
kind delete cluster
```