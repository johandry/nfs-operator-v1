package storage

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	appName                   = "nfs-provisioner"
	imageName                 = "quay.io/kubernetes_incubator/nfs-provisioner:latest"
	storageClassName          = "ibmcloud-nfs"
	provisionerName           = "ibmcloud/nfs"
	persistentVolumeClaimName = "nfs"
)

// const deploymentReplicas = 1

// NfsProvisionerOpt
// type NfsProvisionerOpt struct {
// 	Replicas *int32
// }

// NfsProvisioner provision a NFS server on a Pod
type NfsProvisioner struct {
	client    client.Client
	Namespace string
	applyFn   []NfsProvisionerApplyFn
	// opt       NfsProvisionerOpt
}

// NfsProvisionerApplyFn is a generic function to apply a resource
type NfsProvisionerApplyFn func() (string, metav1.Object, error)

// NfsProvisionerResource object to create
type NfsProvisionerResource struct {
	Obj metav1.Object
	Err error
}

// NfsProvisionerResources list of object created
type NfsProvisionerResources map[string]*NfsProvisionerResource

// NewNFSProvisioner creates a new NFS Provisioner to apply the required resources if doesn't exists
func NewNFSProvisioner(client client.Client, namespace string /* , opt *NfsProvisionerOpt */) *NfsProvisioner {
	// if opt == nil {
	// 	replicas := deploymentReplicas

	// 	opt := NfsProvisionerOpt{
	// 		Replicas: &replicas,
	// 	}
	// }

	nfsProvisioner := &NfsProvisioner{
		client:    client,
		Namespace: namespace,
		// opt:       opt,
	}

	nfsProvisioner.applyFn = []NfsProvisionerApplyFn{
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
		// claim
		nfsProvisioner.applyPersistentVolumeClaim,
	}

	return nfsProvisioner
}

// Apply create all the NFS Provisioner resources if they do not exists
func (p *NfsProvisioner) Apply() NfsProvisionerResources {
	list := NfsProvisionerResources{}

	for _, fn := range p.applyFn {
		name, serviceAccount, err := fn()
		list[name] = &NfsProvisionerResource{
			Obj: serviceAccount,
			Err: err,
		}
	}

	return list
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
