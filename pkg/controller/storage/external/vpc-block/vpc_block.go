package vpcblock

import (
	ibmcloudv1alpha1 "github.com/johandry/nfs-operator/pkg/apis/ibmcloud/v1alpha1"
	"github.com/johandry/nfs-operator/pkg/controller/storage"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	persistentVolumeClaimName = "nfs-block-custom"
)

// VpcBlock is used to create and manage an IBM Cloud VPC Block storage used as
// backing external storage for the provisioner
type VpcBlock struct {
	Namespace        string
	name             string
	storageClassName string
	storageSize      string
	client           client.Client
	applyFn          []storage.ApplyFn
}

// New create a new VPC Block struct with the given storageClassName and storage size
func New(client client.Client, namespace string, spec ibmcloudv1alpha1.BackingStorageSpec) *VpcBlock {
	vpcBlock := &VpcBlock{
		Namespace:        namespace,
		name:             spec.PvcName,
		storageClassName: spec.StorageClass,
		storageSize:      spec.StorageSize,
		client:           client,
	}

	vpcBlock.applyFn = []storage.ApplyFn{}

	return vpcBlock
}

// Apply create all the VPC Block resources if they do not exists
func (b *VpcBlock) Apply() storage.Resources {
	resources := storage.Resources{}

	for _, fn := range b.applyFn {
		name, serviceAccount, err := fn()
		resources[name] = &storage.Resource{
			Obj: serviceAccount,
			Err: err,
		}
	}

	return resources
}
