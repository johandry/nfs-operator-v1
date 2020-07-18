package storage

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApplyFn is a generic function to apply a resource
type ApplyFn func() (string, metav1.Object, error)

// Resource object to create
type Resource struct {
	Obj metav1.Object
	Err error
}

// Resources list of object created
type Resources map[string]*Resource
