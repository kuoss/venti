VERSION := v0.2.0
IMAGE := ghcr.io/kuoss/venti:$(VERSION)

MAKEFLAGS += -j2

# dev server (port 5173)
dev:
	hack/dev.sh

cleanup-dev:
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


docker:
	docker build -t $(IMAGE) --build-arg VERSION=$(VERSION) . && docker push $(IMAGE)

test:
	hack/test-failfast.sh

testall:
	hack/test-all.sh

cover:
	hack/test-cover.sh

checks:
	hack/checks.sh

misspell:
	hack/misspell.sh

gocyclo:
	hack/gocyclo.sh
