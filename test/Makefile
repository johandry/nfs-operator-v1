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

default: apply

all: apply test

## Checks:

check-ibmcloud:
	@if ibmcloud --version | grep -q 'ibmcloud version'; then $(ECHO) "$(P) $(OK) ibmcloud"; else $(ECHO) "$(P) $(ERROR) ibmcloud"; exit 1; fi
	@if ibmcloud plugin show infrastructure-service | grep -q 'Plugin Version'; then $(ECHO) "$(P) $(OK) plugin infrastructure-service"; else $(ECHO) "$(P) $(ERROR) plugin infrastructure-service"; exit 1; fi
	@if ibmcloud plugin show kubernetes-service | grep -q 'Plugin Version'; then $(ECHO) "$(P) $(OK) plugin kubernetes-service"; else $(ECHO) "$(P) $(ERROR) plugin kubernetes-service"; exit 1; fi

check-api-key:
	@if [[ -n $$IC_API_KEY ]]; then $(ECHO) "$(P) $(OK) IC_API_KEY"; else $(ECHO) "$(P) $(ERROR) IC_API_KEY";  exit 1; fi

check-terraform:
	@if terraform version | grep -q 'Terraform v'; then $(ECHO) "$(P) $(OK) terraform"; else $(ECHO) "$(P) $(ERROR) terraform"; exit 1; fi
	@if ls ~/.terraform.d/plugins/terraform-provider-ibm_* | grep -q 'provider'; then $(ECHO) "$(P) $(OK) IBM provider"; else $(ECHO) "$(P) $(ERROR) IBM provider"; exit 1; fi

check-kubectl:
	@if kubectl version | grep -q 'Client Version'; then $(ECHO) "$(P) $(OK) kubectl"; else $(ECHO) "$(P) $(ERROR) kubectl"; exit 1; fi

check-kubernetes:
	@if TERM=dumb kubectl cluster-info | grep -q 'Kubernetes master is running at'; then $(ECHO) "$(P) $(OK) Kubernetes cluster"; else $(ECHO) "$(P) $(ERROR) Kubernetes cluster"; exit 1; fi

wait-pvc:
	@printf "Waiting for IBM Cloud Block PVC to be Bound"
	@while kubectl get pvc nfs-block-custom | grep -q 'Pending'; do printf .; sleep 3; done; echo

check-pvc:
	@if kubectl get pvc nfs-block-custom | grep -q 'Bound'; then $(ECHO) "$(P) $(OK) IBM Cloud Block PVC"; else $(ECHO) "$(P) $(ERROR) IBM Cloud Block PVC"; exit 1; fi

wait-nfs:
	@printf "Waiting for NFS PVC to be Bound"
	@while kubectl get pvc nfs | grep -q 'Pending'; do printf .; sleep 3; done; echo

check-nfs:
	@if kubectl get pvc nfs | grep -q 'Bound'; then $(ECHO) "$(P) $(OK) NFS PVC"; else $(ECHO) "$(P) $(ERROR) NFS PVC"; exit 1; fi

wait-consumer:
	@printf "Waiting for consumer service"
	@while [[ $$($(MAKE) get-endpoint) != *.cloud ]]; do printf .; sleep 3; done; echo
	@printf "Waiting for consumer pod"
	@while kubectl get pods | grep movies | grep -q -v Running; do printf .; sleep 3; done; echo

	kubectl get pods | grep movies | grep -q Running

check-consumer:
	@if [[ $$($(MAKE) get-endpoint) == *.cloud ]]; then $(ECHO) "$(P) $(OK) Consumer Service"; else $(ECHO) "$(P) $(ERROR) Consumer Service"; exit 1; fi

print-target:
	@$(ECHO) "$(P) ibmcloud target: "; ibmcloud target

check: check-ibmcloud check-terraform check-api-key check-kubectl print-target

## Init

api-key:
	ibmcloud iam api-key-create TerraformKey -d "API Key for Terraform" --file .terraform_key.json

ibmcloud-target:
	@if [[ -e .target_account ]]; then $(ECHO) "$(P) $(OK) target account"; else $(ECHO) "$(P) $(ERROR) missing target account"; exit 1; fi
	ibmcloud target -c $$(cat .target_account) -r us-south -g Default

init: check-terraform ibmcloud-target check-api-key
	cd terraform; \
	terraform init && \
	terraform validate && \
	terraform plan

## Apply

apply-terraform: check-terraform check-ibmcloud check-api-key
	cd terraform; terraform apply -auto-approve
	ibmcloud target -g $$(terraform output -state=terraform/terraform.tfstate resource_group)
	ibmcloud ks cluster config -cluster $$(terraform output -state=terraform/terraform.tfstate cluster_id)

apply-$(RESOURCE): check-pvc
	kubectl apply -f kubernetes/nfs-$(RESOURCE)/deployment.yaml
	kubectl apply -f kubernetes/nfs-$(RESOURCE)/rbac.yaml
	kubectl apply -f kubernetes/nfs-$(RESOURCE)/class.yaml
	kubectl apply -f kubernetes/nfs-$(RESOURCE)/claim.yaml
	$(MAKE) wait-nfs

apply-pvc:
	kubectl apply -f kubernetes/pvc.yaml
	@$(MAKE) wait-pvc

apply-consumer: docker-container check-nfs
	kubectl apply -f kubernetes/consumer/cm.yaml
	kubectl apply -f kubernetes/consumer/deployment.yaml
	$(MAKE) wait-consumer

apply-kubernetes: check-kubectl check-kubernetes apply-pvc

apply: init apply-terraform apply-kubernetes apply-$(RESOURCE) apply-consumer

## Build

docker-container:
	docker build -t johandry/movies:1.0 -f kubernetes/consumer/Dockerfile .
	docker push johandry/movies:1.0

get-endpoint:
	@kubectl get svc movies -o=jsonpath='{.status.loadBalancer.ingress[0].hostname}'

## Test

test: check-pvc check-nfs
	@$(ECHO) "$(P) Testing:"
	@ADDRESS=$$($(MAKE) get-endpoint); \
	if [[ $$(curl -s $$ADDRESS/movies/1067 | grep title | sed 's/.*: "\(.*\)",/\1/') == "Seven Samurai" ]]; then $(ECHO) "\t$(PASS) Check Movie"; else $(ECHO) "\t$(FAIL) Check Movie"; exit 1; fi

## Cleanup

destroy:
	cd terraform; terraform destroy -auto-approve

delete:
	kubectl delete -f kubernetes/consumer
	kubectl delete -f kubernetes/nfs-provisioner
	kubectl delete -f kubernetes/pvc.yaml

remove:
	$(RM) -rf terraform/.terraform
	$(RM) -rf terraform/terraform.tfstate*

clean: delete destroy remove
