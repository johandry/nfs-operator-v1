package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BackingStorageSpec defines the desired state of the Backing Storage
type BackingStorageSpec struct {
	PvcName      string `json:"pvcName,omitempty"`
	StorageClass string `json:"storageClass,omitempty"`
	StorageSize  string `json:"storageSize,omitempty"`
}

// NfsSpec defines the desired state of Nfs
type NfsSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// +optional
	// +kubebuilder:default=example-nfs
	StorageClass string `json:"storageClass,omitempty"`

	// +optional
	// +kubebuilder:default=example.com/nfs
	ProvisionerAPI string `json:"provisionerAPI,omitempty"`

	// +optional
	BackingStorage BackingStorageSpec `json:"backingStorage,omitempty"`
}

// NfsStatus defines the observed state of Nfs
type NfsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	Capacity   string `json:"capacity,omitempty"`
	AccessMode string `json:"accessMode,omitempty"`
	Status     string `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Nfs is the Schema for the nfs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=nfs,scope=Namespaced
// +kubebuilder:printcolumn:JSONPath=".status.capacity",name=Capacity,type=string
// +kubebuilder:printcolumn:JSONPath=".spec.storageclass",name=StorageClass,type=string
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
