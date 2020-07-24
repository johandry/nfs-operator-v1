package nfs

// import (
// 	"context"
// 	"fmt"

// 	"github.com/go-logr/logr"
// 	ibmcloudv1alpha1 "github.com/johandry/nfs-operator/pkg/apis/ibmcloud/v1alpha1"
// 	"github.com/johandry/nfs-operator/pkg/resources"
// 	rbacv1 "k8s.io/api/rbac/v1"
// 	"k8s.io/apimachinery/pkg/api/errors"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/runtime"
// 	"sigs.k8s.io/controller-runtime/pkg/client"
// 	"sigs.k8s.io/controller-runtime/pkg/controller"
// 	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
// 	"sigs.k8s.io/controller-runtime/pkg/handler"
// 	"sigs.k8s.io/controller-runtime/pkg/reconcile"
// 	"sigs.k8s.io/controller-runtime/pkg/source"
// )

// var _ resources.Reconcilable = &ResClusterRoleBinding{}
// var _ resources.Watchable = &ResClusterRoleBinding{}

// // ResClusterRoleBinding is the resource ClusterRoleBinding
// type ResClusterRoleBinding struct {
// 	object *rbacv1.ClusterRoleBinding
// 	client client.Client
// 	scheme *runtime.Scheme
// 	owner  *ibmcloudv1alpha1.Nfs
// 	log    logr.Logger
// }

// var contentClusterRoleBinding = []byte(`
// kind: ClusterRoleBinding
// apiVersion: rbac.authorization.k8s.io/v1
// metadata:
//   name: run-nfs-provisioner
// subjects:
//   - kind: ServiceAccount
//     name: nfs-provisioner
//      # replace with namespace where provisioner is deployed
//     namespace: default
// roleRef:
//   kind: ClusterRole
//   name: nfs-provisioner-runner
// 	apiGroup: rbac.authorization.k8s.io
// `)

// // ClusterRoleBinding creates a ClusterRoleBinding
// func ClusterRoleBinding(owner *ibmcloudv1alpha1.Nfs, client client.Client, scheme *runtime.Scheme, log logr.Logger) *ResClusterRoleBinding {
// 	res := &ResClusterRoleBinding{
// 		client: client,
// 		scheme: scheme,
// 		owner:  owner,
// 		log:    log,
// 	}
// 	res.Object = res.newClusterRoleBinding()
// 	res.log = res.log.WithValues("Resource.Name", res.Object.GetName(), "Resource.Namespace", res.Object.GetNamespace())

// 	return res
// }

// // Object returns the created resource as an API Object
// func (r *ResClusterRoleBinding) Object() metav1.Object {
// 	return r.Object
// }

// // EmptyObject returns an empty resource as a Runtime Object
// func (r *ResClusterRoleBinding) EmptyObject() runtime.Object {
// 	return &rbacv1.ClusterRoleBinding{}
// }

// // Get returns the Object from the cluster
// func (r *ResClusterRoleBinding) Get() (runtime.Object, error) {
// 	return r.getClusterRoleBinding()
// }

// // Exists returns true if the Object exists in the cluster
// func (r *ResClusterRoleBinding) Exists() (bool, error) {
// 	return r.existsClusterRoleBinding()
// }

// // Apply creates the Object if it does not exists
// func (r *ResClusterRoleBinding) Apply() error {
// 	return r.applyClusterRoleBinding()
// }

// // Watch watches the resource for changes, if any will send a Reconcile request
// // the owner
// func (r *ResClusterRoleBinding) Watch(c controller.Controller) error {
// 	return c.Watch(
// 		&source.Kind{Type: &rbacv1.ClusterRoleBinding{}},
// 		&handler.EnqueueRequestForOwner{
// 			IsController: true,
// 			OwnerType:    &ibmcloudv1alpha1.Nfs{},
// 		})
// }

// // Reconcile creates the Object if it does not exists and sets the Owner as an
// // owner reference on the Object
// func (r *ResClusterRoleBinding) Reconcile() (reconcile.Result, error) {
// 	if r.Owner == nil {
// 		return reconcile.Result{}, fmt.Errorf("the resource %s/%s does not have an owner", r.Object.Namespace, r.Object.Name)
// 	}
// 	if err := controllerutil.SetControllerReference(r.Owner, r.Object, r.Scheme); err != nil {
// 		return reconcile.Result{}, err
// 	}
// 	err := r.applyClusterRoleBinding()

// 	return reconcile.Result{}, err
// }

// // newClusterRoleBinding returns the definition of this resource as should exists
// func (r *ResClusterRoleBinding) newClusterRoleBinding() *rbacv1.ClusterRoleBinding {
// 	return &rbacv1.ClusterRoleBinding{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: "run-" + appName,
// 		},
// 		Subjects: []rbacv1.Subject{
// 			{
// 				Kind:      "ServiceAccount",
// 				Name:      appName,
// 				Namespace: r.Owner.Namespace,
// 			},
// 		},
// 		RoleRef: rbacv1.RoleRef{
// 			Kind:     "ClusterRole",
// 			Name:     appName + "-runner",
// 			APIGroup: "rbac.authorization.k8s.io",
// 		},
// 	}
// }

// func (r *ResClusterRoleBinding) getClusterRoleBinding() (*rbacv1.ClusterRoleBinding, error) {
// 	found := &rbacv1.ClusterRoleBinding{}
// 	objKey, err := client.ObjectKeyFromObject(r.Object)
// 	if err != nil {
// 		return nil, fmt.Errorf("fail to retreive the object key. %s", err)
// 	}
// 	// 	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: r.Object.Name, Namespace: r.Owner.Namespace}, found)
// 	err = r.Client.Get(context.TODO(), objKey, found)
// 	if err == nil {
// 		return found, nil
// 	}
// 	return nil, err
// }

// func (r *ResClusterRoleBinding) existsClusterRoleBinding() (bool, error) {
// 	_, err := r.getClusterRoleBinding()

// 	if err == nil {
// 		// found, no error
// 		return true, nil
// 	}

// 	if errors.IsNotFound(err) {
// 		// not found. No error
// 		return false, nil
// 	}

// 	// unknown, there is an error
// 	return false, err
// }

// // applyClusterRoleBinding creates this resource if does not exists
// func (r *ResClusterRoleBinding) applyClusterRoleBinding() error {
// 	exists, err := r.existsClusterRoleBinding()
// 	if exists || err != nil {
// 		return err
// 	}

// 	// if not exists and no error, then create
// 	return r.Client.Create(context.TODO(), r.Object)
// }
