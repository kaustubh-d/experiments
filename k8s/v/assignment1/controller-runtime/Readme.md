# Program to demonstrate controller-runtime usage.

- Reconciles Replicaset resource and updates a label "pod-count" with pods managed by replicaset.
- When replicaset is scaled, updates label "pod-count"
- When label "delete=true" is added to replicaset, deletes the replicaset with this label.

## How to run?

### On Dev machine

```sh
go run main.go 

# Create replicaset
kuebctl create -f k8s/replica-set.yaml 

# Check label
kc get rs --show-labels 
NAME                      DESIRED   CURRENT   READY   AGE   LABELS
kd-nginx-replicaset       3         3         3       19s   pod-count=3

# Scale Replicaset
kc scale --replicas=2 rs/kd-nginx-replicaset

# Check label
kc get rs --show-labels 
NAME                      DESIRED   CURRENT   READY   AGE   LABELS
kd-nginx-replicaset       2         2         2       11m   pod-count=2

# Add label delete=true
kc label rs kd-nginx-replicaset  delete=true

# Check replicaset is deleted by controller app.
 kc get rs --show-labels                     
NAME                      DESIRED   CURRENT   READY   AGE   LABELS

```

Controller runtime app snippet.
```sh
{"level":"info","ts":1635397588.448792,"logger":"controller-runtime.metrics","msg":"metrics server is starting to listen","addr":":8080"}
{"level":"info","ts":1635397588.449545,"msg":"starting metrics server","path":"/metrics"}
{"level":"info","ts":1635397588.449727,"logger":"controller.replicaset","msg":"Starting EventSource","reconciler group":"apps","reconciler kind":"ReplicaSet","source":"kind source: /, Kind="}
{"level":"info","ts":1635397588.4497972,"logger":"controller.replicaset","msg":"Starting EventSource","reconciler group":"apps","reconciler kind":"ReplicaSet","source":"kind source: /, Kind="}
{"level":"info","ts":1635397588.449822,"logger":"controller.replicaset","msg":"Starting Controller","reconciler group":"apps","reconciler kind":"ReplicaSet"}
{"level":"info","ts":1635397588.553189,"logger":"controller.replicaset","msg":"Starting workers","reconciler group":"apps","reconciler kind":"ReplicaSet","worker count":1}

... create replicaset
Reconciling ReplicaSet kd-nginx-replicaset: pod-count 0
Reconciling ReplicaSet kd-nginx-replicaset: pod-count 0
Reconciling ReplicaSet kd-nginx-replicaset: pod-count 1
Reconciling ReplicaSet kd-nginx-replicaset: pod-count 3
Reconciling ReplicaSet kd-nginx-replicaset: pod-count 3

... scale replicas = 2
Reconciling ReplicaSet kd-nginx-replicaset: pod-count 3
Reconciling ReplicaSet kd-nginx-replicaset: pod-count 2
Reconciling ReplicaSet kd-nginx-replicaset: pod-count 2
Reconciling ReplicaSet kd-nginx-replicaset: pod-count 2

... add label to replicaset delete=true
Reconciling ReplicaSet Deleting kd-nginx-replicaset
Reconciling ReplicaSet Deleted kd-nginx-replicaset


```

### Run as POD.

Build image
```sh
docker build -t kaustubhd/controllerruntimeapp:latest . 
```

Push the image
```sh
docker push kaustubhd/controllerruntimeapp:latest
```

Create ClusterRole and ClusterRoleBinding so that POD can hit API server using controller runtime.
```sh
kubectl create -f k8s/role.yaml

kubectl create -f k8s/role_binding.yaml

kubectl get -f k8s/role.yaml -f k8s/role_binding.yaml

NAME                                                          CREATED AT
clusterrole.rbac.authorization.k8s.io/controllerruntime-app   2021-10-28T06:34:00Z

NAME                                                                 ROLE                                AGE
clusterrolebinding.rbac.authorization.k8s.io/controllerruntime-app   ClusterRole/controllerruntime-app   2s

```

Start the POD to run controller app and watch logs.
```sh
kubectl create -f k8s/pod.yaml

kubectl logs -f sample-controllerruntime-app

```

Create Replicate, scale and add label and watch the logs.
```sh
kubectl create -f k8s/replica-set.yaml
kubectl get rs --show-labels 

kubectl scale --replicas=2 rs/kd-nginx-replicaset
kubectl get rs --show-labels 

kubectl label rs kd-nginx-replicaset  delete=true
kubectl get rs --show-labels 

```


## Learnings

Initially used Role and RoleBinding instead of cluster role and hit following errors in controller app within POD.
```
{"level":"info","ts":1635402296.2555816,"logger":"controller.replicaset","msg":"Starting Controller","reconciler group":"apps","reconciler kind":"ReplicaSet"}
E1028 06:24:56.256266       1 reflector.go:138] pkg/mod/k8s.io/client-go@v0.22.2/tools/cache/reflector.go:167: Failed to watch *v1.Pod: failed to list *v1.Pod: pods is forbidden: User "system:serviceaccount:default:default" cannot list resource "pods" in API group "" at the cluster scope
E1028 06:24:56.256268       1 reflector.go:138] pkg/mod/k8s.io/client-go@v0.22.2/tools/cache/reflector.go:167: Failed to watch *v1.ReplicaSet: failed to list *v1.ReplicaSet: replicasets.apps is forbidden: User "system:serviceaccount:default:default" cannot list resource "replicasets" in API group "apps" at the cluster scope
E1028 06:24:57.474477       1 reflector.go:138] pkg/mod/k8s.io/client-go@v0.22.2/tools/cache/reflector.go:167: Failed to watch *v1.ReplicaSet: failed to list *v1.ReplicaSet: replicasets.apps is forbidden: User "system:serviceaccount:default:default" cannot list resource "replicasets" in API group "apps" at the cluster scope
```

How to check if service account has access?
```sh
kubectl auth can-i list pods --as "system:serviceaccount:default:default" --all-namespaces
no
```

Now use ClusterRole and ClusterRoleBinding and check access.
```sh
kubectl auth can-i list pods --as "system:serviceaccount:default:default" --all-namespaces
yes
```

All good, try running controller app POD and above errors will be fixed.