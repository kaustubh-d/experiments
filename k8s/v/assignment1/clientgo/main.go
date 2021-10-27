package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	deploymentops "kd.example.com/clntgo/deployments"
	podsops "kd.example.com/clntgo/pods"
)

// Program to demonstrate client-go usage.
// List all pods, deployments.
// Create nginx deployment (with 3 pod) and service.
// Update pods in deployment and label 3 pods as Red, Green, Blue

func main() {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	// Default cli option is no config, so in-cluster.
	var listOp, createOp, deleteOp, updateOp bool
	var ns string
	flag.StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig file")
	flag.StringVar(&ns, "namespace", "default", "namespace to operate on")
	flag.BoolVar(&listOp, "list", false, "Perform only list operations.")
	flag.BoolVar(&createOp, "create", false, "Perform only create operations.")
	flag.BoolVar(&updateOp, "update", false, "Perform only update operations.")
	flag.BoolVar(&deleteOp, "delete", false, "Perform only delete operations.")

	flag.Parse()

	var config *rest.Config
	var err error
	if kubeconfig == "" {
		// in-cluster config
		fmt.Println("Using InClusterConfig")
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		fmt.Println()
	} else {
		// bootstrap config
		fmt.Println()
		fmt.Println("Using kubeconfig: ", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println()
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	if listOp {
		deploymentops.ListAll(clientset, ns)
		podsops.ListAll(clientset, ns)
	}

	if createOp {
		deploymentops.Create(clientset, ns, "kd-nginx-rgb", "nginx:1.14.2", 3, "app", "rgb")
	}

	// Lets color the pods in this deployment.
	if updateOp {
		podsops.UpdateRGB(clientset, ns, "app=rgb")
	}

	if deleteOp {
		deploymentops.Delete(clientset, ns, "kd-nginx-rgb")
	}
}
