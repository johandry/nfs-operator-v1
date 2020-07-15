package storage

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
)

// newPersistentVolumeClaim returns the definition of this resource as should exists
func (p *NfsProvisioner) newPersistentVolumeClaim() *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: persistentVolumeClaimName,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.NewQuantity(1*2^20, resource.BinarySI),
				},
			},
		},
	}
}

// applyPersistentVolumeClaim creates this resource if does not exists
// nil, nil => exists
// nil, err => fail to retreive
// ok,  nil => created
// ok,  err => fail to create
func (p *NfsProvisioner) applyPersistentVolumeClaim() (string, metav1.Object, error) {
	persistentVolumeClaim := p.newPersistentVolumeClaim()
	name := persistentVolumeClaim.Name

	found := &corev1.PersistentVolumeClaim{}
	err := p.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: p.Namespace}, found)
	if err == nil { // exists
		return name, nil, nil
	}

	if errors.IsNotFound(err) { // does not exists, not found
		if err := p.client.Create(context.TODO(), persistentVolumeClaim); err != nil {
			return name, persistentVolumeClaim, fmt.Errorf("fail to create the object. %s", err)
		}
		return name, persistentVolumeClaim, nil
	}

	return name, nil, fmt.Errorf("fail to retreive the object. %s", err)
}

/*
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: nfs
  annotations:
    volume.beta.kubernetes.io/storage-class: "ibmcloud-nfs"
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1Mi
*/
