# Program to demonstrate client-go usage.

Uses custom "rgb" namespace and all operations are performed within this namespace.

- List all pods, deployments.
- Create nginx deployment (with 3 pod) and service.
- Update pods in deployment and label 3 pods as Red, Green, Blue

## How to run?
### Out of cluster using kube config.

List all pods and deployments, there will be none.
```sh
go run main.go -kubeconfig ~/.kube/config -list

Using kubeconfig:  /Users/kaustubh/.kube/config

There are 0 deployments in the cluster with rgb namespace

There are 0 pods in the cluster with **rgb** namespace
```

Create deployment with nginx image and label app=rgb
```sh
go run main.go -kubeconfig ~/.kube/config -create
Using kubeconfig:  /Users/kaustubh/.kube/config

Created deployment kd-nginx-rgb in the cluster
```

List all pods and deployments, there will be 1 deployment and 3 pods.
```sh
go run main.go -kubeconfig ~/.kube/config -list
Using kubeconfig:  /Users/kaustubh/.kube/config

There are 1 deployments in the cluster with rgb namespace
Name = kd-nginx-rgb Labels = map[]

There are 3 pods in the cluster with **rgb** namespace
Name = kd-nginx-rgb-8458644588-5xcdr Labels = map[app:rgb pod-template-hash:8458644588] Status =  Running
Name = kd-nginx-rgb-8458644588-lq7fs Labels = map[app:rgb pod-template-hash:8458644588] Status =  Running
Name = kd-nginx-rgb-8458644588-np72p Labels = map[app:rgb pod-template-hash:8458644588] Status =  Running
```

Update each pod within created deployment to have one of labels color=[red green blue]
```sh
go run main.go -kubeconfig ~/.kube/config -update

Using kubeconfig:  /Users/kaustubh/.kube/config

There are 3 pods in the cluster
Updated pod kd-nginx-rgb-8458644588-5xcdr with label color=red
Updated pod kd-nginx-rgb-8458644588-lq7fs with label color=green
Updated pod kd-nginx-rgb-8458644588-np72p with label color=blue
```

List all pods and deployments, there will be 1 deployment and 3 pods.
Each pod will have a label added from color=[red green blue]
```sh
go run main.go -kubeconfig ~/.kube/config -list

Using kubeconfig:  /Users/kaustubh/.kube/config

There are 1 deployments in the cluster with rgb namespace
Name = kd-nginx-rgb Labels = map[]

There are 3 pods in the cluster with **rgb** namespace
Name = kd-nginx-rgb-8458644588-5xcdr Labels = map[app:rgb color:red pod-template-hash:8458644588] Status =  Running
Name = kd-nginx-rgb-8458644588-lq7fs Labels = map[app:rgb color:green pod-template-hash:8458644588] Status =  Running
Name = kd-nginx-rgb-8458644588-np72p Labels = map[app:rgb color:blue pod-template-hash:8458644588] Status =  Running
```
See above for example `color:red` label added.