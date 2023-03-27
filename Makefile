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

# on codespace
stage:
	skaffold dev --namespace=kube-system --default-repo=ghcr.io/kuoss

#go-build:
#	go mod download -x && go build -ldflags '$(LDFLAGS)' -o /app/venti

docker-build:
	docker build -t ${IMAGE_REPO}/venti:${VENTI_VERSION} --build-arg VENTI_VERSION=${VENTI_VERSION} . && docker push ${IMAGE_REPO}/venti:${VENTI_VERSION} 


checks: fmt vet staticcheck golangci-lint go-licenses js-licenses test-cover

fmt:
	go fmt ./...

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

golangci-lint:
	golangci-lint run

go-licenses:
	# go install github.com/google/go-licenses@latest
	go-licenses report github.com/kuoss/venti | tee docs/go-licenses.csv;\
	go-licenses check github.com/kuoss/venti && echo OK

js-licenses:
	# npm install -g js-green-licenses
	cd web && jsgl --local . && echo OK
	
test-cover:
	./hack/test-cover.sh

