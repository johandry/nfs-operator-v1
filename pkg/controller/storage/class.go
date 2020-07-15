package storage

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	storagev1 "k8s.io/api/storage/v1"
)

// newStorageClass returns the definition of this resource as should exists
func (p *NfsProvisioner) newStorageClass() *storagev1.StorageClass {
	return &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: storageClassName,
		},
		Provisioner: provisionerName,
		MountOptions: []string{
			"vers=4.1",
		},
	}
}

// applyStorageClass creates this resource if does not exists
// nil, nil => exists
// nil, err => fail to retreive
// ok,  nil => created
// ok,  err => fail to create
func (p *NfsProvisioner) applyStorageClass() (string, metav1.Object, error) {
	storageClass := p.newStorageClass()
	name := storageClass.Name

	found := &storagev1.StorageClass{}
	err := p.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: p.Namespace}, found)
	if err == nil { // exists
		return name, nil, nil
	}

	if errors.IsNotFound(err) { // does not exists, not found
		if err := p.client.Create(context.TODO(), storageClass); err != nil {
			return name, storageClass, fmt.Errorf("fail to create the object. %s", err)
		}
		return name, storageClass, nil
	}

	return name, nil, fmt.Errorf("fail to retreive the object. %s", err)
}

/*
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: ibmcloud-nfs
provisioner: ibmcloud/nfs
mountOptions:
  - vers=4.1
*/
