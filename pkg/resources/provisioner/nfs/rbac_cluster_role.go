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

// var _ resources.Reconcilable = &ResClusterRole{}
// var _ resources.Watchable = &ResClusterRole{}

// // ResClusterRole is the resource ClusterRole
// type ResClusterRole struct {
// 	object *rbacv1.ClusterRole
// 	client client.Client
// 	scheme *runtime.Scheme
// 	owner  *ibmcloudv1alpha1.Nfs
// 	log    logr.Logger
// }

// var contentClusterRole = []byte(`
// kind: ClusterRole
// apiVersion: rbac.authorization.k8s.io/v1
// metadata:
//   name: nfs-provisioner-runner
// rules:
//   - apiGroups: [""]
//     resources: ["persistentvolumes"]
//     verbs: ["get", "list", "watch", "create", "delete"]
//   - apiGroups: [""]
//     resources: ["persistentvolumeclaims"]
//     verbs: ["get", "list", "watch", "update"]
//   - apiGroups: ["storage.k8s.io"]
//     resources: ["storageclasses"]
//     verbs: ["get", "list", "watch"]
//   - apiGroups: [""]
//     resources: ["events"]
//     verbs: ["create", "update", "patch"]
//   - apiGroups: [""]
//     resources: ["services", "endpoints"]
//     verbs: ["get"]
//   - apiGroups: ["extensions"]
//     resources: ["podsecuritypolicies"]
//     resourceNames: ["nfs-provisioner"]
//     verbs: ["use"]
// `)

// // ClusterRole creates a ClusterRole
// func ClusterRole(owner *ibmcloudv1alpha1.Nfs, client client.Client, scheme *runtime.Scheme, log logr.Logger) *ResClusterRole {
// 	res := &ResClusterRole{
// 		client: client,
// 		scheme: scheme,
// 		owner:  owner,
// 		log:    log,
// 	}
// 	res.Object = res.newClusterRole()
// 	res.log = res.log.WithValues("Resource.Name", res.Object.GetName(), "Resource.Namespace", res.Object.GetNamespace())

// 	return res
// }

// // Object returns the created resource as an API Object
// func (r *ResClusterRole) Object() metav1.Object {
// 	return r.Object
// }

// // EmptyObject returns an empty resource as a Runtime Object
// func (r *ResClusterRole) EmptyObject() runtime.Object {
// 	return &rbacv1.ClusterRole{}
// }

// // Get returns the Object from the cluster
// func (r *ResClusterRole) Get() (runtime.Object, error) {
// 	return r.getClusterRole()
// }

// // Exists returns true if the Object exists in the cluster
// func (r *ResClusterRole) Exists() (bool, error) {
// 	return r.existsClusterRole()
// }

// // Apply creates the Object if it does not exists
// func (r *ResClusterRole) Apply() error {
// 	return r.applyClusterRole()
// }

// // Watch watches the resource for changes, if any will send a Reconcile request
// // the owner
// func (r *ResClusterRole) Watch(c controller.Controller) error {
// 	return c.Watch(
// 		&source.Kind{Type: &rbacv1.ClusterRole{}},
// 		&handler.EnqueueRequestForOwner{
// 			IsController: true,
// 			OwnerType:    &ibmcloudv1alpha1.Nfs{},
// 		})
// }

// // Reconcile creates the Object if it does not exists and sets the Owner as an
// // owner reference on the Object
// func (r *ResClusterRole) Reconcile() (reconcile.Result, error) {
// 	if r.Owner == nil {
// 		return reconcile.Result{}, fmt.Errorf("the resource %s/%s does not have an owner", r.Object.Namespace, r.Object.Name)
// 	}
// 	if err := controllerutil.SetControllerReference(r.Owner, r.Object, r.Scheme); err != nil {
// 		return reconcile.Result{}, err
// 	}
// 	err := r.applyClusterRole()

// 	return reconcile.Result{}, err
// }

// // newClusterRole returns the definition of this resource as should exists
// func (r *ResClusterRole) newClusterRole() *rbacv1.ClusterRole {
// 	return &rbacv1.ClusterRole{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: appName + "-runner",
// 		},
// 		Rules: []rbacv1.PolicyRule{
// 			{
// 				APIGroups: []string{""},
// 				Resources: []string{"persistentvolumes"},
// 				Verbs:     []string{"get", "list", "watch", "create", "delete"},
// 			},
// 			{
// 				APIGroups: []string{""},
// 				Resources: []string{"persistentvolumeclaims"},
// 				Verbs:     []string{"get", "list", "watch", "update"},
// 			},
// 			{
// 				APIGroups: []string{"storage.k8s.io"},
// 				Resources: []string{"storageclasses"},
// 				Verbs:     []string{"get", "list", "watch"},
// 			},
// 			{
// 				APIGroups: []string{""},
// 				Resources: []string{"events"},
// 				Verbs:     []string{"create", "update", "patch"},
// 			},
// 			{
// 				APIGroups: []string{""},
// 				Resources: []string{"services", "endpoints"},
// 				Verbs:     []string{"get"},
// 			},
// 			{
// 				APIGroups:     []string{"extensions"},
// 				Resources:     []string{"podsecuritypolicies"},
// 				ResourceNames: []string{appName},
// 				Verbs:         []string{"use"},
// 			},
// 		},
// 	}
// }

// func (r *ResClusterRole) getClusterRole() (*rbacv1.ClusterRole, error) {
// 	found := &rbacv1.ClusterRole{}
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

// func (r *ResClusterRole) existsClusterRole() (bool, error) {
// 	_, err := r.getClusterRole()

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

// // applyClusterRole creates this resource if does not exists
// func (r *ResClusterRole) applyClusterRole() error {
// 	exists, err := r.existsClusterRole()
// 	if exists || err != nil {
// 		return err
// 	}

// 	// if not exists and no error, then create
// 	return r.Client.Create(context.TODO(), r.Object)
// }
