package nfs

import (
	"context"

	ibmcloudv1alpha1 "github.com/johandry/nfs-operator/pkg/apis/ibmcloud/v1alpha1"
	vpcblockbackend "github.com/johandry/nfs-operator/pkg/resources/backend/vpc-block"
	nfsprovisioner "github.com/johandry/nfs-operator/pkg/resources/provisioner/nfs"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_nfs")

// Add creates a new Nfs Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNfs{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("nfs-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Nfs
	err = c.Watch(&source.Kind{Type: &ibmcloudv1alpha1.Nfs{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Nfs
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &ibmcloudv1alpha1.Nfs{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileNfs implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileNfs{}

// ReconcileNfs reconciles a Nfs object
type ReconcileNfs struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Nfs object and makes changes based on the state read
// and what is in the Nfs.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileNfs) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Nfs")

	// Fetch the Nfs instance
	instance := &ibmcloudv1alpha1.Nfs{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	result, err := vpcblockbackend.New(instance, r.client, r.scheme, log).Reconcile()
	if err != nil {
		return result, err
	}
	result, err = nfsprovisioner.New(instance, r.client, r.scheme, log).Reconcile()
	if err != nil {
		return result, err
	}

	// resources := []resources.Reconcilable{}
	// resources = append(resources, nfsprovisioner.Resources(instance, r.client, r.scheme, log)...)
	// resources = append(resources, vpcblockbackend.Resources(instance, r.client, r.scheme, log)...)

	// for _, resource := range resources {
	// 	result, err := resource.Reconcile()
	// 	if err != nil {
	// 		return result, err
	// 	}
	// }

	// vpcBlockResources := vpcblock.New(r.client, instance.Namespace, instance.Spec.BackingStorage).Apply()
	// if result, err := r.setReference("VPC Block", instance, vpcBlockResources); err != nil {
	// 	return result, err
	// }

	// nfsProvisionerResources := nfsprovisioner.New(r.client, instance.Namespace, instance.Spec, instance.Status).Apply()
	// if result, err := r.setReference("Nfs Provisioner", instance, nfsProvisionerResources); err != nil {
	// 	return result, err
	// }

	return reconcile.Result{}, nil
}

// func (r *ReconcileNfs) reconcileResources(group string, instance *ibmcloudv1alpha1.Nfs, objects []metav1.Object) (reconcile.Result, err error) {
// 	result := reconcile.Result{}

// 	groupLogger := log.WithValues("Resource.GroupName", group)
// 	groupLogger.Info("Reconciling "+group+" resources")

// 	for _, obj := range objects {
// 		resLogger := groupLogger.WithValues("Resource.Name", obj.Name)
// 		resLogger.Info("Reconciling "+obj.Name+" resource")

// 		// Get the Object Key to identify name and namespace required to search it
// 		objKey, err := client.ObjectKeyFromObject(obj)
// 		if err != nil {
// 			resLogger.Error(err, "Failed to retreive the object key")
// 			return result, err
// 		}

// 		if err = p.client.Get(context.TODO(), objKey, found); err == nil { // exists
// 			return fullName, nil, nil
// 		}
// 	}
// 	return reconcile.Result{}, err
// }

// func (r *ReconcileNfs) setReference(name string, instance *ibmcloudv1alpha1.Nfs, resources resources.Resources) (reconcile.Result, error) {
// 	reqLogger := log.WithValues("Resources.Group", name)
// 	reqLogger.Info("Reconciling " + name)

// 	errs := []string{}

// 	// Set Nfs instance as the owner and controller
// 	for name, obj := range resources {
// 		if ok := obj.CanSetControllerReference(); ok {
// 			if err := controllerutil.SetControllerReference(instance, obj.Object, r.scheme); err != nil {
// 				reqLogger.Error(err, "Failed to set controller reference to resource", "Resource.Name", name)
// 				errs = append(errs, fmt.Sprintf("%s: cannot set reference in controler. %s", name, err))
// 			}
// 		}
// 		if obj.Err != nil {
// 			errs = append(errs, fmt.Sprintf("%s: %s", name, obj.Err))
// 			n, gv, k := getGVK(obj.Obj, r.scheme)
// 			reqLogger.Error(obj.Err, "Failed to reconcile the resource", "Resource.Name", n, "Resource.Content", obj.Obj, "Resource.APIVersion", gv, "Resource.Kind", k)
// 			continue
// 		}

// 		if obj.Obj == nil {
// 			reqLogger.Info("Skip reconcile: Resource already exists", "Resource.Name", name)
// 			continue
// 		}

// 		reqLogger.Info("Created a new resource", "Resource.Name", name)

// 	}

// 	if len(errs) == 0 {
// 		return reconcile.Result{}, nil
// 	}

// 	errStr := ""
// 	for n, err := range errs {
// 		errStr = fmt.Sprintf("%s(%d) - %s.", errStr, n, err)
// 	}
// 	err := fmt.Errorf("Failed to reconcile. Errors (%d): %s", len(errs), errStr)

// 	return reconcile.Result{}, err
// }

// func getGVK(o metav1.Object, scheme *runtime.Scheme) (name string, apiVersion string, kind string) {
// 	ro, ok := o.(runtime.Object)
// 	if !ok {
// 		return
// 	}
// 	gvk, err := apiutil.GVKForObject(ro, scheme)
// 	if err != nil {
// 		return
// 	}

// 	name = o.GetName()
// 	apiVersion = gvk.GroupVersion().String()
// 	kind = gvk.Kind

// 	return
// }
