package storage

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// newServiceAccount returns the definition of this resource as should exists
func (p *NfsProvisioner) newServiceAccount() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: appLabelnName,
		},
	}
}

// applyServiceAccount creates this resource if does not exists
// nil, nil => exists
// nil, err => fail to retreive
// ok,  nil => created
// ok,  err => fail to create
func (p *NfsProvisioner) applyServiceAccount() (string, metav1.Object, error) {
	serviceAccount := p.newServiceAccount()
	name := serviceAccount.Name
	fullName := serviceAccount.GetObjectKind().GroupVersionKind().Kind + "/" + name

	found := &corev1.ServiceAccount{}
	err := p.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: p.Namespace}, found)
	if err == nil { // exists
		return fullName, nil, nil
	}

	if errors.IsNotFound(err) { // does not exists, not found
		if err := p.client.Create(context.TODO(), serviceAccount); err != nil {
			return fullName, serviceAccount, fmt.Errorf("fail to create the object. %s", err)
		}
		return fullName, serviceAccount, nil
	}

	return fullName, nil, fmt.Errorf("fail to retreive the object. %s", err)
}

/*
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nfs-provisioner
*/

// newService returns the definition of this resource as should exists
func (p *NfsProvisioner) newService() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: appLabelnName,
			Labels: map[string]string{
				"app": appLabelnName,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "nfs",
					Port: int32(2049),
				},
				{
					Name:     "nfs-udp",
					Port:     int32(2049),
					Protocol: corev1.ProtocolUDP,
				},
				{
					Name: "nlockmgr",
					Port: int32(32803),
				},
				{
					Name:     "nlockmgr-udp",
					Port:     int32(32803),
					Protocol: corev1.ProtocolUDP,
				},
				{
					Name: "mountd",
					Port: int32(20048),
				},
				{
					Name:     "mountd-udp",
					Port:     int32(20048),
					Protocol: corev1.ProtocolUDP,
				},
				{
					Name: "rquotad",
					Port: int32(875),
				},
				{
					Name:     "rquotad-udp",
					Port:     int32(875),
					Protocol: corev1.ProtocolUDP,
				},
				{
					Name: "rpcbind",
					Port: int32(111),
				},
				{
					Name:     "rpcbind-udp",
					Port:     int32(111),
					Protocol: corev1.ProtocolUDP,
				},
				{
					Name: "statd",
					Port: int32(662),
				},
				{
					Name:     "statd-udp",
					Port:     int32(662),
					Protocol: corev1.ProtocolUDP,
				},
			},
			Selector: map[string]string{
				"app": appLabelnName,
			},
		},
	}
}

// applyService creates this resource if does not exists
// nil, nil => exists
// nil, err => fail to retreive
// ok,  nil => created
// ok,  err => fail to create
func (p *NfsProvisioner) applyService() (string, metav1.Object, error) {
	service := p.newService()
	name := service.Name
	fullName := service.GetObjectKind().GroupVersionKind().Kind + "/" + name

	found := &corev1.Service{}
	err := p.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: p.Namespace}, found)
	if err == nil { // exists
		return fullName, nil, nil
	}

	if errors.IsNotFound(err) { // does not exists, not found
		if err := p.client.Create(context.TODO(), service); err != nil {
			return fullName, service, fmt.Errorf("fail to create the object. %s", err)
		}
		return fullName, service, nil
	}

	return fullName, nil, fmt.Errorf("fail to retreive the object. %s", err)
}

/*
kind: Service
apiVersion: v1
metadata:
  name: nfs-provisioner
  labels:
    app: nfs-provisioner
spec:
  ports:
    - name: nfs
      port: 2049
    - name: nfs-udp
      port: 2049
      protocol: UDP
    - name: nlockmgr
      port: 32803
    - name: nlockmgr-udp
      port: 32803
      protocol: UDP
    - name: mountd
      port: 20048
    - name: mountd-udp
      port: 20048
      protocol: UDP
    - name: rquotad
      port: 875
    - name: rquotad-udp
      port: 875
      protocol: UDP
    - name: rpcbind
      port: 111
    - name: rpcbind-udp
      port: 111
      protocol: UDP
    - name: statd
      port: 662
    - name: statd-udp
      port: 662
      protocol: UDP
  selector:
    app: nfs-provisioner
*/

func (p *NfsProvisioner) newDeployment() *appsv1.Deployment {
	replicas := int32(1)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: appLabelnName,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": appLabelnName,
				},
			},
			Replicas: &replicas,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": appLabelnName,
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: appLabelnName,
					Containers: []corev1.Container{
						{
							Name:  appLabelnName,
							Image: imageName,
							Ports: []corev1.ContainerPort{
								{
									Name:          "nfs",
									ContainerPort: 2049,
								},
								{
									Name:          "nfs-udp",
									ContainerPort: 2049,
									Protocol:      corev1.ProtocolUDP,
								},
								{
									Name:          "nlockmgr",
									ContainerPort: 32803,
								},
								{
									Name:          "nlockmgr-udp",
									ContainerPort: 32803,
									Protocol:      corev1.ProtocolUDP,
								},
								{
									Name:          "mountd",
									ContainerPort: 20048,
								},
								{
									Name:          "mountd-udp",
									ContainerPort: 20048,
									Protocol:      corev1.ProtocolUDP,
								},
								{
									Name:          "rquotad",
									ContainerPort: 875,
								},
								{
									Name:          "rquotad-udp",
									ContainerPort: 875,
									Protocol:      corev1.ProtocolUDP,
								},
								{
									Name:          "rpcbind",
									ContainerPort: 111,
								},
								{
									Name:          "rpcbind-udp",
									ContainerPort: 111,
									Protocol:      corev1.ProtocolUDP,
								},
								{
									Name:          "statd",
									ContainerPort: 662,
								},
								{
									Name:          "statd-udp",
									ContainerPort: 662,
									Protocol:      corev1.ProtocolUDP,
								},
							},
							SecurityContext: &corev1.SecurityContext{
								Capabilities: &corev1.Capabilities{
									Add: []corev1.Capability{
										"DAC_READ_SEARCH",
										"SYS_RESOURCE",
									},
								},
							},
							Args: []string{
								"-provisioner=ibmcloud/nfs",
							},
							Env: []corev1.EnvVar{
								{
									Name: "POD_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
								{
									Name:  "SERVICE_NAME",
									Value: appLabelnName,
								},
								{
									Name: "POD_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
							ImagePullPolicy: corev1.PullIfNotPresent,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "export-volume",
									MountPath: "/export",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "export-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									// TODO: Change the ClaimName for the user provided PVC
									ClaimName: "nfs-block-custom",
								},
							},
						},
					},
				},
			},
		},
	}
}

// applyDeployment creates this resource if does not exists
func (p *NfsProvisioner) applyDeployment() (string, metav1.Object, error) {
	deployment := p.newDeployment()
	name := deployment.Name
	fullName := deployment.GetObjectKind().GroupVersionKind().Kind + "/" + name

	found := &appsv1.Deployment{}
	err := p.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: p.Namespace}, found)
	if err == nil { // exists
		return fullName, nil, nil
	}

	if errors.IsNotFound(err) { // does not exists, not found
		if err := p.client.Create(context.TODO(), deployment); err != nil {
			return fullName, deployment, fmt.Errorf("fail to create the object. %s", err)
		}
		return fullName, deployment, nil
	}

	return fullName, nil, fmt.Errorf("fail to retreive the object. %s", err)
}

/*
kind: Deployment
apiVersion: apps/v1
metadata:
  name: nfs-provisioner
spec:
  selector:
    matchLabels:
      app: nfs-provisioner
  replicas: 1
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: nfs-provisioner
    spec:
      serviceAccount: nfs-provisioner
      containers:
        - name: nfs-provisioner
          image: quay.io/kubernetes_incubator/nfs-provisioner:latest
          ports:
            - name: nfs
              containerPort: 2049
            - name: nfs-udp
              containerPort: 2049
              protocol: UDP
            - name: nlockmgr
              containerPort: 32803
            - name: nlockmgr-udp
              containerPort: 32803
              protocol: UDP
            - name: mountd
              containerPort: 20048
            - name: mountd-udp
              containerPort: 20048
              protocol: UDP
            - name: rquotad
              containerPort: 875
            - name: rquotad-udp
              containerPort: 875
              protocol: UDP
            - name: rpcbind
              containerPort: 111
            - name: rpcbind-udp
              containerPort: 111
              protocol: UDP
            - name: statd
              containerPort: 662
            - name: statd-udp
              containerPort: 662
              protocol: UDP
          securityContext:
            capabilities:
              add:
                - DAC_READ_SEARCH
                - SYS_RESOURCE
          args:
            - "-provisioner=ibmcloud/nfs"
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: SERVICE_NAME
              value: nfs-provisioner
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
          imagePullPolicy: "IfNotPresent"
          volumeMounts:
            - name: export-volume
              mountPath: /export
      volumes:
        - name: export-volume
          persistentVolumeClaim:
            claimName: nfs-block-custom
*/
