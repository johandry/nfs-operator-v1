package storage

import (
	"context"
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var contentClusterRole = []byte(`
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
`)

// newClusterRole returns the definition of this resource as should exists
func (p *NfsProvisioner) newClusterRole() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: appName + "-runner",
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{},
				Resources: []string{"persistentvolumes"},
				Verbs:     []string{"get", "list", "watch", "create", "delete"},
			},
			{
				APIGroups: []string{},
				Resources: []string{"persistentvolumeclaims"},
				Verbs:     []string{"get", "list", "watch", "update"},
			},
			{
				APIGroups: []string{"storage.k8s.io"},
				Resources: []string{"storageclasses"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{},
				Resources: []string{"events"},
				Verbs:     []string{"create", "update", "patch"},
			},
			{
				APIGroups: []string{},
				Resources: []string{"services", "endpoints"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups:     []string{"extensions"},
				Resources:     []string{"podsecuritypolicies"},
				ResourceNames: []string{"nfs-provisioner"},
				Verbs:         []string{"use"},
			},
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
	fullName := clusterRole.GetObjectKind().GroupVersionKind().Kind + "/" + name

	found := &rbacv1.ClusterRole{}
	objKey, err := client.ObjectKeyFromObject(clusterRole)
	if err != nil {
		return fullName, nil, fmt.Errorf("fail to retreive the object. %s", err)
	}
	if err = p.client.Get(context.TODO(), objKey, found); err == nil { // exists
		return fullName, nil, nil
	}

	if errors.IsNotFound(err) { // does not exists, not found
		if err := p.client.Create(context.TODO(), clusterRole); err != nil {
			return fullName, clusterRole, fmt.Errorf("fail to create the object. %s", err)
		}
		return fullName, clusterRole, nil
	}

	return fullName, nil, fmt.Errorf("fail to retreive the object. %s", err)
}

var contentClusterRoleBinding = []byte(`
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
`)

// newClusterRoleBinding returns the definition of this resource as should exists
func (p *NfsProvisioner) newClusterRoleBinding() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "run-" + appName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      appName,
				Namespace: p.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     appName + "-runner",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
}

// func (p *NfsProvisioner) getClusterRoleBinding() runtime.Object {
// 	r := p.ObjectForContent(contentClusterRoleBinding)

// }

// applyClusterRoleBinding creates this resource if does not exists
// nil, nil => exists
// nil, err => fail to retreive
// ok,  nil => created
// ok,  err => fail to create
func (p *NfsProvisioner) applyClusterRoleBinding() (string, metav1.Object, error) {
	clusterRoleBinding := p.newClusterRoleBinding()
	name := clusterRoleBinding.Name
	fullName := clusterRoleBinding.GetObjectKind().GroupVersionKind().Kind + "/" + name

	found := &rbacv1.ClusterRoleBinding{}
	objKey, err := client.ObjectKeyFromObject(clusterRoleBinding)
	if err != nil {
		return fullName, nil, fmt.Errorf("fail to retreive the object. %s", err)
	}
	if err = p.client.Get(context.TODO(), objKey, found); err == nil { // exists
		return fullName, nil, nil
	}

	if errors.IsNotFound(err) { // does not exists, not found
		if err := p.client.Create(context.TODO(), clusterRoleBinding); err != nil {
			return fullName, clusterRoleBinding, fmt.Errorf("fail to create the object. %s", err)
		}
		return fullName, clusterRoleBinding, nil
	}

	return fullName, nil, fmt.Errorf("fail to retreive the object. %s", err)
}

var contentRole = []byte(`
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: leader-locking-nfs-provisioner
rules:
  - apiGroups: [""]
    resources: ["endpoints"]
    verbs: ["get", "list", "watch", "create", "update", "patch"]
`)

// newRole returns the definition of this resource as should exists
func (p *NfsProvisioner) newRole() *rbacv1.Role {
	return &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name: "leader-locking-" + appName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{},
				Resources: []string{"endpoints"},
				Verbs:     []string{"get", "list", "watch", "create", "update", "patch"},
			},
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
	fullName := role.GetObjectKind().GroupVersionKind().Kind + "/" + name

	found := &rbacv1.Role{}
	err := p.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: p.Namespace}, found)
	if err == nil { // exists
		return fullName, nil, nil
	}

	if errors.IsNotFound(err) { // does not exists, not found
		if err := p.client.Create(context.TODO(), role); err != nil {
			return fullName, role, fmt.Errorf("fail to create the object. %s", err)
		}
		return fullName, role, nil
	}

	return fullName, nil, fmt.Errorf("fail to retreive the object. %s", err)
}

var contentRoleBinding = []byte(`
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
`)

// newRoleBinding returns the definition of this resource as should exists
func (p *NfsProvisioner) newRoleBinding() *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "leader-locking-" + appName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      appName,
				Namespace: p.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     "leader-locking-" + appName,
			APIGroup: "rbac.authorization.k8s.io",
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
	fullName := roleBinding.GetObjectKind().GroupVersionKind().Kind + "/" + name

	found := &rbacv1.RoleBinding{}
	err := p.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: p.Namespace}, found)
	if err == nil { // exists
		return fullName, nil, nil
	}

	if errors.IsNotFound(err) { // does not exists, not found
		if err := p.client.Create(context.TODO(), roleBinding); err != nil {
			return fullName, roleBinding, fmt.Errorf("fail to create the object. %s", err)
		}
		return fullName, roleBinding, nil
	}

	return fullName, nil, fmt.Errorf("fail to retreive the object. %s", err)
}
