SHELL				= /bin/bash
ECHO 				= echo -e
C_STD 			= $(shell echo -e "\033[0m")
C_RED				= $(shell echo -e "\033[91m")
C_GREEN 		= $(shell echo -e "\033[92m")
P 			 		= $(shell echo -e "\033[92m> \033[0m")
OK 			 		= $(shell echo -e "\033[92m[ OK  ]\033[0m")
ERROR		 		= $(shell echo -e "\033[91m[ERROR] \033[0m")
PASS		 		= $(shell echo -e "\033[92m[PASS ]\033[0m")
FAIL		 		= $(shell echo -e "\033[91m[FAIL ] \033[0m")

# RESOURCE sets the object to build, deploy or test. Available options are: provisioner
# Soon will be operator
R 			 			 ?= provisioner
RESOURCE 			 ?= $(R)

default: test

all: init test

## Validations:

check-kubectl:
	@if kubectl version | grep -q 'Client Version'; then $(ECHO) "$(P) $(OK) kubectl"; else $(ECHO) "$(P) $(ERROR) kubectl"; exit 1; fi

check-kubernetes:
	$(MAKE) -f ./terraform test

check: check-kubectl check-kubernetes

## Init

environment:
	$(MAKE) -f ./terraform all

deploy-%:
	$(MAKE) -C kubernetes $*

deploy: check deploy-pvc deploy-nfs deploy-consumer

init: environment deploy

## Test

# Test the NFS Provisioner not using the Operator. Use this test to verify the NFS Provisioner if
# the Operator is failing.
test-provisioner-deployment: deploy-pvc deploy-nfs-provisioner deploy-consumer

# Test the NFS Provisioner Operator. The Operator uses the PVC created.
test-operator-pvc-deployment: deploy-pvc deploy-nfs deploy-consumer

# Test the NFS Provisioner Operator. The Operator create a PVC.
test-operator-deployment: deploy-nfs deploy-consumer

# Test the Consumer once the NFS Provisioner exists with the Operator or not, with an existing PVC
# or created
test-consumer:
	$(MAKE) -C kubernetes test

# Default test: deploy the NFS Provisioner Operator, it creates the PVC and the Consumer should work
test: deploy test-consumer

## Cleanup

destroy:
	$(MAKE) -C terraform clean

delete:
	$(MAKE) -C kubernetes clean
	$(MAKE) -C kubernetes clean-pvc

remove:
	$(RM) -rf terraform/.terraform
	$(RM) -rf terraform/.terraform.tfstate*
	$(RM) -rf terraform/terraform.tfstate*
	$(RM) -rf terraform/terraform.tfvars
	$(RM) -rf .terraform_key.json
	$(RM) -rf .target_account

clean: delete destroy

purge: clean remove
