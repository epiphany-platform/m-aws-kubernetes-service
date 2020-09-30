include $(M_RESOURCES)/templates.mk
include $(M_RESOURCES)/consts.mk
include $(M_RESOURCES)/defaults.mk

export

#medatada method is printing static metadata information about module
metadata: guard-M_RESOURCES
	#AWSKS | metadata | should print component metadata
	@echo "$$M_METADATA_CONTENT"

#init method is used to initialize module configuration and check if state is providing strong (and weak) dependencies
#TODO should also validate state if strong requirements are met
init: guard-M_RESOURCES guard-M_SHARED setup ensure-state-file template-config-file initialize-state-file display-config-file

#plan method would get config file and environment state file and compare them and calculate what would be done o apply stage
plan: guard-M_RESOURCES guard-M_SHARED setup validate-config validate-state template-tfvars module-plan terraform-plan

#apply method runs module provider logic using config file
apply: guard-M_RESOURCES guard-M_SHARED setup module-plan terraform-apply update-state-after-apply terraform-output

#audit method should call logic to check if remote components are in "known" state
#TODO implement validation if remote resources are as expected, possibly with terraform plan
audit:
	#AWSKS | audit | should output current state of remote components

destroy: template-tfvars terraform-destroy update-state-after-destroy

plan-destroy: template-tfvars terraform-plan-destroy

output: terraform-output

#TODO consider parsing terraform plan output
terraform-plan:
	#AWSKS | terraform-plan | will run plan
	@cd $(M_RESOURCES)/terraform ; \
	TF_IN_AUTOMATION=true \
	AWS_ACCESS_KEY_ID=$(M_AWS_ACCESS_KEY) \
	AWS_SECRET_ACCESS_KEY=$(M_AWS_SECRET_KEY) \
		terraform plan \
		-no-color \
		-input=false \
		-var-file=$(M_RESOURCES)/terraform/vars.tfvars.json \
		-state=$(M_SHARED)/$(M_MODULE_SHORT)/terraform.tfstate \
		-out=$(M_SHARED)/$(M_MODULE_SHORT)/terraform-apply.tfplan \
		$(M_RESOURCES)/terraform

terraform-plan-json:
	#AWSKS | terraform-plan-json | will show plan in json
	@cd $(M_RESOURCES)/terraform ; \
	TF_IN_AUTOMATION=true \
		terraform show \
		-no-color \
		-json \
		$(M_SHARED)/$(M_MODULE_SHORT)/terraform-apply.tfplan

terraform-apply:
	#AWSKS | terraform-apply | will run terraform apply
	@cd $(M_RESOURCES)/terraform ; \
	TF_IN_AUTOMATION=true \
	AWS_ACCESS_KEY_ID=$(M_AWS_ACCESS_KEY) \
	AWS_SECRET_ACCESS_KEY=$(M_AWS_SECRET_KEY) \
		terraform apply \
		-no-color \
		-input=false \
		-auto-approve \
		-state=$(M_SHARED)/$(M_MODULE_SHORT)/terraform.tfstate \
		$(M_SHARED)/$(M_MODULE_SHORT)/terraform-apply.tfplan

terraform-plan-destroy:
	#AWSKS | terraform-plan-destroy | will prepare plan of destruction
	@cd $(M_RESOURCES)/terraform ; \
	TF_IN_AUTOMATION=true \
	AWS_ACCESS_KEY_ID=$(M_AWS_ACCESS_KEY) \
	AWS_SECRET_ACCESS_KEY=$(M_AWS_SECRET_KEY) \
		terraform plan \
		-destroy \
		-no-color \
		-input=false \
		-var-file=$(M_RESOURCES)/terraform/vars.tfvars.json \
		-state=$(M_SHARED)/$(M_MODULE_SHORT)/terraform.tfstate \
		-out=$(M_SHARED)/$(M_MODULE_SHORT)/terraform-destroy.tfplan \
		$(M_RESOURCES)/terraform

terraform-destroy:
	#AWSKS | terraform-destroy | terraform-destroy is just about to begin ...
	@cd $(M_RESOURCES)/terraform ; \
	TF_IN_AUTOMATION=true \
	TF_WARN_OUTPUT_ERRORS=1 \
	AWS_ACCESS_KEY_ID=$(M_AWS_ACCESS_KEY) \
	AWS_SECRET_ACCESS_KEY=$(M_AWS_SECRET_KEY) \
		terraform apply \
		-no-color \
		-input=false \
		-auto-approve \
		-state=$(M_SHARED)/$(M_MODULE_SHORT)/terraform.tfstate \
		$(M_SHARED)/$(M_MODULE_SHORT)/terraform-destroy.tfplan

terraform-output:
	#AWSKS | terraform-output | will prepare terraform output
	@cd $(M_RESOURCES)/terraform ; \
	TF_IN_AUTOMATION=true \
		terraform output \
		-no-color \
		-json \
		-state=$(M_SHARED)/$(M_MODULE_SHORT)/terraform.tfstate > $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.json
	@yq r --prettyPrint $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.json | yq r - --printMode pv '*.value' > $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.yml
	@yq p -i $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.yml $(M_MODULE_SHORT).output
	@yq m -x -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.yml
	@rm $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.json $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.yml

setup: $(M_SHARED)/$(M_MODULE_SHORT)
	#AWSKS | setup | ensure required directories

template-tfvars:
	#AWSKS | template-tfvars | will template .tfvars.json file
	@yq read -jP $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME) '$(M_MODULE_SHORT)*' > $(M_RESOURCES)/terraform/vars.tfvars.json

template-config-file:
	#AWSKS | template-config-file | will template config file (and backup previous if exists)
	@if test -f $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME); then mv $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME).backup ; fi
	@echo "$$M_CONFIG_CONTENT" > $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME)
	#AWSKS | template-config-file | will replace arguments with values from state file
	@yq r $(M_SHARED)/$(M_STATE_FILE_NAME) 'awsbi.output[vpc_id.value]' \
	| xargs -r yq w -i $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME) '$(M_MODULE_SHORT).vpc_id'
	@yq r $(M_SHARED)/$(M_STATE_FILE_NAME) 'awsbi.output[public_subnet_id.value]' \
	| xargs -r yq w -i $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME) '$(M_MODULE_SHORT).public_subnet_id'

display-config-file:
	#AWSKS | display-config-file | config file content is:
	@cat $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME)

ensure-state-file: $(M_SHARED)/$(M_STATE_FILE_NAME)
	#AWSKS | ensure-state-file | checks if state file exists

initialize-state-file:
	#AWSKS | initialize-state-file | will initialize state file
	@echo "$$M_STATE_INITIAL" > $(M_SHARED)/$(M_MODULE_SHORT)/awsbi-state.tmp
	@yq m -i -x $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/awsbi-state.tmp
	@rm $(M_SHARED)/$(M_MODULE_SHORT)/awsbi-state.tmp

update-state-after-apply:
	#AWSKS | update-state-after-apply | will update state file after apply
	@cp $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/awsbi-config.tmp.yml
	@yq d -i $(M_SHARED)/$(M_MODULE_SHORT)/awsbi-config.tmp.yml kind
	@yq m -x -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/awsbi-config.tmp.yml
	@yq w -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_MODULE_SHORT).status applied
	@rm $(M_SHARED)/$(M_MODULE_SHORT)/awsbi-config.tmp.yml

update-state-after-destroy:
	#AWSKS | update-state-after-destroy | will clean state file after destroy
	@yq d -i $(M_SHARED)/$(M_STATE_FILE_NAME) '$(M_MODULE_SHORT)'
	@yq w -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_MODULE_SHORT).status destroyed

#TODO check if there is state file
#TODO check if there is config
assert-init-completed:
	#AWSKS | assert-init-completed | will check if all initialization steps are completed

#TODO validate if config is correct
#TODO consider https://github.com/santhosh-tekuri/jsonschema as it's small
validate-config:
	#AWSKS | validate-config | will perform config validation

#TODO validate if state file is correct
#TODO consider https://github.com/santhosh-tekuri/jsonschema as it's small
validate-state:
	#AWSKS | validate-state | will perform state file validation

module-plan:
	#AWSKS | module-plan | will perform module plan
	@yq m -x $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME) > $(M_SHARED)/$(M_MODULE_SHORT)/awsbi-future-state.tmp
	@yq w -i $(M_SHARED)/$(M_MODULE_SHORT)/awsbi-future-state.tmp kind state
	@- yq compare $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/awsbi-future-state.tmp
	@rm $(M_SHARED)/$(M_MODULE_SHORT)/awsbi-future-state.tmp

kubeconfig:
	#AWSKS | kubeconfig | will store kubeconfig in a file
	@yq r $(M_SHARED)/$(M_STATE_FILE_NAME) 'awsks.output[kubeconfig.value]' > $(M_SHARED)/kubeconfig

guard-%:
	@if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

$(M_SHARED)/$(M_STATE_FILE_NAME):
	@mkdir -p $(dir $@)
	@touch $@

$(M_SHARED)/$(M_MODULE_SHORT):
	@mkdir -p $@
