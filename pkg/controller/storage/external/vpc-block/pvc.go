package vpcblock

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var contentPersistentVolumeClaim = []byte(`
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: nfs-block-custom
spec:
  storageClassName: ibmc-vpc-block-general-purpose
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
`)

// newPersistentVolumeClaim returns the definition of this resource as should exists
func (b *VpcBlock) newPersistentVolumeClaim() *corev1.PersistentVolumeClaim {
	storageClassNameStr := b.storageClassName
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      persistentVolumeClaimName,
			Namespace: b.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassNameStr,
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(b.storageSize),
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
func (b *VpcBlock) applyPersistentVolumeClaim() (string, metav1.Object, error) {
	persistentVolumeClaim := b.newPersistentVolumeClaim()
	name := persistentVolumeClaim.Name
	fullName := persistentVolumeClaim.GetObjectKind().GroupVersionKind().Kind + "/" + name

	found := &corev1.PersistentVolumeClaim{}
	err := b.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: b.Namespace}, found)
	if err == nil { // exists
		return fullName, nil, nil
	}

	if errors.IsNotFound(err) { // does not exists, not found
		if err := b.client.Create(context.TODO(), persistentVolumeClaim); err != nil {
			return fullName, persistentVolumeClaim, fmt.Errorf("fail to create the object. %s", err)
		}
		return fullName, persistentVolumeClaim, nil
	}

	return fullName, nil, fmt.Errorf("fail to retreive the object. %s", err)
}
