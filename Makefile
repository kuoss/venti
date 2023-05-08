VERSION := v0.1.13
IMAGE := ghcr.io/kuoss/venti:$(VERSION)

MAKEFLAGS += -j2

install-dev:
	go mod tidy
	cd web && npm install
	which air   || go install github.com/cosmtrek/air@latest
	which godef || go install github.com/rogpeppe/godef@lates

mock-prometheus:
	docker rm -f prometheus; docker run -d -p9090:9090 --name prometheus prom/prometheus

# dev server (port 5173)
run-dev:
	cp docs/examples/datasources_dev.yml etc/datasources.yml
	hack/run-dev.sh

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


docker:
	docker build -t $(IMAGE) --build-arg VERSION=$(VERSION) . && docker push $(IMAGE)

test:
	cp docs/examples/datasources_test.yml etc/datasources.yml
	hack/test-failfast.sh

cover:
	cp docs/examples/datasources_test.yml etc/datasources.yml
	hack/test-cover.sh

checks:
	cp docs/examples/datasources_test.yml etc/datasources.yml
	hack/checks.sh

