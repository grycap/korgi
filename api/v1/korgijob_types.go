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
	// KorgiJob status options
	KorgiJobPending      = "PENDING"
	KorgiJobRunning      = "RUNNING"
	KorgiJobRescheduling = "RESCHEDULING"
	KorgiJobCompleted    = "COMPLETED"
	KorgiJobFailed       = "FAILED"
)

// KorgiJobSpec defines the desired state of KorgiJob
type KorgiJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of KorgiJob. Edit korgijob_types.go to remove/update
	// Foo string `json:"foo,omitempty"`

	Image   string   `json:"image"`
	Command []string `json:"command"`
}

// KorgiJobStatus defines the observed state of KorgiJob
type KorgiJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//Status describes the status of the KorgiJob:
	//PENDING, RUNNING, RESCHEDULING, COMPLETED, FAILED
	Status string `json:"status"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:JSONPath=".status.status",name="Status",type="string"

// KorgiJob is the Schema for the korgijobs API
type KorgiJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KorgiJobSpec   `json:"spec,omitempty"`
	Status KorgiJobStatus `json:"status,omitempty"`

	//GPUInfo string `json:""`
}

//+kubebuilder:object:root=true

// KorgiJobList contains a list of KorgiJob
type KorgiJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KorgiJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KorgiJob{}, &KorgiJobList{})
}

func (kj KorgiJob) GetStatus() string {
	return kj.Status.Status
}
