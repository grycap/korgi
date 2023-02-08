/*
Copyright 2023.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	// KorgiJobScheduler status options
	KorgiJobSchedulerRunning   = "RUNNING"
	KorgiJobSchedulerCompleted = "COMPLETED"
	KorgiJobSchedulerFailed    = "FAILED"
)

// KorgiJobSchedulerSpec defines the desired state of KorgiJobScheduler
type KorgiJobSchedulerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of KorgiJobScheduler. Edit korgijobscheduler_types.go to remove/update
	// Foo string `json:"foo,omitempty"`

	KorgiJobName string `json:"korgijobname"`
	GPURecources string `json:"gpu"`
}

// KorgiJobSchedulerStatus defines the observed state of KorgiJobScheduler
type KorgiJobSchedulerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//Status describes the status of the KorgiJob:
	//RUNNING, COMPLETED, FAILED
	Status string `json:"status"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KorgiJobScheduler is the Schema for the korgijobschedulers API
type KorgiJobScheduler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KorgiJobSchedulerSpec   `json:"spec,omitempty"`
	Status KorgiJobSchedulerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KorgiJobSchedulerList contains a list of KorgiJobScheduler
type KorgiJobSchedulerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KorgiJobScheduler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KorgiJobScheduler{}, &KorgiJobSchedulerList{})
}

func (kjs KorgiJobScheduler) GetStatus() string {
	return kjs.Status.Status
}
