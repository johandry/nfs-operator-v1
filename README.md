# NFS Operator

- [NFS Operator](#nfs-operator)
  - [How to use](#how-to-use)
  - [Build and Tests](#build-and-tests)
    - [Requirements](#requirements)
    - [Build](#build)
    - [Test](#test)
    - [Cleanup](#cleanup)
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

### Requirements

Before execute the tests you need the following requirements:

1. Have an IBM Cloud account with required privileges
2. [Install IBM Cloud CLI](https://ibm.github.io/cloud-enterprise-examples/iac/setup-environment#install-ibm-cloud-cli)
3. [Install the IBM Cloud CLI Plugins](https://ibm.github.io/cloud-enterprise-examples/iac/setup-environment#ibm-cloud-cli-plugins) `infrastructure-service`, `schematics` and `container-registry`.
4. [Login to IBM Cloud with the CLI](https://ibm.github.io/cloud-enterprise-examples/iac/setup-environment#login-to-ibm-cloud)
5. [Install Terraform](https://ibm.github.io/cloud-enterprise-examples/iac/setup-environment#install-terraform)
6. [Install IBM Cloud Terraform Provider](https://ibm.github.io/cloud-enterprise-examples/iac/setup-environment#configure-access-to-ibm-cloud)
7. [Configure access to IBM Cloud](https://ibm.github.io/cloud-enterprise-examples/iac/setup-environment#configure-access-to-ibm-cloud) for Terraform and the IBM Cloud CLI setting up the `IC_API_KEY` environment variable.
8. Install the following tools:
   1. [jq](https://stedolan.github.io/jq/download/)
   2. [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

If you have an API Key but is not set neither have the JSON file when it was created, you must recreate the key. Delete the old one if won't be in use anymore. Then execute `make api-key`, set the `IC_API_KEY` and, optionally, validate the requirements with `make check`.

```bash
# Delete the old one, if won't be in use anymore
ibmcloud iam api-keys       # Identify your old API Key Name
ibmcloud iam api-key-delete OLD-NAME

# Create a new one and set it as environment variable
cd test
make api-key

export IC_API_KEY=$(grep '"apikey":' terraform_key.json | sed 's/.*: "\(.*\)".*/\1/')
# Or
export IC_API_KEY=$(jq -r .apikey terraform_key.json)

make check
```

However, the Terraform variables validation is made by terraform, before continue set them up either using the environmet variables `TF_VAR_project_name` and `TF_VAR_owner` (i.e. `export TF_VAR_owner=$USER`) or the `terraform/terraform.tfvars` file, like this:

```hcl
project_name = "nfs-op-ja"
owner        = "johandry"
```

You may add more variables to customize the cluster, for example like this one, to have a larger cluster:

```hcl
region         = "us-south"
vpc_zone_names = ["us-south-1", "us-south-2", "us-south-3"]
flavors        = ["cx2.2x4", "cx2.4x8", "cx2.8x16"]
workers_count  = [3, 2, 1]
k8s_version    = "1.18"
```

### Build

To create the IKS cluster and deploy the required resources execute the `apply` rule defining wich resource to deploy with the variable `RESOURCE`, the possible values are: `provisioner` and ~~`operator`~~.

```bash
cd test
make apply RESOURCE=provisioner
```

To make it simpler, you can also execute `make R=provisioner` having the same results.

### Test

To test it's required to complete the build rules in the previous section, then execute the `test` rule.

```bash
cd test
make apply RESOURCE=provisioner

make test
```

Or, you can simplly execute: `make all R=provisioner` to get the build and test done.

### Cleanup

To destroy your environment and cleanup what you have created, execute:

```bash
cd test
make clean
```

To cleanup the Kubernetes cluster of resources without destroying it, execute the `delete` rule:

```bash
make delete
```

## External Resources

- [NFS Provisioner](https://github.com/kubernetes-incubator/external-storage/tree/master/nfs)
- CSI implementation for EFS and NFS
- [CSI Driver NFS](https://github.com/kubernetes-csi/csi-driver-nfs)
- Rook NFS operator
- [Rook Operator Kit](https://github.com/rook/operator-kit)
