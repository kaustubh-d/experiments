/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"errors"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kdv1 "kb.example.com/rgbcrd/api/v1"
)

// RGBResourceManagerReconciler reconciles a RGBResourceManager object
type RGBResourceManagerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kd.kb.example.com,resources=rgbresourcemanagers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kd.kb.example.com,resources=rgbresourcemanagers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kd.kb.example.com,resources=rgbresourcemanagers/finalizers,verbs=update

// Additional rbac rules so we can manage Pods and deployments.
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/status,verbs=get
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RGBResourceManager object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *RGBResourceManagerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Logic
	// RGB Resource will create resources as per spec and once child resources (pod or deployment)
	// is created and ready, we will mark RBG resource as ready.

	// Get the RGB Resource.
	var rgb_resource kdv1.RGBResourceManager
	err := r.Get(ctx, req.NamespacedName, &rgb_resource)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("Reconciling RGB", "Color", rgb_resource.Spec.Color)

	if rgb_resource.Spec.Kind == kdv1.RGBSupportedKind(kdv1.PodRc) {
		// Managing PODs
		var childPods corev1.PodList
		if err := r.List(ctx, &childPods,
			client.InNamespace(req.Namespace),
			client.MatchingLabels(map[string]string{"app": "rgb"})); err != nil {

			log.Error(err, "unable to list child pods")
			return ctrl.Result{}, err
		}

		count := len(childPods.Items)
		log.Info("Reconciling RGB", "Kind", "Pod", "Count", count)

		if count == int(rgb_resource.Spec.Count) {
			// Final state achieved, mark rgb as ready
			rgb_resource.Status.Result = kdv1.RGBStatus(kdv1.RGBReady)
			return ctrl.Result{}, nil
		}

		// Reconcile to ensure spec
		if count == int(rgb_resource.Spec.Count) {
			// Final state achieved, mark rgb as ready
			return r.markRGBReady(ctx, log, &rgb_resource)
		} else if count < int(rgb_resource.Spec.Count) {
			// Total Pods less then expected, create
			newCntToCreate := int(rgb_resource.Spec.Count) - count
			log.Info("Reconciling RGB", "operation", "create-pod", "count", newCntToCreate)
			for i := 0; i < newCntToCreate; i++ {

				name := rgb_resource.Name + "-" + uuid.New().String()
				d := createPodObj("default", name, "app", "rgb")
				d.Labels["color"] = string(rgb_resource.Spec.Color)
				// Set owner reference
				if err := ctrl.SetControllerReference(&rgb_resource, d, r.Scheme); err != nil {
					return ctrl.Result{}, err
				}
				log.Info("Reconciling RGB", "operation", "create-pod", "Name", name)
				err = r.Create(ctx, d, &client.CreateOptions{})
				if err != nil {
					// Requeue
					log.Info("Reconciling RGB", "operation", "create-pod", "Failed", name)
					return ctrl.Result{}, err
				}
				log.Info("Reconciling RGB", "operation", "create-pod", "Success", name)
			}
		} else {
			newCntToDelete := int(rgb_resource.Spec.Count) - count
			log.Info("Reconciling RGB", "operation", "delete-pod", "count", newCntToDelete)
			for i := 0; i < newCntToDelete; i++ {
				log.Info("Reconciling RGB", "operation", "delete-pod", "Name", childPods.Items[i].Name)
				err = r.Delete(ctx, &childPods.Items[i], &client.DeleteOptions{})
				if err != nil {
					log.Info("Reconciling RGB", "operation", "delete-pod", "Failed", childPods.Items[i].Name)
					return ctrl.Result{}, err
				}
				log.Info("Reconciling RGB", "operation", "delete-pod", "Success", childPods.Items[i].Name)
			}
		}

	} else if rgb_resource.Spec.Kind == kdv1.RGBSupportedKind(kdv1.DeploymentRc) {
		// Managing deployments.
		var childDeployments appsv1.DeploymentList
		if err := r.List(ctx, &childDeployments,
			client.InNamespace(req.Namespace),
			client.MatchingLabels(map[string]string{"app": "rgb"})); err != nil {

			log.Error(err, "unable to list child deployments")
			return ctrl.Result{}, err
		}

		count := len(childDeployments.Items)
		log.Info("Reconciling RGB", "Deployments", count)

		if count == int(rgb_resource.Spec.Count) {
			log.Info("Reconciling RGB", "operation", "update", "rgb-Status", "Ready")
			// Final state achieved, mark rgb as ready
			return r.markRGBReady(ctx, log, &rgb_resource)
		} else if count < int(rgb_resource.Spec.Count) {
			// Total Deployments less then expected, create
			newCntToCreate := int(rgb_resource.Spec.Count) - count
			log.Info("Reconciling RGB", "operation", "create-deployment", "count", newCntToCreate)
			for i := 0; i < newCntToCreate; i++ {

				name := rgb_resource.Name + "-" + uuid.New().String()
				d := createDeploymentObj("default", name, 1, "app", "rgb")
				d.Labels["color"] = string(rgb_resource.Spec.Color)
				// Set owner reference
				if err := ctrl.SetControllerReference(&rgb_resource, d, r.Scheme); err != nil {
					return ctrl.Result{}, err
				}
				log.Info("Reconciling RGB", "operation", "create-deployment", "Name", name)
				err = r.Create(ctx, d, &client.CreateOptions{})
				if err != nil {
					// Requeue
					log.Info("Reconciling RGB", "operation", "create-deployment", "Failed", name)
					return ctrl.Result{}, err
				}
				log.Info("Reconciling RGB", "operation", "create-deployment", "Success", name)
			}

		} else {
			newCntToDelete := int(rgb_resource.Spec.Count) - count
			log.Info("Reconciling RGB", "operation", "delete-deployment", "count", newCntToDelete)
			for i := 0; i < newCntToDelete; i++ {
				log.Info("Reconciling RGB", "operation", "delete-deployment", "Name", childDeployments.Items[i].Name)
				err = r.Delete(ctx, &childDeployments.Items[i], &client.DeleteOptions{})
				if err != nil {
					log.Info("Reconciling RGB", "operation", "delete-deployment", "Failed", childDeployments.Items[i].Name)
					return ctrl.Result{}, err
				}
				log.Info("Reconciling RGB", "operation", "delete-deployment", "Success", childDeployments.Items[i].Name)
			}
		}
		// Reconcile to ensure spec
	} else {
		return ctrl.Result{}, errors.New("unsupported kind in rgb")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RGBResourceManagerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kdv1.RGBResourceManager{}).
		Owns(&corev1.Pod{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}

func createDeploymentObj(namespace string, name string, replicas int32, labelkey string, labelvalue string) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				labelkey: labelvalue,
			},
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					labelkey: labelvalue,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						labelkey: labelvalue,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: "nginx",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	return deployment
}

func createPodObj(namespace string, name string, labelkey string, labelvalue string) *corev1.Pod {
	deployment := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				labelkey: labelvalue,
			},
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  name,
					Image: "nginx",
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							Protocol:      corev1.ProtocolTCP,
							ContainerPort: 80,
						},
					},
				},
			},
		},
	}
	return deployment
}

func (r *RGBResourceManagerReconciler) markRGBReady(ctx context.Context, log logr.Logger, rgb_resource *kdv1.RGBResourceManager) (ctrl.Result, error) {
	log.Info("Reconciling RGB", "operation", "update", "rgb-Status", "Ready")
	// Final state achieved, mark rgb as ready
	rgb_resource.Status.Result = kdv1.RGBStatus(kdv1.RGBReady)
	err := r.Status().Update(ctx, rgb_resource)
	if err != nil {
		log.Info("Reconciling RGB", "operation", "update", "rgb", "Failed")
		return ctrl.Result{}, err
	}
	log.Info("Reconciling RGB", "operation", "Updated", "rgb-Status", "Ready")
	return ctrl.Result{}, nil
}
