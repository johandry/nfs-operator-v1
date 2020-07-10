SHELL		= /bin/bash

OPERATOR_NAME	?= nfs-operator
KIND 					?= nfs
APINAME 			?= ibmcloud.ibm.com/v1alpha1
VERSION				?= 1.0

# It's the value of KIND capitalized
CAP_KIND 	= $(shell K=$(KIND); echo $$(tr '[:lower:]' '[:upper:]' <<< $${K:0:1})$${K:1})
IMAGE 		= $(OPERATOR_NAME):$(VERSION)

CR_NAME		= $(shell echo $(APINAME)_$(KIND)_cr | sed 's|/|_|')
CRD_NAME	= $(shell echo $$(echo $(APINAME) | cut -d/ -f1)_$(KIND)_crd)


ECHO 			= echo -e
C_STD 		= $(shell echo -e "\033[0m")
C_RED			= $(shell echo -e "\033[91m")
C_GREEN 	= $(shell echo -e "\033[92m")
P 			 	= $(shell echo -e "\033[92m> \033[0m")
OK 				= $(shell echo -e "\033[92m[ OK  ]\033[0m")
ERROR	 		= $(shell echo -e "\033[91m[ERROR]\033[0m")
PASS	 		= $(shell echo -e "\033[92m[PASS ]\033[0m")
FAIL	 		= $(shell echo -e "\033[91m[FAIL ]\033[0m")

default: generate build-image

all:

generate:
	operator-sdk generate crds
	operator-sdk generate k8s
	kubectl apply -f deploy/crds/$(CRD_NAME).yaml

test-local:
	OPERATOR_NAME=$(OPERATOR_NAME) operator-sdk run local --watch-namespace=default

build-image:
	operator-sdk build johandry/$(IMAGE)
	docker push johandry/$(IMAGE)

setup-deploy:
	if [[ -e deploy/operator.yaml.org ]]; then mv deploy/operator.yaml.org deploy/operator.yaml 2>/dev/null; fi
	sed -i.org 's|REPLACE_IMAGE|$(IMAGE)|g' deploy/operator.yaml

deploy: setup-deploy
	kubectl create -f deploy/service_account.yaml
	kubectl create -f deploy/role.yaml
	kubectl create -f deploy/role_binding.yaml
	kubectl create -f deploy/operator.yaml

deploy-cr:
	kubectl create -f deploy/crds/$(CR_NAME).yaml

clean:
	kubectl delete -f deploy/operator.yaml
	kubectl delete -f deploy/role_binding.yaml
	kubectl delete -f deploy/role.yaml
	kubectl delete -f deploy/service_account.yaml
	kubectl delete -f deploy/crds/$(CR_NAME).yaml
	kubectl delete -f deploy/crds/$(CRD_NAME).yaml

add-api:
	operator-sdk add api --kind $(CAP_KIND) --api-version $(APINAME)

add-controller:
	operator-sdk add controller --kind $(CAP_KIND) --api-version $(APINAME)

add: add-api add-controller

install:
	brew install operator-sdk
