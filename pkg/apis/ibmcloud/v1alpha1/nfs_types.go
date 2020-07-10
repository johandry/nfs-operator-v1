package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type BackingStorage struct {
	PvcName      string `json:"pvcName,omitempty"`
	StorageClass string `json:"storageClass,omitempty"`
	StorageSize  string `json:"storageSize,omitempty"`
}

// NfsSpec defines the desired state of Nfs
type NfsSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	StorageClass   string         `json:"storageClass,omitempty"`
	ProvisionerAPI string         `json:"provisionerAPI,omitempty"`
	BackingStorage BackingStorage `json:"backingStorage,omitempty"`
}

// NfsStatus defines the observed state of Nfs
type NfsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Nfs is the Schema for the nfs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=nfs,scope=Namespaced
type Nfs struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NfsSpec   `json:"spec,omitempty"`
	Status NfsStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NfsList contains a list of Nfs
type NfsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Nfs `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Nfs{}, &NfsList{})
}
