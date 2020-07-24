package nfs

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	ibmcloudv1alpha1 "github.com/johandry/nfs-operator/pkg/apis/ibmcloud/v1alpha1"
	"github.com/johandry/nfs-operator/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ resources.Reconcilable = &ResService{}

// ResService is the resource Service
type ResService struct {
	Object *corev1.Service
	resources.Resource
}

var contentService = []byte(`
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
`)

// Service creates a Service
func Service(owner *ibmcloudv1alpha1.Nfs, client client.Client, scheme *runtime.Scheme, log logr.Logger) *ResService {
	res := &ResService{}
	res.Resource = resources.New(owner, client, scheme, log)
	res.Object = res.newService()
	apiVersion, kind := resources.GVK(res.Object, res.Scheme)
	res.Log = res.Log.WithValues("Resource.Name", res.Object.GetName(), "Resource.Namespace", res.Object.GetNamespace(), "Resource.APIVersion", apiVersion, "Resource.Kind", kind)

	return res
}

// Get returns the Object from the cluster
func (r *ResService) Get() (runtime.Object, error) {
	return r.getService()
}

// Apply creates the Object if it does not exists
func (r *ResService) Apply() error {
	_, err := r.getService()
	exists, err := resources.Exists(err)
	if exists {
		r.Log.Info("Skip reconcile: Resource already exists")
		return nil
	}
	if err != nil {
		r.Log.Error(err, "Failed to reconcile the resource")
		return err
	}

	// if not exists and no error, then create
	r.Log.Info("Created a new resource")
	return r.Client.Create(context.TODO(), r.Object)
}

// Reconcile creates the Object if it does not exists and sets the Owner as an
// owner reference on the Object
func (r *ResService) Reconcile() (reconcile.Result, error) {
	if r.Owner == nil {
		return reconcile.Result{}, fmt.Errorf("the resource %s/%s does not have an owner", r.Object.Namespace, r.Object.Name)
	}
	r.Log.Info("Reconciling " + r.Object.Name + " resource")
	if err := controllerutil.SetControllerReference(r.Owner, r.Object, r.Scheme); err != nil {
		r.Log.Error(err, "Failed to set controller reference to resource")
		return reconcile.Result{}, err
	}
	err := r.Apply()

	return reconcile.Result{}, err
}

// newService returns the definition of this resource as should exists
func (r *ResService) newService() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appName,
			Namespace: r.Owner.Namespace,
			Labels: map[string]string{
				"app": appName,
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
				"app": appName,
			},
		},
	}
}

func (r *ResService) getService() (*corev1.Service, error) {
	found := &corev1.Service{}
	objKey, err := client.ObjectKeyFromObject(r.Object)
	if err != nil {
		return nil, fmt.Errorf("fail to retreive the object key. %s", err)
	}
	// 	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: r.Object.Name, Namespace: r.Owner.Namespace}, found)
	err = r.Client.Get(context.TODO(), objKey, found)
	if err == nil {
		return found, nil
	}
	return nil, err
}
