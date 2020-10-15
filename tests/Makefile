SHELL := $(shell which bash)
SELF  := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

-include $(SELF)/../awscreds.mk

export

.PHONY: all

all: test

.PHONY: test TestInit TestPlan TestApply

test: prereq
	@cd $(SELF)/ && go test -v -timeout 45m

TestInit: prereq
	@cd $(SELF)/ && go test -v -timeout 5m -run $@

TestPlan: prereq
	@cd $(SELF)/ && go test -v -timeout 10m -run $@

TestApply: prereq
	@cd $(SELF)/ && go test -v -timeout 30m -run $@

.PHONY: prereq

prereq: needs-aws-iam-authenticator \
        needs-docker \
        needs-go \
        needs-kubectl

needs-%:
	@which $*
