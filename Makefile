VERSION := development
IMAGE := ghcr.io/kuoss/venti:$(VERSION)
GOLANGCI_LINT_VER := v1.59.1

MAKEFLAGS += -j2

.PHONY: datasources
datasources:
	hack/datasources/datasources.sh

.PHONY: install-dev
install-dev:
	hack/install-dev.sh

.PHONY: dev
dev:
	hack/dev.sh

.PHONY: kill-dev
kill-dev:
	fuser 3030/tcp && kill -9 `fuser 3030/tcp | awk '{print $1}'` || true
	fuser 5173/tcp && kill -9 `fuser 5173/tcp | awk '{print $1}'` || true

.PHONY: run-watch
run-watch: run-watch-go run-watch-web

.PHONY: run-watch-go
run-watch-go:
	cd web && npm run watch

.PHONY: run-watch-web
run-watch-web:
	sleep 15 && air

.PHONY: run-air
run-air:
	cd web && npm run build
	air

.PHONY: docker
docker: docker-build docker-push

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE) --build-arg VERSION=$(VERSION) .

.PHONY: docker-push
docker-push:
	docker push $(IMAGE)

.PHONY: test
test:
	hack/test.sh

.PHONY: cover
cover:
	hack/test-cover.sh

.PHONY: checks
checks: test lint

.PHONY: lint
lint:
	go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VER) || true
	$(shell go env GOPATH)/bin/golangci-lint run
