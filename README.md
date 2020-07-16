# NFS Operator

- [NFS Operator](#nfs-operator)
  - [Deploy the NFS Operator](#deploy-the-nfs-operator)
  - [Usage](#usage)
  - [Build](#build)
  - [Testing](#testing)
    - [Test with the NFS Provisioner](#test-with-the-nfs-provisioner)
    - [Cleanup](#cleanup)
  - [External Resources](#external-resources)

The PersistenVolumeClaim available on IBM Cloud Gen 2 - at this time - only allows access mode `ReadWriteOnce`. Also, there is a limit of block storages you can have in a VPC, as well as there are limitations about the volume size. This operator allows you to have a volume available to many Pods using the same block storage, reducing cost and improving the management of resources. The operator creates a Pod mounting the created or requested PVC and sharing that storage to the cluster using NFS.

## Deploy the NFS Operator

_To be completed ..._

Optionally you can create a volume to be used as backup storage. The NFS Provisioner will share the storage from this volume. If it's not created the NFS Provisioner Operator will create it for you with the provided specifications.

To create the volume use the file `kubernetes/pvc.yaml` with the definition of a Persisten Volume Claim with the profile `ibmc-vpc-block-5iops-tier`. As it's a dynamic provisioner, the Persisten Volume will be created. The size cannot be less than 10Gb, it's the minimun.

```yaml
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
```

Modify the Custom Resource file, i.e. `deploy/crds/ibmcloud.ibm.com_v1alpha1_nfs_cr.yaml`, to provide the name of the created PVC, in this case it's `nfs-block-custom`. This is an example of how this file would be:

```yaml
apiVersion: ibmcloud.ibm.com/v1alpha1
kind: Nfs
metadata:
  name: nfs
spec:
  storageClass: ibmcloud-nfs
  provisionerAPI: some-provisioner-api
  backingStorage:
    pvcName: nfs-block-custom
```

If the PVC is not created, then modify the Custom Resource providing all the volume specifications required. This is an example:

```yaml
apiVersion: ibmcloud.ibm.com/v1alpha1
kind: Nfs
metadata:
  name: nfs
spec:
  storageClass: example-nfs
  provisionerAPI: example.com/nfs
  backingStorage:
    pvcName: nfs-block-custom
    storageClass: ibmc-vpc-block-general-purpose
    storageSize: 10Gi
```

_To be completed ..._

## Usage

_To be completed ..._

To use the NFS Provider from a container create a volume from a PVC using the same name as defined in the parameter `name` of the NFS Custom Resource. Then, as a regular consumption of the volume, mount it in any directory in the container and use it.

This is a simple example of a container that is creating a file in the NFS volume:

```yaml
kind: Pod
apiVersion: v1
metadata:
  name: consumer
spec:
  containers:
    - name: consumer
      image: busybox
      command:
        - "/bin/sh"
      args:
        - "-c"
        - "touch /mnt/SUCCESS && exit 0 || exit 1"
      volumeMounts:
        - name: nfs-pvc
          mountPath: "/mnt"
  restartPolicy: "Never"
  volumes:
    - name: nfs-pvc
      persistentVolumeClaim:
        claimName: nfs
```

A demo application to use the NFS service can be found in the `kubernetes/consumer` folder and it's a simple API for movies. The database - a single JSON file - is stored in the shared volumen. The deployment uses a initContainer to move the JSON database/file to the shared volume.

_To be completed ..._

## Build

The development of the operator focus on basically the files:

- `pkg/apis/ibmcloud/v1alpha1/nfs_types.go`: defines the operator specs and status, modifying this file requires to execute `make generate`
- `pkg/controller/nfs/nfs_controller.go`: containg the `Reconcile` function to create or delete all the required resources.
- `pkg/controller/storage`: package with all the logic to create NFS Provisioner and the Backing Storage (PVC)

After modify any of the files it's recommended to execute `make` to generate the CR and CRD's, and to build the Docker container with the NFS Operator and finally push it to the Docker Registry.

To quick test the operator (build, deploy and test locally), execute:

```bash
make all
```

Refer to the [DEVELOPMENT](./DEVELOPMENT.md) document for more information. Optionally, read the `Makefile` to be familiar with all the tasks you can execute for testing.

## Testing

The tests require a Kubernetes cluster on IBM Cloud, to get one follow the instructions from the testing [README](./test/README.md) document or follow the following quick start instructions.

```bash
make environment

# Optionally, create the PVC
make deploy-pvc

# Optionally, edit the Custom Resource or, at least, confirm the specifications
vim deploy/crds/*_cr.yaml

make deploy
make deploy-consumer

# Test the Operator locally
make test-local
# Or test it with the consumer application
make test
```

### Test with the NFS Provisioner

If the Operator is not working, verify the NFS Provider (the one the operator deploy) works correctlly. You can do this test with the following instructions:

```bash
cd test
make test-provisioner-deployment
make test-consumer
```

### Cleanup

To know the status of the resources created by the NFS Operator or the NFS Provisioner, execute `make list` and to know all the resources in the cluster, either created by the code or external, execute `make list-all`.

To cleanup the cluster, deleting all the created resources, execute:

```bash
make delete
```

To wipe out everything, including the IKS cluster execute `make clean` or to get the cluster to the original state (not recommended) execute `make purge`.

Refer to the `test/` folder and [README](./test/README.md) document for more information. Optionally, read the `Makefile` to be familiar with all the tasks you can execute for testing.

## External Resources

- [NFS Provisioner](https://github.com/kubernetes-incubator/external-storage/tree/master/nfs)
- CSI implementation for EFS and NFS
- [CSI Driver NFS](https://github.com/kubernetes-csi/csi-driver-nfs)
- Rook NFS operator
- [Rook Operator Kit](https://github.com/rook/operator-kit)
- [External Storage Provisioners](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner)
- [Operator SDK](https://github.com/operator-framework/operator-sdk)
- [Operator SDK Docs for Go](https://sdk.operatorframework.io/docs/golang/quickstart/)
