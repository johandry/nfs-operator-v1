package vpcblock

import (
	"github.com/go-logr/logr"
	ibmcloudv1alpha1 "github.com/johandry/nfs-operator/pkg/apis/ibmcloud/v1alpha1"
	"github.com/johandry/nfs-operator/pkg/resources"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	persistentVolumeClaimName = "nfs-block-custom"
)

// Resources implements the resources.Group interface
type Resources struct {
	resources []resources.Reconcilable
}

// New creates a resources group for the NFS Provisioner
func New(owner *ibmcloudv1alpha1.Nfs, client client.Client, scheme *runtime.Scheme, log logr.Logger) *Resources {
	log = log.WithName("vpc-block")
	resources := []resources.Reconcilable{
		PersistentVolumeClaim(owner, client, scheme, log),
	}

	return &Resources{
		resources: resources,
	}
}

// Resources returns the group of reconcilable resources required to
// have a NFS Provisioner
func (r *Resources) Resources() []resources.Reconcilable {
	return r.resources
}

// Reconcile creates the the Resources that does not exists and sets the Owner as an
// owner reference on the Object
func (r *Resources) Reconcile() (reconcile.Result, error) {
	for _, resource := range r.resources {
		result, err := resource.Reconcile()
		if err != nil {
			return result, err
		}
	}
	return reconcile.Result{}, nil
}
