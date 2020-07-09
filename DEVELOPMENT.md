# Development Guide of the Operador using the Operator SDK

- [Development Guide of the Operador using the Operator SDK](#development-guide-of-the-operador-using-the-operator-sdk)
  - [Install the Operador SDK](#install-the-operador-sdk)
  - [Bootstrap the Operator](#bootstrap-the-operator)
  - [Local Test](#local-test)
  - [Testing on Kubernetes](#testing-on-kubernetes)
  - [Deployment](#deployment)
  - [Deployment with the Operator Lifecycle Manager (OLM)](#deployment-with-the-operator-lifecycle-manager-olm)
  - [Cleanup](#cleanup)
  - [Reference to Advance Topics](#reference-to-advance-topics)

## Install the Operador SDK

```bash
brew install operator-sdk
```

## Bootstrap the Operator

```bash
NAME=nfs-operator
GH_REPO=github.com/johandry/nfs-operator
KIND=nfs
ORG=ibmcloud

# It's the KIND capitalized
CAP_KIND="$(tr '[:lower:]' '[:upper:]' <<< ${KIND:0:1})${KIND:1}"

operator-sdk new $NAME --type go --repo $GH_REPO

cd $NAME
operator-sdk add api --kind $CAP_KIND --api-version $ORG.com/v1alpha1

vim pkg/apis/$ORG/v1alpha1/$KIND_types.go
```

Edit in `pkg/apis/$ORG/v1alpha1/$KIND_types.go` the struct `${KIND}Spec` to add the parameters:

- `StorageClass`: the storage class that the provisioner will listen for requests default: example-nfs
- `ProvisionerApi`: (optional) specify the provisioner api default: example.com/nfs
- `BackingStorage`: structure to define the storage used to provide the NFS service
  - `useExistingPvc`: (_Boolean_) use an existing claim (default `false`)
  - `pvc`: Name of the claim
  - `storageClass`: cloud provider storage class to use for new claim
  - `storageSize`: size of block volume to request from Cloud provider

Optionally edit the struct `${KIND}Status`.

Example:

```go
type NfsSpec struct {
  StorageClass       string `json:"storageClass,omitempty"`
  ProvisionerAPI     string `json:"provisionerAPI,omitempty"`
  BackingStorageName string `json:"backingStorageName,omitempty"`
}
```

Do this everytime a `*_types.go` file is modified:

```bash
operator-sdk generate crds
operator-sdk generate k8s

kubectl apply -f deploy/crds/$ORG.com_$KIND_crd.yaml

operator-sdk add controller --kind $CAP_KIND --api-version $ORG.com/v1alpha1

vim pkg/controller/$KIND/$KIND_controller.go
```

Edit in `pkg/controller/$KIND/$KIND_controller.go` the function `Reconcile` which is call every time the CR is created, changed or deleted.

```go
reqLogger.Info("Instance Specs: %v", instance.Spec)
```

## Local Test

```bash
export OPERATOR_NAME=$NAME

operator-sdk run local --watch-namespace=default
```

## Testing on Kubernetes

```bash
kubectl apply -f <(echo "
apiVersion: $ORG.com/v1alpha1
kind: $CAP_KIND
metadata:
  name: example-$KIND
spec:
  storageClass: some-class
  provisionerAPI: some-provisioner-api
  backingStorageName: some-backing-storage-name
")

kubectl get $KIND

kubectl get pods
```

## Deployment

```bash
operator-sdk build johandry/$KIND-$ORG-operator
docker push johandry/$KIND-$ORG-operator

sed -i.org "s|REPLACE_IMAGE|johandry/$KIND-$ORG-operator|g" deploy/operator.yaml

kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/operator.yaml

kubectl get deployment
```

## Deployment with the Operator Lifecycle Manager (OLM)

Edit `deploy/crds/$ORG.com_v1alpha1_$KIND_cr.yaml` to add/modify the CRD fields/parameters, then apply and validate the changes.

```bash
vim deploy/crds/$ORG.com_v1alpha1_$KIND_cr.yaml

kubectl apply -f deploy/crds/$ORG.com_v1alpha1_$KIND_cr.yaml

kubectl get deployment
kubectl get pods

kubectl get $KIND -o yaml
```

## Cleanup

```bash
kubectl delete -f deploy/crds/$ORG.com_v1alpha1_$KIND_cr.yaml
kubectl delete -f deploy/operator.yaml
kubectl delete -f deploy/role_binding.yaml
kubectl delete -f deploy/role.yaml
kubectl delete -f deploy/service_account.yaml
```

## Reference to Advance Topics

These topics or references may be required for the development of the Operator:

- [Handle Cleanup on Deletion](https://sdk.operatorframework.io/docs/golang/quickstart/#handle-cleanup-on-deletion)
- [Unit Testing](https://sdk.operatorframework.io/docs/golang/unit-testing/)
- [E2E Tests](https://sdk.operatorframework.io/docs/golang/e2e-tests/)
- [Monitoring with Prometheus](https://sdk.operatorframework.io/docs/golang/monitoring/prometheus/)
- [Controller Runtime Client API](https://sdk.operatorframework.io/docs/golang/references/client/)
- [Logging](https://sdk.operatorframework.io/docs/golang/references/logging/)
