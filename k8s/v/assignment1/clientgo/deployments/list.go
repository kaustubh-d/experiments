package deploymentops

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListAll(clientset *kubernetes.Clientset, ns string) {
	deployments, err := clientset.AppsV1().Deployments(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d deployments in the cluster with %v namespace\n", len(deployments.Items), ns)
	for _, v := range deployments.Items {
		fmt.Println("Name =", v.Name, "Labels =", v.Labels)
	}
	fmt.Println()
}
