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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RGBColor describes describes which color is applied to a resource.
// Only one of the following colors may be specified.
// If none of the following colors is specified, the default one
// is Red.
// +kubebuilder:validation:Enum=Red;Green;Blue
type RGBColor string

const (
	RedColor   RGBColor = "Red"
	GreenColor RGBColor = "Green"
	Blue       RGBColor = "Blue"
)

// +kubebuilder:validation:Enum=core;apps
type RGBSupportedGroup string

const (
	CoreGrp string = "core"
	AppsGrp string = "apps"
)

// +kubebuilder:validation:Enum=v1
type RGBSupportedVersion string

const (
	VerV1 string = "v1"
)

// +kubebuilder:validation:Enum=Pod;Deployment
type RGBSupportedKind string

const (
	PodRc        string = "Pod"
	DeploymentRc string = "Deployment"
)

// +kubebuilder:validation:Enum=Initial;Ready
type RGBStatus string

const (
	RGBInitial string = "Initial"
	RGBReady   string = "Ready"
)

// RGBResourceManagerSpec defines the desired state of RGBResourceManager
type RGBResourceManagerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Color that will be applied to created resources by RGBResourceManager.
	Color RGBColor `json:"color,omitempty"`

	Group   RGBSupportedGroup   `json:"group"`
	Version RGBSupportedVersion `json:"version"`
	Kind    RGBSupportedKind    `json:"kind"`

	// Number of instances
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=5
	Count int32 `json:"count"`
}

// RGBResourceManagerStatus defines the observed state of RGBResourceManager
type RGBResourceManagerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// A list of pointers to currently managed resources.
	// +optional
	Active []corev1.ObjectReference `json:"active,omitempty"`

	Result RGBStatus `json:"result"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=rgb

// RGBResourceManager is the Schema for the rgbresourcemanagers API
type RGBResourceManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RGBResourceManagerSpec   `json:"spec,omitempty"`
	Status RGBResourceManagerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RGBResourceManagerList contains a list of RGBResourceManager
type RGBResourceManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RGBResourceManager `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RGBResourceManager{}, &RGBResourceManagerList{})
}
