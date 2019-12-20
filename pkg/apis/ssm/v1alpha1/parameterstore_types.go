package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ParameterStoreSpec defines the desired state of ParameterStore
// +k8s:openapi-gen=true
type ParameterStoreSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	ValueFrom ValueFrom `json:"valueFrom"`
}

type ValueFrom struct {
	ParameterStoreRef ParameterStoreRef `json:"parameterStoreRef"`
}

type ParameterStoreRef struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// ParameterStoreStatus defines the observed state of ParameterStore
// +k8s:openapi-gen=true
type ParameterStoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ParameterStore is the Schema for the parameterstores API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type ParameterStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ParameterStoreSpec   `json:"spec,omitempty"`
	Status ParameterStoreStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ParameterStoreList contains a list of ParameterStore
type ParameterStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ParameterStore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ParameterStore{}, &ParameterStoreList{})
}
