package nfs

import (
	"context"
	"fmt"
	"strings"

	ibmcloudv1alpha1 "github.com/johandry/nfs-operator/pkg/apis/ibmcloud/v1alpha1"
	"github.com/johandry/nfs-operator/pkg/controller/storage"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_nfs")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

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
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
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

	objList := storage.NewNFSProvisioner(r.client, instance.Namespace).Apply()

	errs := []string{}    // list of errors
	created := []string{} // list of objects created
	skipped := []string{} // list of objects that already exists

	// Set Nfs instance as the owner and controller
	for name, obj := range objList {
		if obj.Err != nil {
			errs = append(errs, fmt.Sprintf("%s: %s", name, obj.Err))
			// reqLogger.Error(fmt.Errorf("failed to reconcile the resource %q, Namespace: %s", name, instance.Namespace))
			continue
		} else if obj.Obj == nil {
			skipped = append(skipped, name)
			reqLogger.Info("Skip reconcile: Resource already exists", "Namespace", instance.Namespace, "Name", name)
			continue
		}
		created = append(created, name)
		reqLogger.Info("Created a new resource", "Namespace", obj.Obj.GetNamespace(), "Name", name)

		if err := controllerutil.SetControllerReference(instance, obj.Obj, r.scheme); err != nil {
			errs = append(errs, fmt.Sprintf("%s: cannot set reference in controler. %s", name, err))
		}
	}

	if len(skipped) == 0 {
		reqLogger.Info("Skip reconcile for: Resources already exists", "Namespace", instance.Namespace, "Resources", strings.Join(skipped, ", "))
	}
	if len(created) == 0 {
		reqLogger.Info("Resources created", "Namespace", request.Namespace, "Resources", strings.Join(created, ", "))
	}

	if len(errs) == 0 {
		return reconcile.Result{}, nil
	}

	errStr := ""
	for _, err := range errs {
		errStr = fmt.Sprintf("%s\n\t%s", errStr, err)
	}
	err = fmt.Errorf("Failed to reconcile. Errors: %s", errStr)
	return reconcile.Result{}, err
}
