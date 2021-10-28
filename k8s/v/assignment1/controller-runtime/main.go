package main

import (
	"context"
	"fmt"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func init() {
	log.SetLogger(zap.New())
}

func main() {
	entryLog := log.Log.WithName("entrypoint")

	manager, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
	if err != nil {
		entryLog.Error(err, "could not create manager")
		os.Exit(1)
	}

	err = ctrl.
		NewControllerManagedBy(manager). // Create the Controller
		For(&appsv1.ReplicaSet{}).       // ReplicaSet is the Application API
		Owns(&corev1.Pod{}).             // ReplicaSet owns Pods created by it
		Complete(&ReplicaSetReconciler{Client: manager.GetClient()})
	if err != nil {
		entryLog.Error(err, "could not create controller")
		os.Exit(1)
	}

	if err := manager.Start(ctrl.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "could not start manager")
		os.Exit(1)
	}
}

// ReplicaSetReconciler is a simple Controller example implementation.
type ReplicaSetReconciler struct {
	client.Client
}

// Implement the business logic:
// This function will be called when there is a change to a ReplicaSet or a Pod with an OwnerReference
// to a ReplicaSet.
//
// * Read the ReplicaSet
// * Read the Pods
// * Set a Label on the ReplicaSet with the Pod count.
func (a *ReplicaSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// log := log.FromContext(ctx)
	rs := &appsv1.ReplicaSet{}
	err := a.Get(ctx, req.NamespacedName, rs)
	if err != nil {
		// log.Error(err, "**Error: Replicaset get failed...")
		fmt.Println(err, "**Error: Replicaset get failed...")
		return ctrl.Result{}, err
	}

	pods := &corev1.PodList{}
	err = a.List(ctx, pods, client.InNamespace(req.Namespace), client.MatchingLabels(rs.Spec.Template.Labels))
	if err != nil {
		fmt.Println(err, "**Error: pods listing failed...")
		return ctrl.Result{}, err
	}

	fmt.Println("Reconciling ReplicaSet", rs.Name+": pod-count", len(pods.Items))

	// Set the label if it is missing
	if rs.Labels == nil {
		rs.Labels = map[string]string{}
	}
	rs.Labels["pod-count"] = fmt.Sprintf("%v", len(pods.Items))
	err = a.Update(context.TODO(), rs)
	if err != nil {
		fmt.Println(err, "**Error: Update failed...")
		return ctrl.Result{}, err
	}

	// If replicaset has a label of "delete=true" then delete replicaset.
	if val, ok := rs.Labels["delete"]; ok && val == "true" {
		fmt.Println("Reconciling ReplicaSet", "Deleting", rs.Name)
		err = a.Delete(context.TODO(), rs, &client.DeleteOptions{})
		if err != nil {
			fmt.Println(err, "Reconciling ReplicaSet", "Failed to delete", rs.Name)
			return ctrl.Result{}, err
		}
		fmt.Println("Reconciling ReplicaSet", "Deleted", rs.Name)
	}
	return ctrl.Result{}, nil
}
