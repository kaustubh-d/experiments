package deploymentops

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Delete(clientset *kubernetes.Clientset, ns string, name string) {
	err := clientset.AppsV1().Deployments(ns).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Deleted %v deployments in the cluster with **%v** namespace\n", name, ns)

	fmt.Println()
}
