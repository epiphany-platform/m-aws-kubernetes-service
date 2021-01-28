ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

VERSION ?= 0.0.1
USER := epiphanyplatform
IMAGE := awsks

IMAGE_NAME := $(USER)/$(IMAGE):$(VERSION)

define AWS_CREDENTIALS_CONTENT
AWS_ACCESS_KEY ?= $(ACCESS_KEY)
AWS_SECRET_KEY ?= $(SECRET_KEY)
endef

-include ./awscreds.mk

export

#used for correctly setting shared folder permissions
HOST_UID := $(shell id -u)
HOST_GID := $(shell id -g)

.PHONY: build test test-release prepare-aws-credentials metadata

build: guard-VERSION guard-IMAGE guard-USER needs-docker
	docker build \
		--build-arg ARG_M_VERSION=$(VERSION) \
		--build-arg ARG_HOST_UID=$(HOST_UID) \
		--build-arg ARG_HOST_GID=$(HOST_GID) \
		-t $(IMAGE_NAME) \
		.

#prepare AWS credentials file before running this target using `ACCESS_KEY=xxx SECRET_KEY=yyy make prepare-aws-credentials`
test: test-prerequisite \
	build
	@go test -v -timeout 90m

test-release: test-prerequisite \
	release
	@CGO_ENABLED=0 go test -v -timeout 45m

prepare-aws-credentials: guard-ACCESS_KEY guard-SECRET_KEY
	@echo "$$AWS_CREDENTIALS_CONTENT" > $(ROOT_DIR)/awscreds.mk

release: guard-VERSION guard-IMAGE guard-USER
	docker build \
		--build-arg ARG_M_VERSION=$(VERSION) \
		-t $(IMAGE_NAME) \
		.

metadata: guard-VERSION guard-IMAGE guard-USER needs-docker
	@docker run --rm \
		-t $(IMAGE_NAME) \
		metadata

print-%:
	@echo "$($*)"

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

test-prerequisite: guard-AWS_ACCESS_KEY guard-AWS_SECRET_KEY \
		needs-docker \
		needs-go \
		needs-kubectl \
		needs-aws

needs-%:
	@which $*
