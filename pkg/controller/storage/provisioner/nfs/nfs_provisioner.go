package storage

import (
	ibmcloudv1alpha1 "github.com/johandry/nfs-operator/pkg/apis/ibmcloud/v1alpha1"
	"github.com/johandry/nfs-operator/pkg/controller/storage"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	appName                   = "nfs-provisioner"
	imageName                 = "quay.io/kubernetes_incubator/nfs-provisioner:latest"
	storageClassName          = "ibmcloud-nfs"
	provisionerName           = "ibmcloud/nfs"
	persistentVolumeClaimName = "nfs"
)

// NfsProvisioner provision a NFS server on a Pod
type NfsProvisioner struct {
	Namespace string
	spec      ibmcloudv1alpha1.NfsSpec
	status    ibmcloudv1alpha1.NfsStatus
	client    client.Client
	applyFn   []storage.ApplyFn
}

// New creates a new NFS Provisioner to apply the required resources if doesn't exists
func New(client client.Client, namespace string, spec ibmcloudv1alpha1.NfsSpec, status ibmcloudv1alpha1.NfsStatus) *NfsProvisioner {
	nfsProvisioner := &NfsProvisioner{
		client:    client,
		Namespace: namespace,
		spec:      spec,
		status:    status,
	}

	nfsProvisioner.applyFn = []storage.ApplyFn{
		// deployment.go
		nfsProvisioner.applyServiceAccount,
		nfsProvisioner.applyService,
		nfsProvisioner.applyDeployment,
		// rbac.go
		nfsProvisioner.applyClusterRole,
		nfsProvisioner.applyClusterRoleBinding,
		nfsProvisioner.applyRole,
		nfsProvisioner.applyRoleBinding,
		// class.go
		nfsProvisioner.applyStorageClass,
		// claim.go
		nfsProvisioner.applyPersistentVolumeClaim,
	}

	return nfsProvisioner
}

// Apply create all the NFS Provisioner resources if they do not exists
func (p *NfsProvisioner) Apply() storage.Resources {
	resources := storage.Resources{}

	for _, fn := range p.applyFn {
		name, serviceAccount, err := fn()
		resources[name] = &storage.Resource{
			Obj: serviceAccount,
			Err: err,
		}
	}

	return resources
}

// func (p *NfsProvisioner) ObjectFromContent(content []byte) runtime.Object {
// 	r := bytes.NewBuffer(content)
// 	result := p.client.NewBuilder().
// 		Unstructured().
// 		Schema(validation.Schema).
// 		ContinueOnError().
// 		NamespaceParam(p.Namespace).DefaultNamespace().
// 		Stream(r, "").
// 		Flatten().
// 		Do()

// 	return result.Object
// }
