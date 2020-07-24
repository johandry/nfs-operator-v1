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

var _ resources.Reconcilable = &ResRoleBinding{}

// ResRoleBinding is the resource RoleBinding
type ResRoleBinding struct {
	Object *rbacv1.RoleBinding
	resources.Resource
}

var contentRoleBinding = []byte(`
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: leader-locking-nfs-provisioner
subjects:
  - kind: ServiceAccount
    name: nfs-provisioner
    # replace with namespace where provisioner is deployed
    namespace: default
roleRef:
  kind: Role
  name: leader-locking-nfs-provisioner
  apiGroup: rbac.authorization.k8s.io
`)

// RoleBinding creates a RoleBinding
func RoleBinding(owner *ibmcloudv1alpha1.Nfs, client client.Client, scheme *runtime.Scheme, log logr.Logger) *ResRoleBinding {
	res := &ResRoleBinding{}
	res.Resource = resources.New(owner, client, scheme, log)
	res.Object = res.newRoleBinding()
	apiVersion, kind := resources.GVK(res.Object, res.Scheme)
	res.Log = res.Log.WithValues("Resource.Name", res.Object.GetName(), "Resource.Namespace", res.Object.GetNamespace(), "Resource.APIVersion", apiVersion, "Resource.Kind", kind)

	return res
}

// Get returns the Object from the cluster
func (r *ResRoleBinding) Get() (runtime.Object, error) {
	return r.getRoleBinding()
}

// Apply creates the Object if it does not exists
func (r *ResRoleBinding) Apply() error {
	_, err := r.getRoleBinding()
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
func (r *ResRoleBinding) Reconcile() (reconcile.Result, error) {
	if r.Owner == nil {
		return reconcile.Result{}, fmt.Errorf("the resource %s/%s does not have an owner", r.Object.Namespace, r.Object.Name)
	}
	if err := controllerutil.SetControllerReference(r.Owner, r.Object, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}
	err := r.Apply()

	return reconcile.Result{}, err
}

// newRoleBinding returns the definition of this resource as should exists
func (r *ResRoleBinding) newRoleBinding() *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "leader-locking-" + appName,
			Namespace: r.Owner.Namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      appName,
				Namespace: r.Owner.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     "leader-locking-" + appName,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
}

func (r *ResRoleBinding) getRoleBinding() (*rbacv1.RoleBinding, error) {
	found := &rbacv1.RoleBinding{}
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
