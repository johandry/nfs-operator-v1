SHELL		= /bin/bash

OPERATOR_NAME	?= nfs-operator
KIND 					?= nfs
APINAME 			?= ibmcloud.ibm.com/v1alpha1
VERSION				?= 1.0.2
REGISTRY 			?= johandry
IMAGE 				 = $(REGISTRY)/$(OPERATOR_NAME):$(VERSION)
MUTABLE_IMAGE  = $(REGISTRY)/$(OPERATOR_NAME):latest

ECHO 			= echo -e
C_STD 		= $(shell echo -e "\033[0m")
C_RED			= $(shell echo -e "\033[91m")
C_GREEN 	= $(shell echo -e "\033[92m")
P 			 	= $(shell echo -e "\033[92m> \033[0m")
OK 				= $(shell echo -e "\033[92m[ OK  ]\033[0m")
ERROR	 		= $(shell echo -e "\033[91m[ERROR]\033[0m")
PASS	 		= $(shell echo -e "\033[92m[PASS ]\033[0m")
FAIL	 		= $(shell echo -e "\033[91m[FAIL ]\033[0m")

default: build-operator

all: build-operator release-operator deploy test-local

## Build

generate:
	operator-sdk generate crds
	operator-sdk generate k8s

build-image:
	operator-sdk build $(MUTABLE_IMAGE)
	docker tag  $(MUTABLE_IMAGE) $(IMAGE)

build-operator: generate build-image

## Deploy

deploy-operator: release-to-test
	kubectl apply -f deploy/service_account.yaml
	kubectl apply -f deploy/role.yaml
	kubectl apply -f deploy/role_binding.yaml
	kubectl apply -f deploy/operator.yaml

deploy-crd:
	kubectl apply -f deploy/crds/*_crd.yaml
	kubectl apply -f deploy/crds/*_cr.yaml

deploy-pvc:
	$(MAKE) -C test/kubernetes deploy-pvc

deploy-consumer:
	$(MAKE) -C test/kubernetes deploy-consumer

deploy: deploy-operator deploy-crd

## Release

reset-operator-yaml:
	if [[ -e deploy/operator.yaml.org ]]; then mv deploy/operator.yaml.org deploy/operator.yaml 2>/dev/null; fi

init-operator-yaml: reset-operator-yaml
	sed -i.org 's|REPLACE_IMAGE|$(IMAGE)|g' deploy/operator.yaml

release-to-test: init-operator-yaml
	cp -R deploy/*.yaml test/kubernetes/nfs-operator/
	$(RM) -f test/kubernetes/nfs-operator/operator.yaml.org
	@$(MAKE) reset-operator-yaml

release-operator:
	docker push $(IMAGE)
	docker push $(MUTABLE_IMAGE)

release: build-operator release-operator init-operator-yaml
	cat deploy/service_account.yaml  > docs/nfs_provisioner.yaml
	@echo "---"			 								>> docs/nfs_provisioner.yaml
	cat deploy/role.yaml 						>> docs/nfs_provisioner.yaml
	@echo "---"			 								>> docs/nfs_provisioner.yaml
	cat deploy/role_binding.yaml		>> docs/nfs_provisioner.yaml
	@echo "---"			 								>> docs/nfs_provisioner.yaml
	cat deploy/operator.yaml 				>> docs/nfs_provisioner.yaml
	@echo "---"			 								>> docs/nfs_provisioner.yaml
	cat deploy/crds/*_crd.yaml 			>> docs/nfs_provisioner.yaml
	@$(MAKE) reset-operator-yaml

## Test

environment:
	$(MAKE) -C test environment

test-local:
	OPERATOR_NAME=$(OPERATOR_NAME) operator-sdk run local --watch-namespace=default

test:
	$(MAKE) -C test

## List Resources

# list the following resources: consumer, pvc, nfs, nfs-operator, nfs-provisioner, all
list-%:
	@$(MAKE) -C test/kubernetes $@

list: list-pvc list-nfs list-consumer

## Cleanup

delete-operator:
	kubectl delete -f deploy/operator.yaml
	kubectl delete -f deploy/role_binding.yaml
	kubectl delete -f deploy/role.yaml
	kubectl delete -f deploy/service_account.yaml

delete-crd:
	kubectl delete -f deploy/crds/*_cr.yaml
	kubectl delete -f deploy/crds/*_crd.yaml

delete-pvc:
	$(MAKE) -C test/kubernetes delete-pvc

delete-consumer:
	$(MAKE) -C test/kubernetes delete-consumer

delete-nfs-provisioner:
	$(MAKE) -C test/kubernetes delete-nfs-provisioner

delete:  delete-crd delete-operator

delete-all: delete delete-consumer delete-nfs-provisioner delete-pvc

destroy:
	$(MAKE) -C test/terraform clean

clean: delete destroy

purge: reset-operator-yaml clean
	$(MAKE) -C test remove

## Bootstrap

# Don't use the following rules, they are here mainly for documentation

add-api:
	K=$(KIND); CKIND=$$(tr '[:lower:]' '[:upper:]' <<< $${K:0:1})$${K:1}; \
	operator-sdk add api --kind $${CKIND} --api-version $(APINAME)

add-controller:
	K=$(KIND); CKIND=$$(tr '[:lower:]' '[:upper:]' <<< $${K:0:1})$${K:1}; \
	operator-sdk add controller --kind $${CKIND} --api-version $(APINAME)

add: add-api add-controller

install:
	brew install operator-sdk
