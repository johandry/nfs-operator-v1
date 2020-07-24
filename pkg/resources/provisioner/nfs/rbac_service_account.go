package nfs

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	ibmcloudv1alpha1 "github.com/johandry/nfs-operator/pkg/apis/ibmcloud/v1alpha1"
	"github.com/johandry/nfs-operator/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ resources.Reconcilable = &ResServiceAccount{}

// ResServiceAccount is the resource ServiceAccount
type ResServiceAccount struct {
	Object *corev1.ServiceAccount
	resources.Resource
}

var contentServiceAccount = []byte(`
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nfs-provisioner
`)

// ServiceAccount creates a ServiceAccount
func ServiceAccount(owner *ibmcloudv1alpha1.Nfs, client client.Client, scheme *runtime.Scheme, log logr.Logger) *ResServiceAccount {
	res := &ResServiceAccount{}
	res.Resource = resources.New(owner, client, scheme, log)
	res.Object = res.newServiceAccount()
	apiVersion, kind := resources.GVK(res.Object, res.Scheme)
	res.Log = res.Log.WithValues("Resource.Name", res.Object.GetName(), "Resource.Namespace", res.Object.GetNamespace(), "Resource.APIVersion", apiVersion, "Resource.Kind", kind)

	return res
}

// Get returns the Object from the cluster
func (r *ResServiceAccount) Get() (runtime.Object, error) {
	return r.getServiceAccount()
}

// Apply creates the Object if it does not exists
func (r *ResServiceAccount) Apply() error {
	_, err := r.getServiceAccount()
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
func (r *ResServiceAccount) Reconcile() (reconcile.Result, error) {
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

// newServiceAccount returns the definition of this resource as should exists
func (r *ResServiceAccount) newServiceAccount() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appName,
			Namespace: r.Owner.Namespace,
		},
	}
}

func (r *ResServiceAccount) getServiceAccount() (*corev1.ServiceAccount, error) {
	found := &corev1.ServiceAccount{}
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
