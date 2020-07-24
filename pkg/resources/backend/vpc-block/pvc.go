package vpcblock

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	ibmcloudv1alpha1 "github.com/johandry/nfs-operator/pkg/apis/ibmcloud/v1alpha1"
	"github.com/johandry/nfs-operator/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ resources.Reconcilable = &ResPersistentVolumeClaim{}

// ResPersistentVolumeClaim is the resource PersistentVolumeClaim
type ResPersistentVolumeClaim struct {
	Object *corev1.PersistentVolumeClaim
	resources.Resource
}

var contentPersistentVolumeClaim = []byte(`
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: nfs-block-custom
spec:
  storageClassName: ibmc-vpc-block-general-purpose
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
`)

// PersistentVolumeClaim creates a PersistentVolumeClaim
func PersistentVolumeClaim(owner *ibmcloudv1alpha1.Nfs, client client.Client, scheme *runtime.Scheme, log logr.Logger) *ResPersistentVolumeClaim {
	res := &ResPersistentVolumeClaim{}
	res.Resource = resources.New(owner, client, scheme, log)
	res.Object = res.newPersistentVolumeClaim()
	apiVersion, kind := resources.GVK(res.Object, res.Scheme)
	res.Log = res.Log.WithValues("Resource.Name", res.Object.GetName(), "Resource.Namespace", res.Object.GetNamespace(), "Resource.APIVersion", apiVersion, "Resource.Kind", kind)

	return res
}

// Get returns the Object from the cluster
func (r *ResPersistentVolumeClaim) Get() (runtime.Object, error) {
	return r.getPersistentVolumeClaim()
}

// Apply creates the Object if it does not exists
func (r *ResPersistentVolumeClaim) Apply() error {
	_, err := r.getPersistentVolumeClaim()
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
func (r *ResPersistentVolumeClaim) Reconcile() (reconcile.Result, error) {
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

// newPersistentVolumeClaim returns the definition of this resource as should exists
func (r *ResPersistentVolumeClaim) newPersistentVolumeClaim() *corev1.PersistentVolumeClaim {
	storageClassNameStr := r.Owner.Spec.BackingStorage.StorageClass
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      persistentVolumeClaimName,
			Namespace: r.Owner.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassNameStr,
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(r.Owner.Spec.BackingStorage.StorageSize),
				},
			},
		},
	}
}

func (r *ResPersistentVolumeClaim) getPersistentVolumeClaim() (*corev1.PersistentVolumeClaim, error) {
	found := &corev1.PersistentVolumeClaim{}
	objKey, err := client.ObjectKeyFromObject(r.Object)
	if err != nil {
		return nil, fmt.Errorf("fail to retreive the object key. %s", err)
	}
	// 	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: r.Object.Name, Namespace: r.Owner.Namespace}, found)
	err = r.Client.Get(context.TODO(), objKey, found)
	if err == nil {
		return found, nil
	}
	return nil, err
}
