package nfs

// import (
// 	"context"
// 	"fmt"

// 	corev1 "k8s.io/api/core/v1"
// 	"k8s.io/apimachinery/pkg/api/errors"
// 	"k8s.io/apimachinery/pkg/api/resource"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/types"
// )

// /*
// kind: PersistentVolumeClaim
// apiVersion: v1
// metadata:
//   name: nfs
// spec:
// 	storageClassName: ibmcloud-nfs
//   accessModes:
//     - ReadWriteMany
//   resources:
//     requests:
//       storage: 1Mi
// */

// // newPersistentVolumeClaim returns the definition of this resource as should exists
// func (p *Provisioner) newPersistentVolumeClaim() *corev1.PersistentVolumeClaim {
// 	storageClassNameStr := storageClassName
// 	return &corev1.PersistentVolumeClaim{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      persistentVolumeClaimName,
// 			Namespace: r.Owner.Namespace,
// 		},
// 		Spec: corev1.PersistentVolumeClaimSpec{
// 			StorageClassName: &storageClassNameStr,
// 			AccessModes: []corev1.PersistentVolumeAccessMode{
// 				corev1.ReadWriteMany,
// 			},
// 			Resources: corev1.ResourceRequirements{
// 				Requests: corev1.ResourceList{
// 					corev1.ResourceStorage: resource.MustParse("1Mi"),
// 				},
// 			},
// 		},
// 	}
// }

// // applyPersistentVolumeClaim creates this resource if does not exists
// // nil, nil => exists
// // nil, err => fail to retreive
// // ok,  nil => created
// // ok,  err => fail to create
// func (p *Provisioner) applyPersistentVolumeClaim() (string, metav1.Object, error) {
// 	persistentVolumeClaim := p.newPersistentVolumeClaim()
// 	name := persistentVolumeClaim.Name
// 	fullName := persistentVolumeClaim.APIVersion + "/" + name

// 	found := &corev1.PersistentVolumeClaim{}
// 	err := p.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: r.Owner.Namespace}, found)
// 	if err == nil { // exists
// 		return fullName, nil, nil
// 	}

// 	if errors.IsNotFound(err) { // does not exists, not found
// 		if err := p.client.Create(context.TODO(), persistentVolumeClaim); err != nil {
// 			return fullName, persistentVolumeClaim, fmt.Errorf("fail to create the object. %s", err)
// 		}
// 		return fullName, persistentVolumeClaim, nil
// 	}

// 	return fullName, nil, fmt.Errorf("fail to retreive the object. %s", err)
// }
