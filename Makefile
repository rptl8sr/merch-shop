REPO_NAME := $(shell basename $(CURDIR))
PROJECT := $(CURDIR)
LOCAL_BIN := $(CURDIR)/bin
MIGRATIONS_DIR := $(CURDIR)/migrations

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

# PROJECT
.PHONY: blueprint
blueprint:
	mkdir $(LOCAL_BIN)
	mkdir $(MIGRATIONS_DIR)
	mkdir -p cmd/api && echo 'package main' > cmd/api/main.go
	mkdir -p internal/config && echo 'package config' > internal/config/config.go
	mkdir -p internal/handler && echo 'package handler' > internal/handler/handler.go
	mkdir -p internal/model && echo 'package model' > internal/model/model.go
	mkdir -p internal/repository && echo 'package repository' > internal/repository/repository.go
	mkdir -p internal/middleware && echo 'package middleware' > internal/middleware/middleware.go
	mkdir -p pkg/jwt && echo 'package jwt' > pkg/jwt/jwt.go
	mkdir -p pkg/database && echo 'package database' > pkg/database/database.go


