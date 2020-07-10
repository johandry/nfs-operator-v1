# Development Guide of the Operador using the Operator SDK

- [Development Guide of the Operador using the Operator SDK](#development-guide-of-the-operador-using-the-operator-sdk)
  - [QuickStart](#quickstart)
  - [Install the Operador SDK](#install-the-operador-sdk)
  - [Bootstrap the Operator](#bootstrap-the-operator)
  - [Local Test](#local-test)
  - [Testing on Kubernetes](#testing-on-kubernetes)
  - [Deployment](#deployment)
  - [Deployment with the Operator Lifecycle Manager (OLM)](#deployment-with-the-operator-lifecycle-manager-olm)
  - [Cleanup](#cleanup)
  - [Reference to Advance Topics](#reference-to-advance-topics)

## QuickStart

After [install the Operator SDK](#install-the-operador-sdk) and [Boostrap the Operator](#bootstrap-the-operator) adding the API and the Controller, the development process might be just to modify them.

After modify the API file `pkg/apis/ibmcloud/v1alpha1/nfs_types.go`, or the Controller file `pkg/controller/nfs/nfs_controller.go`, update (optionally) the `VERSION` variable in the `Makefile` and execute `make` t generate the CRDs, then build and push the Docker container.

```bash
make
```

To test, modify the file `deploy/crds/ibmcloud.ibm.com_v1alpha1_nfs_cr.yaml` and apply it executing `make deploy-cr`. Then execute `make test-local`.

```bash
make deploy-cr
make test-local
```

To deploy everything, execute `make deploy` then you can verify the reults with:

```bash
make deploy

kubectl get deployment nfs-operator
kubectl get pods
kubectl get nfs
...
```

To remove everything, execute `make clean`

## Install the Operador SDK

```bash
brew install operator-sdk
```

## Bootstrap the Operator

```bash
export OPERATOR_NAME=nfs-operator
export KIND=nfs
export APINAME=ibmcloud.ibm.com/v1alpha1

# It's the KIND capitalized
export CAP_KIND="$(tr '[:lower:]' '[:upper:]' <<< ${KIND:0:1})${KIND:1}" # Nfs
export SHORT_APINAME=$(echo $APINAME | cut -d. -f1) # ibmcloud

operator-sdk new ${OPERATOR_NAME} --type go --repo github.com/johandry/nfs-operator
cd ${OPERATOR_NAME}
git init

operator-sdk add api --kind $CAP_KIND --api-version $APINAME

vim pkg/apis/$SHORT_APINAME/v1alpha1/${KIND}_types.go
```

Edit in `pkg/apis/$SHORT_APINAME/v1alpha1/$KIND_types.go` the struct `${KIND}Spec` to add the parameters:

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
type BackingStorage struct {
  PvcName      string `json:"pvcName,omitempty"`
  StorageClass string `json:"storageClass,omitempty"`
  StorageSize  string `json:"storageSize,omitempty"`
}
type NfsSpec struct {
  StorageClass   string         `json:"storageClass,omitempty"`
  ProvisionerAPI string         `json:"provisionerAPI,omitempty"`
  BackingStorage BackingStorage `json:"backingStorage,omitempty"`
}
```

Do this everytime a `*_types.go` file is modified:

```bash
operator-sdk generate crds
operator-sdk generate k8s

CRD_NAME="$(echo ${APINAME} | cut -d/ -f1)_${KIND}_crd"
kubectl apply -f deploy/crds/$CRD_NAME.yaml

operator-sdk add controller --kind $CAP_KIND --api-version $APINAME

vim pkg/controller/$KIND/${KIND}_controller.go
```

Edit in `pkg/controller/$KIND/${KIND}_controller.go` the function `Reconcile` which is call every time the CR is created, changed or deleted.

```go
reqLogger.Info(fmt.Sprintf("Instance Specs: %+v", instance.Spec))
```

## Local Test

```bashr
opeator-sdk run local --watch-namespace=default
```

## Testing on Kubernetes

```bash
kubectl apply -f <(echo "
apiVersion: $APINAME
kind: $CAP_KIND
metadata:
  name: example-$KIND
spec:
  storageClass: some-class
  provisionerAPI: some-provisioner-api
  backingStorage:
    pvcName: some-backing-storage-pvc-name
    storageClass: some-backing-storage-class-name
    storageSize: some-backing-storage-size
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
vim deploy/crds/$(echo ${APINAME} | sed 's|/|_|')_${KIND}_cr.yaml

kubectl apply -f deploy/crds/$(echo ${APINAME} | sed 's|/|_|')_${KIND}_cr.yaml

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
