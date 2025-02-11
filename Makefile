REPO_NAME := $(shell basename $(CURDIR))
PROJECT := $(CURDIR)
LOCAL_BIN := $(CURDIR)/bin

ifneq (,$(wildcard ./.env))
ENV_FILE := .env
else
ENV_FILE :=
endif

ifneq ($(ENV_FILE),)
include $(ENV_FILE)
export
endif

# GIT
.PHONY: git-init
git-init:
	gh repo create $(USER)/$(REPO_NAME) --private
	git init
	git config user.name "$(USER)"
	git config user.email "$(EMAIL)"
	git add Makefile go.mod
	git commit -m "Init commit"
	git remote add origin git@github.com:$(USER)/$(REPO_NAME).git
	git remote -v
	git push -u origin master

BN ?= dev
.PHONY: git-checkout
git-checkout:
	git checkout -b $(BN)
	
# LINT
.PHONY: golangci-lint-install
lint-install:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2

.PHONY: lint
lint:
	$(LOCAL_BIN)/golangci-lint run ./...