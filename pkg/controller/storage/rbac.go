package storage

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	rbacv1 "k8s.io/api/rbac/v1"
)

// newClusterRole returns the definition of this resource as should exists
func (p *NfsProvisioner) newClusterRole() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: appLabelnName + "-runner",
		},
		Rules: []rbacv1.PolicyRule{
			// TODO
		},
	}
}

// applyClusterRole creates this resource if does not exists
// nil, nil => exists
// nil, err => fail to retreive
// ok,  nil => created
// ok,  err => fail to create
func (p *NfsProvisioner) applyClusterRole() (string, metav1.Object, error) {
	clusterRole := p.newClusterRole()
	name := clusterRole.Name

	found := &rbacv1.ClusterRole{}
	err := p.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: p.Namespace}, found)
	if err == nil { // exists
		return name, nil, nil
	}

	if errors.IsNotFound(err) { // does not exists, not found
		if err := p.client.Create(context.TODO(), clusterRole); err != nil {
			return name, clusterRole, fmt.Errorf("fail to create the object. %s", err)
		}
		return name, clusterRole, nil
	}

	return name, nil, fmt.Errorf("fail to retreive the object. %s", err)
}

/*
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: nfs-provisioner-runner
rules:
  - apiGroups: [""]
    resources: ["persistentvolumes"]
    verbs: ["get", "list", "watch", "create", "delete"]
  - apiGroups: [""]
    resources: ["persistentvolumeclaims"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["storageclasses"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["create", "update", "patch"]
  - apiGroups: [""]
    resources: ["services", "endpoints"]
    verbs: ["get"]
  - apiGroups: ["extensions"]
    resources: ["podsecuritypolicies"]
    resourceNames: ["nfs-provisioner"]
    verbs: ["use"]
*/

// newClusterRoleBinding returns the definition of this resource as should exists
func (p *NfsProvisioner) newClusterRoleBinding() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "run-" + appLabelnName,
		},
		Subjects: []rbacv1.Subject{
			// TODO
		},
		RoleRef: rbacv1.RoleRef{
			// TODO
		},
	}
}

// applyClusterRoleBinding creates this resource if does not exists
// nil, nil => exists
// nil, err => fail to retreive
// ok,  nil => created
// ok,  err => fail to create
func (p *NfsProvisioner) applyClusterRoleBinding() (string, metav1.Object, error) {
	clusterRoleBinding := p.newClusterRoleBinding()
	name := clusterRoleBinding.Name

	found := &rbacv1.ClusterRoleBinding{}
	err := p.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: p.Namespace}, found)
	if err == nil { // exists
		return name, nil, nil
	}

	if errors.IsNotFound(err) { // does not exists, not found
		if err := p.client.Create(context.TODO(), clusterRoleBinding); err != nil {
			return name, clusterRoleBinding, fmt.Errorf("fail to create the object. %s", err)
		}
		return name, clusterRoleBinding, nil
	}

	return name, nil, fmt.Errorf("fail to retreive the object. %s", err)
}

/*
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: run-nfs-provisioner
subjects:
  - kind: ServiceAccount
    name: nfs-provisioner
     # replace with namespace where provisioner is deployed
    namespace: default
roleRef:
  kind: ClusterRole
  name: nfs-provisioner-runner
  apiGroup: rbac.authorization.k8s.io
*/

// newRole returns the definition of this resource as should exists
func (p *NfsProvisioner) newRole() *rbacv1.Role {
	return &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name: "leader-locking-" + appLabelnName,
		},
		Rules: []rbacv1.PolicyRule{
			// TODO
		},
	}
}

// applyRole creates this resource if does not exists
// nil, nil => exists
// nil, err => fail to retreive
// ok,  nil => created
// ok,  err => fail to create
func (p *NfsProvisioner) applyRole() (string, metav1.Object, error) {
	role := p.newRole()
	name := role.Name

	found := &rbacv1.Role{}
	err := p.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: p.Namespace}, found)
	if err == nil { // exists
		return name, nil, nil
	}

	if errors.IsNotFound(err) { // does not exists, not found
		if err := p.client.Create(context.TODO(), role); err != nil {
			return name, role, fmt.Errorf("fail to create the object. %s", err)
		}
		return name, role, nil
	}

	return name, nil, fmt.Errorf("fail to retreive the object. %s", err)
}

/*
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: leader-locking-nfs-provisioner
rules:
  - apiGroups: [""]
    resources: ["endpoints"]
    verbs: ["get", "list", "watch", "create", "update", "patch"]
*/

// newRoleBinding returns the definition of this resource as should exists
func (p *NfsProvisioner) newRoleBinding() *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "leader-locking-" + appLabelnName,
		},
		Subjects: []rbacv1.Subject{
			// TODO
		},
		RoleRef: rbacv1.RoleRef{
			// TODO
		},
	}
}

// applyRoleBinding creates this resource if does not exists
// nil, nil => exists
// nil, err => fail to retreive
// ok,  nil => created
// ok,  err => fail to create
func (p *NfsProvisioner) applyRoleBinding() (string, metav1.Object, error) {
	roleBinding := p.newRoleBinding()
	name := roleBinding.Name

	found := &rbacv1.RoleBinding{}
	err := p.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: p.Namespace}, found)
	if err == nil { // exists
		return name, nil, nil
	}

	if errors.IsNotFound(err) { // does not exists, not found
		if err := p.client.Create(context.TODO(), roleBinding); err != nil {
			return name, roleBinding, fmt.Errorf("fail to create the object. %s", err)
		}
		return name, roleBinding, nil
	}

	return name, nil, fmt.Errorf("fail to retreive the object. %s", err)
}

/*
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: leader-locking-nfs-provisioner
subjects:
  - kind: ServiceAccount
    name: nfs-provisioner
    # replace with namespace where provisioner is deployed
    namespace: default
roleRef:
  kind: Role
  name: leader-locking-nfs-provisioner
  apiGroup: rbac.authorization.k8s.io
*/
