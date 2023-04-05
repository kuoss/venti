VENTI_VERSION=v0.1.13
IMAGE_REPO=ghcr.io/kuoss

LDFLAGS += -X "main.ventiVersion=$(VENTI_VERSION)"
MAKEFLAGS += -j2

install-dev:
	go mod tidy
	cd web && npm install
	go install github.com/cosmtrek/air@latest

# dev server (port 5173)
run-dev: run-dev-go run-dev-web
run-dev-go:
	API_ONLY=1 VENTI_VERSION=${VENTI_VERSION} air
run-dev-web:
	cd web && npm run dev --clearScreen=false

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

# on codespace
stage:
	skaffold dev --namespace=kube-system --default-repo=ghcr.io/kuoss

#go-build:
#	go mod download -x && go build -ldflags '$(LDFLAGS)' -o /app/venti

docker-build:
	docker build -t ${IMAGE_REPO}/venti:${VENTI_VERSION} --build-arg VENTI_VERSION=${VENTI_VERSION} . && docker push ${IMAGE_REPO}/venti:${VENTI_VERSION} 

pre-checks:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

checks:
	./scripts/checks.sh

fmt:
	go fmt ./...

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

golangci-lint:
	golangci-lint run --timeout 5m

test-cover:
	./scripts/test-cover.sh

