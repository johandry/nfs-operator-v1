# NFS Operator

- [NFS Operator](#nfs-operator)
  - [How to use](#how-to-use)
  - [Build and Tests](#build-and-tests)
  - [External Resources](#external-resources)

The PersistenVolumeClaim available on IBM Cloud Gen 2 - at this time - only allows access mode `ReadWriteOnce`. Also, there is a limit of block storages you can have in a VPC, as well as there are limitations about the volume size. This operator allows you to have a volume available to many Pods using the same block storage, reducing cost and improving the management of resources. The operator creates a Pod mounting the created or requested PVC and sharing that storage to the cluster using NFS.

## How to use

Let’s start creating the file `kubernetes/pvc.yaml` with the definition of a Persisten Volume Claim with the profile `ibmc-vpc-block-5iops-tier`. The size cannot be less than 10Gb, it's the minimun.

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

This PVC is the one the NFS Provisioner or Operator uses to provide the NFS service to all the pods in the cluster.

To use the NFS Provider as is, deploy the resources located in the `kubernetes/nfs-provider` folder. The [nfs-provider documentation](https://github.com/kubernetes-incubator/external-storage/blob/master/nfs/README.md#quickstart) explains what and how it does it.

To use the operator _... to be completed ..._

We have an application to use the NFS service, it's in the `kubernetes/consumer` folder and it's a simple API for movies. The database - a single JSON file - is stored in the shared volumen. The deployment uses a initContainer to move the JSON database/file to the shared volume.

## Build and Tests

Refer to the `test/` folder and [README](./test/README.md) document.

## External Resources

- [NFS Provisioner](https://github.com/kubernetes-incubator/external-storage/tree/master/nfs)
- CSI implementation for EFS and NFS
- [CSI Driver NFS](https://github.com/kubernetes-csi/csi-driver-nfs)
- Rook NFS operator
- [Rook Operator Kit](https://github.com/rook/operator-kit)
- [External Storage Provisioners](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner)
- [Operator SDK](https://github.com/operator-framework/operator-sdk)
- [Operator SDK Docs for Go](https://sdk.operatorframework.io/docs/golang/quickstart/)
