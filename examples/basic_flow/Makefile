ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

AWSBI := epiphanyplatform/awsbi:0.0.1
AWSKS := epiphanyplatform/awsks:0.0.1

#makes it easier to replace the value
M_NAME            ?= awsks-basic-flow
M_SHARED_EXTERNAL ?= $(ROOT_DIR)/shared

VMS_RSA_FILENAME ?= vms_rsa

-include ./awsks.mk

export

.PHONY: all apply destroy clean setup

all: apply

apply: init-awsbi plan-awsbi apply-awsbi \
       init-awsks plan-awsks apply-awsks

destroy: plan-destroy-awsks destroy-awsks \
         plan-destroy-awsbi destroy-awsbi

clean:
	@rm -rf $(M_SHARED_EXTERNAL)

.PHONY: init-awsbi plan-awsbi apply-awsbi output-awsbi plan-destroy-awsbi destroy-awsbi

init-awsbi: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSBI) \
		init \
		M_VMS_COUNT=0 \
		M_PUBLIC_IPS=false \
		M_NAME=$(M_NAME)

plan-awsbi: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSBI) \
		plan \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY

apply-awsbi: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSBI) \
		apply \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY

output-awsbi: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSBI) \
		output \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY

plan-destroy-awsbi: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSBI) \
		plan-destroy \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY

destroy-awsbi: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSBI) \
		destroy \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY

.PHONY: init-awsks plan-awsks apply-awsks output-awsks plan-destroy-awsks destroy-awsks

init-awsks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSKS) \
		init \
		M_NAME=$(M_NAME)

plan-awsks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSKS) \
		plan \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY

apply-awsks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSKS) \
		apply \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY \
		kubeconfig

output-awsks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSKS) \
		output \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY

plan-destroy-awsks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSKS) \
		plan-destroy \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY

destroy-awsks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSKS) \
		destroy \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY

setup: $(M_SHARED_EXTERNAL)/$(VMS_RSA_FILENAME)

$(M_SHARED_EXTERNAL)/$(VMS_RSA_FILENAME):
	@mkdir -p $(dir $@)
	@ssh-keygen -t rsa -b 4096 -f $@ -N '' 2>&1 >/dev/null
