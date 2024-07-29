VERSION := development
IMAGE := ghcr.io/kuoss/venti:$(VERSION)
GOLANGCI_LINT_VER := v1.59.1

MAKEFLAGS += -j2

datasources:
	hack/genernate-logs.sh
	docker ps | grep lethe        || docker run -d -p6060:6060 --name lethe -v /tmp/log:/var/data/log ghcr.io/kuoss/lethe
	docker ps | grep prometheus   || docker run -d -p9090:9090 --name prometheus   prom/prometheus
	docker ps | grep alertmanager || docker run -d -p9093:9093 --name alertmanager prom/alertmanager

install-dev:
	hack/install-dev.sh

# dev server (port 5173)
dev:
	hack/dev.sh

kill-dev:
	fuser 3030/tcp && kill -9 `fuser 3030/tcp | awk '{print $1}'` || true
	fuser 5173/tcp && kill -9 `fuser 5173/tcp | awk '{print $1}'` || true

# gin server (port 8080)
run-watch: run-watch-go run-watch-web
run-watch-go:
	cd web && npm run watch
run-watch-web:
	sleep 15 && air

# gin server (port 8080)
run-air:
	cd web && npm run build
	air


.PHONY: docker
docker:
	docker build -t $(IMAGE) --build-arg VERSION=$(VERSION) . && docker push $(IMAGE)

.PHONY: test
test:
	hack/test.sh

.PHONY: cover
cover:
	hack/test-cover.sh

.PHONY: checks
checks:
	hack/checks.sh

.PHONY: misspell
misspell:
	hack/misspell.sh

.PHONY: gocyclo
gocyclo:
	hack/gocyclo.sh

.PHONY: lint
lint:
	go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VER) || true
	golangci-lint run
