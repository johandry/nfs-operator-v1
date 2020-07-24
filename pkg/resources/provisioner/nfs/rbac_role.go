package nfs

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	ibmcloudv1alpha1 "github.com/johandry/nfs-operator/pkg/apis/ibmcloud/v1alpha1"
	"github.com/johandry/nfs-operator/pkg/resources"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ resources.Reconcilable = &ResRole{}

// ResRole is the resource Role
type ResRole struct {
	Object *rbacv1.Role
	resources.Resource
}

var contentRole = []byte(`
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: leader-locking-nfs-provisioner
rules:
  - apiGroups: [""]
    resources: ["endpoints"]
    verbs: ["get", "list", "watch", "create", "update", "patch"]
`)

// Role creates a Role
func Role(owner *ibmcloudv1alpha1.Nfs, client client.Client, scheme *runtime.Scheme, log logr.Logger) *ResRole {
	res := &ResRole{}
	res.Resource = resources.New(owner, client, scheme, log)
	res.Object = res.newRole()
	apiVersion, kind := resources.GVK(res.Object, res.Scheme)
	res.Log = res.Log.WithValues("Resource.Name", res.Object.GetName(), "Resource.Namespace", res.Object.GetNamespace(), "Resource.APIVersion", apiVersion, "Resource.Kind", kind)

	return res
}

// Get returns the Object from the cluster
func (r *ResRole) Get() (runtime.Object, error) {
	return r.getRole()
}

// Apply creates the Object if it does not exists
func (r *ResRole) Apply() error {
	_, err := r.getRole()
	exists, err := resources.Exists(err)
	if exists || err != nil {
		return err
	}

	// if not exists and no error, then create
	r.Log.Info("Created a new resource")
	return r.Client.Create(context.TODO(), r.Object)
}

// Reconcile creates the Object if it does not exists and sets the Owner as an
// owner reference on the Object
func (r *ResRole) Reconcile() (reconcile.Result, error) {
	if r.Owner == nil {
		return reconcile.Result{}, fmt.Errorf("the resource %s/%s does not have an owner", r.Object.Namespace, r.Object.Name)
	}
	if err := controllerutil.SetControllerReference(r.Owner, r.Object, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}
	err := r.Apply()

	return reconcile.Result{}, err
}

// newRole returns the definition of this resource as should exists
func (r *ResRole) newRole() *rbacv1.Role {
	return &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "leader-locking-" + appName,
			Namespace: r.Owner.Namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"endpoints"},
				Verbs:     []string{"get", "list", "watch", "create", "update", "patch"},
			},
		},
	}
}

func (r *ResRole) getRole() (*rbacv1.Role, error) {
	found := &rbacv1.Role{}
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
