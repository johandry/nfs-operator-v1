SHELL		= /bin/bash

OPERATOR_NAME	?= nfs-operator
KIND 					?= nfs
APINAME 			?= ibmcloud.ibm.com/v1alpha1
VERSION				?= 1.0.1
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

all: build-operator deploy test delete

all-local: build-operator deploy test-local

## Build

generate:
	operator-sdk generate crds
	operator-sdk generate k8s

build-image:
	operator-sdk build $(MUTABLE_IMAGE)
	docker tag  $(MUTABLE_IMAGE) $(IMAGE)

push-image:
	docker push $(IMAGE)
	docker push $(MUTABLE_IMAGE)

build-operator: generate build-image push-image

## Deploy

rename-operator-yaml:
	if [[ -e deploy/operator.yaml.org ]]; then mv deploy/operator.yaml.org deploy/operator.yaml 2>/dev/null; fi

init-operator: rename-operator-yaml
	sed -i.org 's|REPLACE_IMAGE|$(IMAGE)|g' deploy/operator.yaml

deploy-operator: init-operator cp-deploy-to-test
	kubectl apply -f deploy/service_account.yaml
	kubectl apply -f deploy/role.yaml
	kubectl apply -f deploy/role_binding.yaml
	kubectl apply -f deploy/operator.yaml

cp-deploy-to-test: init-operator
	cp -R deploy/*.yaml test/kubernetes/nfs-operator/
	$(RM) -f test/kubernetes/nfs-operator/operator.yaml.org

deploy-crds:
	kubectl apply -f deploy/crds/*_crd.yaml
	kubectl apply -f deploy/crds/*_cr.yaml

deploy: deploy-operator deploy-crds

## Test

environment:
	$(MAKE) -C test environment

test-local:
	OPERATOR_NAME=$(OPERATOR_NAME) operator-sdk run local --watch-namespace=default

test:
	$(MAKE) -C test

## List Resources

list-%:
	@$(MAKE) -C test/kubernetes $@

list: list-pvc list-nfs list-consumer

## Cleanup

delete-consumer:
	$(MAKE) -C test/kubernetes delete-consumer

delete-nfs-operator:
	$(MAKE) -C test/kubernetes delete-nfs-operator

delete-nfs-provisioner:
	$(MAKE) -C test/kubernetes delete-nfs-provisioner

delete-pvc:
	$(MAKE) -C test/kubernetes delete-pvc

delete: delete-consumer delete-nfs-operator delete-pvc

destroy:
	$(MAKE) -C test/terraform clean

clean: delete destroy

purge: rename-operator-yaml clean
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
