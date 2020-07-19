# NFS Provisioner Operator

- [NFS Provisioner Operator](#nfs-provisioner-operator)
  - [Usage](#usage)
  - [Build](#build)
  - [Testing](#testing)
    - [Test with the NFS Provisioner](#test-with-the-nfs-provisioner)
    - [Cleanup](#cleanup)
  - [External Resources](#external-resources)

The NFS Provisioner Operator creates a **NFS External Provisioner** which creates a `ReadWriteMany` `PersistentVolumeClaim` to be consumed by any Pod/Container in the cluster. The backend block storage, if not specified, is an IBM Cloud VPC Block using the requested storage class.

The goal of this NFS Provisioner Operator is to make it easier to Kubernetes developers to have a PVC that can be used by many pods (`ReadWriteMany`) using the same volume, saving resources and money.

Refer to the [documentation](./docs/index.md) for more information about the design and architecture of the NFS Provisioner Operator.

## Usage

Before use it you need to deploy the NFS Provisioner Operator, this is usually done, but not necesarilly, when the cluster is created. The deployment can be done with the following `kubectl` command:

```bash
kubectl create -f https://www.johandry.com/nfs-operator/nfs_provisioner.yaml
```

The first step after the NFS Provisioner Operator is deployed is to create a NFS CustomResource defining the `storageClassName` and define the backing block storage. The backend block storage will be created by the operartor or you can provide an existing storage accessible through a PVC.

An example of a regular NFS CustomResource could be like this.

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

Notice the value of `storageClass` and the values of the `backingStorage` specification. The backend block storage will be of `storageClass` name `ibmc-vpc-block-general-purpose` with **10Gb**.

If you have your own block storage to be used by the NFS Provisioner, read the [documentation](./docs/index.md).

To use the storage, create a PVC using the given storage class, in this example it is `example-nfs`. The VPC for this example would be like this.

```yaml
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: nfs
spec:
  storageClassName: example-nfs
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1Mi
```

This VPC request 1Mb and the name is `nfs`, as its access mode is `ReadWriteMany` many containers or Pods can use it.

The following Pod example uses the NFS Provider creating a volume from the PVC using the claim name `nfs`. Then mount the volume in any directory of the container with `volumeMounts`.

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

More information can be found in the [documentation](./docs/index.md).

## Build

The development of the operator focus on basically the files:

- `pkg/apis/ibmcloud/v1alpha1/nfs_types.go`: defines the operator specs and status, modifying this file requires to execute `make generate`
- `pkg/controller/nfs/nfs_controller.go`: containg the `Reconcile` function to create or delete all the required resources.
- `pkg/controller/storage`: packages with all the logic to create NFS Provisioner and the Backing Storage (PVC)

After modify any of the files it's recommended to execute `make` to generate the CR and CRD's, and to build the Docker container with the NFS Operator and finally push it to the Docker Registry.

To quick test the operator (build, deploy and test locally), execute:

```bash
make all
```

To test using the consumer application on the cluster, execute:

```bash
make build-operator deploy test

make delete
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
