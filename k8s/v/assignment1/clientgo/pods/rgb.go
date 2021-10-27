package podsops

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func UpdateRGB(clientset *kubernetes.Clientset, ns string, label string) {

	// Find pods with given label and add color label one each.
	pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{
		LabelSelector: label,
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	colors := []string{"red", "green", "blue"}
	i := 0
	for _, v := range pods.Items {

		// Add the label
		v.Labels["color"] = colors[i]

		// Update the pod
		_, err = clientset.CoreV1().Pods(ns).Update(context.TODO(), &v, metav1.UpdateOptions{})
		if err != nil {
			panic(err.Error())
		}

		fmt.Printf("Updated pod %v with label color=%v\n", v.Name, colors[i])
		if i == 3 {
			i = 0
		} else {
			i++
		}
	}
	fmt.Println()
}
