VERSION := development
IMG ?= ghcr.io/kuoss/venti:$(VERSION)

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
lint: golangci-lint ## Run golangci-lint linter
	$(GOLANGCI_LINT) run -v

##@ Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint

## Tool Versions
GOLANGCI_LINT_VERSION ?= v1.60.3

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef