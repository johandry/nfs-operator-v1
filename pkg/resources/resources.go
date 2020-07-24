package resources

import (
	"github.com/go-logr/logr"
	ibmcloudv1alpha1 "github.com/johandry/nfs-operator/pkg/apis/ibmcloud/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconcilable is a resource that can be reconciled by the controller
type Reconcilable interface {
	Get() (runtime.Object, error)
	Apply() error
	Reconcile() (reconcile.Result, error)
}

// Resource is the resource Unstructured
type Resource struct {
	Client client.Client
	Scheme *runtime.Scheme
	Owner  *ibmcloudv1alpha1.Nfs
	Log    logr.Logger
}

// New creates a Resource which can Reconcile
func New(owner *ibmcloudv1alpha1.Nfs, client client.Client, scheme *runtime.Scheme, log logr.Logger) Resource {
	res := Resource{
		Client: client,
		Scheme: scheme,
		Owner:  owner,
		Log:    log,
	}

	return res
}

// GVK returns the API Version and Kind of a given resource
func GVK(ro runtime.Object, scheme *runtime.Scheme) (apiVersion string, kind string) {
	gvk, err := apiutil.GVKForObject(ro, scheme)
	if err != nil {
		return
	}

	apiVersion = gvk.GroupVersion().String()
	kind = gvk.Kind

	return
}

// Exists return true if the given error from getting the resource is not a
// NotFound error. Otherwise returns false or the same error if it's an unknown
// error
func Exists(err error) (bool, error) {
	if err == nil {
		// found, no error
		return true, nil
	}

	if errors.IsNotFound(err) {
		// not found. No error
		return false, nil
	}

	// unknown, there is an error
	return false, err
}
