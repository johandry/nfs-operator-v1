package resources

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	ibmcloudv1alpha1 "github.com/johandry/nfs-operator/pkg/apis/ibmcloud/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ Reconcilable = &ResUnstructured{}

// ResUnstructured is the resource Unstructured
type ResUnstructured struct {
	group     string
	kind      string
	version   string
	name      string
	namespace string
	Object    *unstructured.Unstructured
	Resource
}

// Unstructured create a Unstructured object Resource
func Unstructured(group, kind, version, name, namespace string, object map[string]interface{}, owner *ibmcloudv1alpha1.Nfs, client client.Client, scheme *runtime.Scheme, log logr.Logger) {
	res := &ResUnstructured{
		group:     group,
		kind:      kind,
		version:   version,
		namespace: namespace,
	}
	res.Resource = New(owner, client, scheme, log)
	res.Object = res.newUnstructured(object)

	apiVersion, kind := GVK(res.Object, res.Scheme)
	res.Log = res.Log.WithValues("Resource.Name", res.Object.GetName(), "Resource.Namespace", res.Object.GetNamespace(), "Resource.APIVersion", apiVersion, "Resource.Kind", kind)
}

// Get returns the Object from the cluster
func (r *ResUnstructured) Get() (runtime.Object, error) {
	return r.getUnstructured()
}

// Apply creates the Object if it does not exists
func (r *ResUnstructured) Apply() error {
	_, err := r.getUnstructured()
	exists, err := Exists(err)
	if exists || err != nil {
		return err
	}

	// if not exists and no error, then create
	r.Log.Info("Created a new resource")
	return r.Client.Create(context.TODO(), r.Object)
}

// Reconcile creates the Object if it does not exists and sets the Owner as an
// owner reference on the Object
func (r *ResUnstructured) Reconcile() (reconcile.Result, error) {
	if r.Owner == nil {
		return reconcile.Result{}, fmt.Errorf("the resource %s/%s does not have an owner", r.Object.GetNamespace(), r.Object.GetName())
	}
	if err := controllerutil.SetControllerReference(r.Owner, r.Object, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}
	err := r.Apply()

	return reconcile.Result{}, err
}

// newUnstructured returns the definition of this resource as should exists
func (r *ResUnstructured) newUnstructured(object map[string]interface{}) *unstructured.Unstructured {
	u := &unstructured.Unstructured{}
	u.Object = object

	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   r.group,
		Kind:    r.kind,
		Version: r.version,
	})

	u.SetName(r.name)
	u.SetNamespace(r.namespace)
	return u
}

func (r *ResUnstructured) getUnstructured() (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   r.group,
		Kind:    r.kind,
		Version: r.version,
	})

	err := r.Client.Get(context.Background(), client.ObjectKey{
		Namespace: r.namespace,
		Name:      r.name,
	}, u)

	if err == nil {
		return u, nil
	}
	return nil, err
}
