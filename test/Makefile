SHELL				= /bin/bash
ECHO 				= echo -e
C_STD 			= $(shell echo -e "\033[0m")
C_RED				= $(shell echo -e "\033[91m")
C_GREEN 		= $(shell echo -e "\033[92m")
P 			 		= $(shell echo -e "\033[92m> \033[0m")
OK 			 		= $(shell echo -e "\033[92m[ OK  ]\033[0m")
ERROR		 		= $(shell echo -e "\033[91m[ERROR] \033[0m")

default: apply

all: check apply

## Checks:

check-ibmcloud:
	@if ibmcloud --version | grep -q 'ibmcloud version'; then $(ECHO) "$(P) $(OK) ibmcloud"; else $(ECHO) "$(P) $(ERROR) ibmcloud"; exit 1; fi
	@if ibmcloud plugin show infrastructure-service | grep -q 'Plugin Version'; then $(ECHO) "$(P) $(OK) plugin infrastructure-service"; else $(ECHO) "$(P) $(ERROR) plugin infrastructure-service"; exit 1; fi

check-api-key:
	@if [[ -n $$IC_API_KEY ]]; then $(ECHO) "$(P) $(OK) IC_API_KEY"; else $(ECHO) "$(P) $(ERROR) IC_API_KEY"; exit 1; fi

check-terraform:
	@if terraform version | grep -q 'Terraform v'; then $(ECHO) "$(P) $(OK) terraform"; else $(ECHO) "$(P) $(ERROR) terraform"; exit 1; fi
	@if ls ~/.terraform.d/plugins/terraform-provider-ibm_* | grep -q 'provider'; then $(ECHO) "$(P) $(OK) IBM provider"; else $(ECHO) "$(P) $(ERROR) IBM provider"; exit 1; fi

print-target:
	@$(ECHO) "$(P) ibmcloud target: "; ibmcloud target

check-kubectl:
	@if kubectl version | grep -q 'Client Version'; then $(ECHO) "$(P) $(OK) kubectl"; else $(ECHO) "$(P) $(ERROR) kubectl"; exit 1; fi

check-kubernetes:
	@if kubectl cluster-info | grep -q 'Kubernetes master is running at'; then $(ECHO) "$(P) $(OK) Kubernetes cluster"; else $(ECHO) "$(P) $(ERROR) Kubernetes cluster"; exit 1; fi

check: check-ibmcloud check-init check-api-key check-kubectl print-target

## Init

api-key: check-ibmcloud
	ibmcloud iam api-key-create TerraformKey -d "API Key for Terraform" --file terraform_key.json

init: check-terraform check-api-key
	cd terraform; \
	terraform init && \
	terraform validate && \
	terraform plan

## Apply

apply-terraform: check-terraform check-ibmcloud check-api-key
	cd terraform; terraform apply -auto-approve
	ibmcloud ks cluster config -cluster $$(terraform output -state=terraform/terraform.tfstate cluster_id)

apply-provisioner:
	kubectl apply -f kubernetes/nfs-provisioner

apply-operator:
	kubectl apply -f kubernetes/nfs-operator

apply-pvc:
	kubectl apply -f kubernetes/pvc.yaml

apply-kubernetes: check-kubectl check-kubernetes apply-pvc apply-provisioner

apply: apply-terraform apply-kubernetes

## Cleanup

destroy:
	cd terraform; terraform destroy -auto-approve

delete:
	kubectl delete -f kubernetes/pvc.yaml
	kubectl delete -f kubernetes/nfs-provisioner
	kubectl delete -f kubernetes/nfs-operator

clean: delete destroy
	$(RM) -rf terraform/.terraform
