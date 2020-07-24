package nfs

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	ibmcloudv1alpha1 "github.com/johandry/nfs-operator/pkg/apis/ibmcloud/v1alpha1"
	"github.com/johandry/nfs-operator/pkg/resources"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ resources.Reconcilable = &ResStorageClass{}

// ResStorageClass is the resource StorageClass
type ResStorageClass struct {
	Object *storagev1.StorageClass
	resources.Resource
}

var contentStorageClass = []byte(`
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: ibmcloud-nfs
provisioner: ibmcloud/nfs
mountOptions:
	- vers=4.1
`)

// StorageClass creates a StorageClass
func StorageClass(owner *ibmcloudv1alpha1.Nfs, client client.Client, scheme *runtime.Scheme, log logr.Logger) *ResStorageClass {
	res := &ResStorageClass{}
	res.Resource = resources.New(owner, client, scheme, log)
	res.Object = res.newStorageClass()
	apiVersion, kind := resources.GVK(res.Object, res.Scheme)
	res.Log = res.Log.WithValues("Resource.Name", res.Object.GetName(), "Resource.Namespace", res.Object.GetNamespace(), "Resource.APIVersion", apiVersion, "Resource.Kind", kind)

	return res
}

// Get returns the Object from the cluster
func (r *ResStorageClass) Get() (runtime.Object, error) {
	return r.getStorageClass()
}

// Apply creates the Object if it does not exists
func (r *ResStorageClass) Apply() error {
	_, err := r.getStorageClass()
	exists, err := resources.Exists(err)
	if exists {
		r.Log.Info("Skip reconcile: Resource already exists")
		return nil
	}
	if err != nil {
		r.Log.Error(err, "Failed to reconcile the resource")
		return err
	}

	// if not exists and no error, then create
	r.Log.Info("Created a new resource")
	return r.Client.Create(context.TODO(), r.Object)
}

// Reconcile creates the Object if it does not exists and sets the Owner as an
// owner reference on the Object
func (r *ResStorageClass) Reconcile() (reconcile.Result, error) {
	if r.Owner == nil {
		return reconcile.Result{}, fmt.Errorf("the resource %s/%s does not have an owner", r.Object.Namespace, r.Object.Name)
	}
	r.Log.Info("Reconciling " + r.Object.Name + " resource")
	if err := controllerutil.SetControllerReference(r.Owner, r.Object, r.Scheme); err != nil {
		r.Log.Error(err, "Failed to set controller reference to resource")
		return reconcile.Result{}, err
	}
	err := r.Apply()

	return reconcile.Result{}, err
}

// newStorageClass returns the definition of this resource as should exists
func (r *ResStorageClass) newStorageClass() *storagev1.StorageClass {
	return &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      storageClassName,
			Namespace: r.Owner.Namespace,
		},
		Provisioner: provisionerName,
		MountOptions: []string{
			"vers=4.1",
		},
	}
}

func (r *ResStorageClass) getStorageClass() (*storagev1.StorageClass, error) {
	found := &storagev1.StorageClass{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: r.Object.Name}, found)
	if err == nil {
		return found, nil
	}
	return nil, err
}
