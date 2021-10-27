package podsops

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListAll(clientset *kubernetes.Clientset, ns string) {
	pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster with **%v** namespace\n", len(pods.Items), ns)
	for _, v := range pods.Items {
		fmt.Println("Name =", v.Name, "Labels =", v.Labels, "Status = ", v.Status.Phase)
	}
	fmt.Println()
}
